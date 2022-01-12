package async

import (
	"context"
	"fmt"
	"time"

	"github.com/fleetdm/fleet/v4/server/datastore/redis"
	"github.com/fleetdm/fleet/v4/server/fleet"
	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	redigo "github.com/gomodule/redigo/redis"
)

const collectorLockKey = "locks:async_collector:{%s}"

type Task struct {
	Datastore fleet.Datastore
	Pool      fleet.RedisPool
	// AsyncEnabled indicates if async processing is enabled in the
	// configuration. Note that Pool can be nil if this is false.
	AsyncEnabled bool

	LockTimeout        time.Duration
	LogStatsInterval   time.Duration
	InsertBatch        int
	DeleteBatch        int
	UpdateBatch        int
	RedisPopCount      int
	RedisScanKeysCount int
	CollectorInterval  time.Duration
}

// Collect runs the various collectors as distinct background goroutines if
// async processing is enabled.  Each collector will stop processing when ctx
// is done.
func (t *Task) StartCollectors(ctx context.Context, jitterPct int, logger kitlog.Logger) {
	if !t.AsyncEnabled {
		level.Debug(logger).Log("task", "async disabled, not starting collectors")
		return
	}
	level.Debug(logger).Log("task", "async enabled, starting collectors", "interval", t.CollectorInterval, "jitter", jitterPct)

	collectorErrHandler := func(name string, err error) {
		level.Error(logger).Log("err", fmt.Sprintf("%s collector", name), "details", err)
	}

	labelColl := &collector{
		name:         "collect_labels",
		pool:         t.Pool,
		ds:           t.Datastore,
		execInterval: t.CollectorInterval,
		jitterPct:    jitterPct,
		lockTimeout:  t.LockTimeout,
		handler:      t.collectLabelQueryExecutions,
		errHandler:   collectorErrHandler,
	}

	policyColl := &collector{
		name:         "collect_policies",
		pool:         t.Pool,
		ds:           t.Datastore,
		execInterval: t.CollectorInterval,
		jitterPct:    jitterPct,
		lockTimeout:  t.LockTimeout,
		handler:      t.collectPolicyQueryExecutions,
		errHandler:   collectorErrHandler,
	}

	colls := []*collector{labelColl, policyColl}
	for _, coll := range colls {
		go coll.Start(ctx)
	}

	// log stats at regular intervals
	if t.LogStatsInterval > 0 {
		go func() {
			tick := time.Tick(t.LogStatsInterval)
			for {
				select {
				case <-tick:
					for _, coll := range colls {
						stats := coll.ReadStats()
						level.Debug(logger).Log("stats", fmt.Sprintf("%#v", stats), "name", coll.name)
					}
				case <-ctx.Done():
					return
				}
			}
		}()
	}
}

func storePurgeActiveHostID(pool fleet.RedisPool, zsetKey string, hid uint, reportedAt, purgeOlder time.Time) error {
	// KEYS[1]: the zsetKey
	// ARGV[1]: the host ID to add
	// ARGV[2]: the added host's reported-at timestamp
	// ARGV[3]: purge any entry with score older than this (purgeOlder timestamp)
	script := redigo.NewScript(1, `
    redis.call('ZADD', KEYS[1], ARGV[2], ARGV[1])
    return redis.call('ZREMRANGEBYSCORE', KEYS[1], '-inf', ARGV[3])
  `)

	conn := pool.Get()
	defer conn.Close()

	if err := redis.BindConn(pool, conn, zsetKey); err != nil {
		return fmt.Errorf("bind redis connection: %w", err)
	}

	if _, err := script.Do(conn, zsetKey, hid, reportedAt.Unix(), purgeOlder.Unix()); err != nil {
		return fmt.Errorf("run redis script: %w", err)
	}
	return nil
}

type hostIDLastReported struct {
	HostID       uint
	LastReported int64 // timestamp in unix epoch
}

func loadActiveHostIDs(pool fleet.RedisPool, zsetKey string, scanCount int) ([]hostIDLastReported, error) {
	conn := redis.ConfigureDoer(pool, pool.Get())
	defer conn.Close()

	// using ZSCAN instead of fetching in one shot, as there may be 100K+ hosts
	// and we don't want to block the redis server too long.
	var hosts []hostIDLastReported
	cursor := 0
	for {
		res, err := redigo.Values(conn.Do("ZSCAN", zsetKey, cursor, "COUNT", scanCount))
		if err != nil {
			return nil, fmt.Errorf("scan active host ids: %w", err)
		}
		var hostVals []uint
		if _, err := redigo.Scan(res, &cursor, &hostVals); err != nil {
			return nil, fmt.Errorf("convert scan results: %w", err)
		}
		for i := 0; i < len(hostVals); i += 2 {
			hosts = append(hosts, hostIDLastReported{HostID: hostVals[i], LastReported: int64(hostVals[i+1])})
		}

		if cursor == 0 {
			// iteration completed
			return hosts, nil
		}
	}
}

func removeProcessedHostIDs(pool fleet.RedisPool, zsetKey string, batch []hostIDLastReported) error {
	// This script removes from the set of active hosts all those that still have
	// the same score as when the batch was read (via loadActiveHostIDs). This is
	// so that any host that would've reported new data since the call to
	// loadActiveHostIDs would *not* get deleted (as the score would change if
	// that was the case).
	//
	// Note that this approach is correct - in that it is safe and won't delete
	// any host that has unsaved reported data - but it is potentially slow, as
	// it needs to check the score of each member before deleting it. Should that
	// become too slow, we have some options:
	//
	// * split the batch in smaller, capped ones (that would be if the redis
	//   server gets blocked for too long processing a single batch)
	// * use ZREMRANGEBYSCORE to remove in one command all members with a score
	//   (reported-at timestamp) lower than the maximum timestamp in batch.
	//   While this would be almost certainly faster, it might be incorrect as
	//   new data could be reported with timestamps older than the maximum one,
	//   e.g. if the clocks are not exactly in sync between fleet instances, or
	//   if hosts report new data while the ZSCAN is going on and don't get picked
	//   up by the SCAN (this is possible, as part of the guarantees of SCAN).

	// KEYS[1]: zsetKey
	// ARGV...: the list of host ID-last reported timestamp pairs
	script := redigo.NewScript(1, `
    local count = 0
    for i = 1, #ARGV, 2 do
      local member, ts = ARGV[i], ARGV[i+1]
      if redis.call('ZSCORE', KEYS[1], member) == ts then
        count = count + 1
        redis.call('ZREM', KEYS[1], member)
      end
    end
    return count
  `)

	conn := pool.Get()
	defer conn.Close()

	if err := redis.BindConn(pool, conn, zsetKey); err != nil {
		return fmt.Errorf("bind redis connection: %w", err)
	}

	args := redigo.Args{zsetKey}
	for _, host := range batch {
		args = args.Add(host.HostID, host.LastReported)
	}
	if _, err := script.Do(conn, args...); err != nil {
		return fmt.Errorf("run redis script: %w", err)
	}
	return nil
}
