package token

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/fleetdm/fleet/v4/orbit/pkg/constant"
	"github.com/fleetdm/fleet/v4/pkg/retry"
	"github.com/google/uuid"
)

type remoteUpdaterFunc func(token string) error

type ReadWriter struct {
	*Reader
	remoteUpdate remoteUpdaterFunc
}

func NewReadWriter(path string) *ReadWriter {
	return &ReadWriter{
		Reader: &Reader{Path: path},
	}
}

// LoadOrGenerate tries to read a token file from disk,
// and if it doesn't exist, generates a new one.
func (rw *ReadWriter) LoadOrGenerate() error {
	_, err := rw.Read()
	switch {
	case err == nil:
		// If the token file exists, but has expired, rotate it.
		if rw.HasExpired() {
			if err := rw.Rotate(); err != nil {
				return fmt.Errorf("rotating token on generation: %w", err)
			}
			return nil
		}

		// otherwise we're good to go, but first ensure the file is readable by
		// other processes, old versions of Orbit used to chmod this file with
		// 0o600
		if err := rw.setChmod(); err != nil {
			return fmt.Errorf("loading token file, chmod %q: %w", rw.Path, err)
		}
	case errors.Is(err, os.ErrNotExist):
		if err := rw.Rotate(); err != nil {
			return fmt.Errorf("rotating token on generation: %w", err)
		}
	default:
		return fmt.Errorf("load identifier file %q: %w", rw.Path, err)
	}
	return nil
}

// Rotate assigns a new value to the token and writes it to disk.
func (rw *ReadWriter) Rotate() error {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("generate identifier: %w", err)
	}

	id := uuid.String()
	if err := rw.Write(id); err != nil {
		return fmt.Errorf("writing token: %w", err)
	}

	attempts := 3
	interval := 5 * time.Second
	err = retry.Do(func() error {
		return rw.Write(id)
	}, retry.WithMaxAttempts(attempts), retry.WithInterval(interval))

	if err != nil {
		return fmt.Errorf("saving token after %d attempts: %w", attempts, err)
	}

	return nil
}

func (rw *ReadWriter) SetRemoteUpdateFunc(f remoteUpdaterFunc) {
	rw.remoteUpdate = f
}

func (rw *ReadWriter) Write(id string) error {
	if rw.remoteUpdate != nil {
		if err := rw.remoteUpdate(id); err != nil {
			return err
		}
	}

	err := os.WriteFile(rw.Path, []byte(id), constant.DefaultWorldReadableFileMode)
	if err != nil {
		return fmt.Errorf("write identifier file %q: %w", rw.Path, err)
	}

	// ensure the file is readable by other processes, os.WriteFile does not
	// modify permissions if the file already exists
	if err := rw.setChmod(); err != nil {
		return fmt.Errorf("write identifier file, chmod %q: %w", rw.Path, err)
	}

	// ensure the `mtime` is updated, we have seen tests fail in some versions of
	// Ubuntu because this value is not update when the file is written
	err = os.Chtimes(rw.Path, time.Now(), time.Now())
	if err != nil {
		return fmt.Errorf("set mtime of identifier file %q: %w", rw.Path, err)
	}

	// make sure we can read the token, and cache the value
	if _, err = rw.Read(); err != nil {
		return fmt.Errorf("read identifier file %q: %w", rw.Path, err)
	}

	return nil
}

func (rw *ReadWriter) setChmod() error {
	return os.Chmod(rw.Path, constant.DefaultWorldReadableFileMode)
}
