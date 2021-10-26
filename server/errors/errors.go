package errors

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/fleetdm/fleet/v4/server/datastore/redis"
	"github.com/fleetdm/fleet/v4/server/fleet"
	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/rotisserie/eris"
)

// ErrorFlusher defines the method to implement to flush all stored errors and
// return them as a slice of JSON-encoded strings. Once flushed, existing
// errors are removed from the store. The *Handler type implements this
// interface.
type ErrorFlusher interface {
	Flush() ([]string, error)
}

var (
	testOnStore func(error) // if set, called each time an error is stored, for tests
	testOnStart func()      // if set, called once the handler is running, for tests
)

// Handler defines an error handler. Call Handler.New to handle an error, and
// Handler.Flush to retrieve all stored errors and clear them from the store.
// It is safe to call those methods concurrently.
type Handler struct {
	pool    fleet.RedisPool
	logger  kitlog.Logger
	running int32
	errCh   chan error
}

// NewHandler creates an error handler using the provided pool and logger,
// storing unique instances of errors in Redis using the pool. It stops storing
// errors when ctx is cancelled.
func NewHandler(ctx context.Context, pool fleet.RedisPool, logger kitlog.Logger) *Handler {
	ch := make(chan error, 1)

	eh := &Handler{
		pool:   pool,
		logger: logger,
		errCh:  ch,
	}
	go eh.handleErrors(ctx)
	return eh
}

// Flush retrieves all stored errors from Redis and returns them as a slice of
// JSON-encoded strings. It is a destructive read - the errors are removed from
// Redis on return.
func (h *Handler) Flush() ([]string, error) {
	errorKeys, err := redis.ScanKeys(h.pool, "error:*")
	if err != nil {
		return nil, err
	}

	keysBySlot := redis.SplitKeysBySlot(h.pool, errorKeys...)
	var errors []string
	for _, qkeys := range keysBySlot {
		if len(qkeys) > 0 {
			gotErrors, err := h.collectBatchErrors(qkeys)
			if err != nil {
				return nil, err
			}
			errors = append(errors, gotErrors...)
		}
	}
	return errors, nil
}

func (h *Handler) collectBatchErrors(errorKeys []string) ([]string, error) {
	conn := redis.ConfigureDoer(h.pool, h.pool.Get())
	defer conn.Close()

	var args redigo.Args
	args = args.AddFlat(errorKeys)
	errorList, err := redigo.Strings(conn.Do("MGET", args...))
	if err != nil {
		return nil, err
	}

	if _, err := conn.Do("DEL", args...); err != nil {
		return nil, err
	}

	return errorList, nil
}

func sha256b64(s string) string {
	src := sha256.Sum256([]byte(s))
	return base64.URLEncoding.EncodeToString(src[:])
}

func hashErrorLocation(err error) string {
	unpackedErr := eris.Unpack(err)
	if unpackedErr.ErrExternal != nil {
		// TODO(mna): this returns the same hash for the same error type+message
		// that happens in different places, e.g. if io.EOF is wrapped in different
		// locations, this would only save one instance. Not sure this is the
		// intention?
		s := fmt.Sprintf("%T\n%s", unpackedErr.ErrExternal, unpackedErr.ErrExternal.Error())
		return sha256b64(s)
	}

	if len(unpackedErr.ErrRoot.Stack) == 0 {
		return sha256b64(unpackedErr.ErrRoot.Msg)
	}

	lastFrame := unpackedErr.ErrRoot.Stack[0]
	return sha256b64(fmt.Sprintf("%s:%d", lastFrame.File, lastFrame.Line))
}

func hashErr(externalErr error) (errHash string, errAsJson string, err error) {
	m := eris.ToJSON(externalErr, true)
	bytes, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return "", "", err
	}
	return hashErrorLocation(externalErr), string(bytes), nil
}

func (h *Handler) handleErrors(ctx context.Context) {
	atomic.StoreInt32(&h.running, 1)
	defer func() {
		atomic.StoreInt32(&h.running, 0)
	}()

	if testOnStart != nil {
		testOnStart()
	}

	for {
		select {
		case <-ctx.Done():
			return
		case err := <-h.errCh:
			h.storeError(ctx, err)
		}
	}
}

func (h *Handler) storeError(ctx context.Context, err error) {
	errorHash, errorJson, err := hashErr(err)
	if err != nil {
		level.Error(h.logger).Log("err", err, "msg", "hashErr failed")
		if testOnStore != nil {
			testOnStore(err)
		}
		return
	}
	jsonKey := fmt.Sprintf("error:%s:json", errorHash)

	conn := redis.ConfigureDoer(h.pool, h.pool.Get())
	defer conn.Close()

	if _, err := conn.Do("SET", jsonKey, errorJson, "EX", int((24 * time.Hour).Seconds())); err != nil {
		level.Error(h.logger).Log("err", err, "msg", "redis SET failed")
		if testOnStore != nil {
			testOnStore(err)
		}
		return
	}

	if testOnStore != nil {
		testOnStore(nil)
	}
}

// New handles the provided error by storing it into Redis if the handler is
// still running. In any case, it always returns the error wrapped with an
// eris error (stack trace and extra information).
//
// If the provided ctx contains additional metadata (user, host, etc.), this is
// added to the error's.  If the ctx is cancelled before the error is stored,
// the call returns without storing the error, otherwise it waits for a
// predefined period of time to try to store the error.
func (h *Handler) New(ctx context.Context, err error) error {
	// TODO: wrap in eris error with other metadata
	err = eris.Wrapf(err, "timestamp: %v", time.Now().Format(time.RFC3339))
	if atomic.LoadInt32(&h.running) == 0 {
		return err
	}

	timer := time.NewTimer(2 * time.Second)
	defer timer.Stop()
	select {
	case h.errCh <- err:
	case <-timer.C:
	case <-ctx.Done():
	}

	return err
}

// NewHttpHandler creates an http.HandlerFunc that flushes the errors stored
// by the provided ErrorFlusher and returns them in the response as JSON.
func NewHttpHandler(eh ErrorFlusher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		errors, err := eh.Flush()
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		// TODO: I think those will be double-marshaled as the errors are stored as JSON
		// in Redis.
		bytes, err := json.Marshal(errors)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			return
		}
		w.Write(bytes)
	}
}
