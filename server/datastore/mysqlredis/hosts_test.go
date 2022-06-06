package mysqlredis

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/fleetdm/fleet/v4/server/datastore/redis/redistest"
	"github.com/fleetdm/fleet/v4/server/fleet"
	"github.com/fleetdm/fleet/v4/server/mock"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/require"
)

func TestEnforceHostLimit(t *testing.T) {
	const hostLimit = 3

	oldBatchSize := redisSetMembersBatchSize
	redisSetMembersBatchSize = 2
	defer func() { redisSetMembersBatchSize = oldBatchSize }()

	runTest := func(t *testing.T, pool fleet.RedisPool) {
		var hostIDSeq uint
		var expiredHostsIDs, incomingHostsIDs []uint

		ctx := context.Background()
		ds := new(mock.Store)
		ds.EnrollHostFunc = func(ctx context.Context, osqueryHostId, nodeKey string, teamID *uint, cooldown time.Duration) (*fleet.Host, error) {
			hostIDSeq++
			return &fleet.Host{
				ID: hostIDSeq, OsqueryHostID: osqueryHostId, NodeKey: nodeKey,
			}, nil
		}
		ds.NewHostFunc = func(ctx context.Context, host *fleet.Host) (*fleet.Host, error) {
			hostIDSeq++
			host.ID = hostIDSeq
			return host, nil
		}
		ds.DeleteHostFunc = func(ctx context.Context, hid uint) error {
			return nil
		}
		ds.DeleteHostsFunc = func(ctx context.Context, ids []uint) error {
			return nil
		}
		ds.CleanupExpiredHostsFunc = func(ctx context.Context) ([]uint, error) {
			return expiredHostsIDs, nil
		}
		ds.CleanupIncomingHostsFunc = func(ctx context.Context, now time.Time) ([]uint, error) {
			return incomingHostsIDs, nil
		}

		requireInvokedAndReset := func(flag *bool) {
			require.True(t, *flag)
			*flag = false
		}

		wrappedDS := New(ds, pool, WithEnforcedHostLimit(hostLimit))

		// create a few hosts within the limit
		h1, err := wrappedDS.NewHost(ctx, &fleet.Host{})
		require.NoError(t, err)
		require.NotNil(t, h1)
		requireInvokedAndReset(&ds.NewHostFuncInvoked)
		h2, err := wrappedDS.EnrollHost(ctx, "osquery-2", "node-2", nil, time.Second)
		require.NoError(t, err)
		require.NotNil(t, h2)
		requireInvokedAndReset(&ds.EnrollHostFuncInvoked)
		h3, err := wrappedDS.EnrollHost(ctx, "osquery-3", "node-3", nil, time.Second)
		require.NoError(t, err)
		require.NotNil(t, h3)
		requireInvokedAndReset(&ds.EnrollHostFuncInvoked)

		// creating a new one fails - the limit is reached
		_, err = wrappedDS.NewHost(ctx, &fleet.Host{})
		require.Error(t, err)
		require.Contains(t, err.Error(), "maximum number of hosts")
		require.False(t, ds.NewHostFuncInvoked)

		_, err = wrappedDS.EnrollHost(ctx, "osquery-4", "node-4", nil, time.Second)
		require.Error(t, err)
		require.Contains(t, err.Error(), "maximum number of hosts")
		require.False(t, ds.EnrollHostFuncInvoked)

		// deleting h1 allows h4 to be created
		err = wrappedDS.DeleteHost(ctx, h1.ID)
		require.NoError(t, err)
		h4, err := wrappedDS.EnrollHost(ctx, "osquery-4", "node-4", nil, time.Second)
		require.NoError(t, err)
		require.NotNil(t, h4)
		requireInvokedAndReset(&ds.EnrollHostFuncInvoked)

		// and then limit is reached again
		_, err = wrappedDS.EnrollHost(ctx, "osquery-5", "node-5", nil, time.Second)
		require.Error(t, err)
		require.Contains(t, err.Error(), "maximum number of hosts")
		require.False(t, ds.EnrollHostFuncInvoked)

		// delete h1-h2-h3 (even if h1 is already deleted) should allow 2 more
		err = wrappedDS.DeleteHosts(ctx, []uint{h1.ID, h2.ID, h3.ID})
		require.NoError(t, err)
		h5, err := wrappedDS.EnrollHost(ctx, "osquery-5", "node-5", nil, time.Second)
		require.NoError(t, err)
		require.NotNil(t, h5)
		requireInvokedAndReset(&ds.EnrollHostFuncInvoked)
		h6, err := wrappedDS.NewHost(ctx, &fleet.Host{})
		require.NoError(t, err)
		require.NotNil(t, h6)
		requireInvokedAndReset(&ds.NewHostFuncInvoked)
		_, err = wrappedDS.EnrollHost(ctx, "osquery-7", "node-7", nil, time.Second)
		require.Error(t, err)
		require.Contains(t, err.Error(), "maximum number of hosts")
		require.False(t, ds.EnrollHostFuncInvoked)

		// cleanup expired removes h4
		expiredHostsIDs = []uint{h4.ID}
		_, err = wrappedDS.CleanupExpiredHosts(ctx)
		require.NoError(t, err)
		// cleanup incoming removes h4, h5
		incomingHostsIDs = []uint{h4.ID, h5.ID}
		_, err = wrappedDS.CleanupIncomingHosts(ctx, time.Now())
		require.NoError(t, err)

		// can now create 2 more
		h7, err := wrappedDS.EnrollHost(ctx, "osquery-7", "node-7", nil, time.Second)
		require.NoError(t, err)
		require.NotNil(t, h7)
		requireInvokedAndReset(&ds.EnrollHostFuncInvoked)
		h8, err := wrappedDS.NewHost(ctx, &fleet.Host{})
		require.NoError(t, err)
		require.NotNil(t, h8)
		requireInvokedAndReset(&ds.NewHostFuncInvoked)
		_, err = wrappedDS.EnrollHost(ctx, "osquery-9", "node-9", nil, time.Second)
		require.Error(t, err)
		require.Contains(t, err.Error(), "maximum number of hosts")
		require.False(t, ds.EnrollHostFuncInvoked)
	}

	t.Run("standalone", func(t *testing.T) {
		pool := redistest.SetupRedis(t, enrolledHostsSetKey, false, false, false)
		runTest(t, pool)
	})

	t.Run("cluster", func(t *testing.T) {
		pool := redistest.SetupRedis(t, enrolledHostsSetKey, true, true, false)
		runTest(t, pool)
	})
}

func TestSyncEnrolledHostIDs(t *testing.T) {
	runTest := func(t *testing.T, pool fleet.RedisPool) {
		var hostIDSeq uint
		var enrolledHostCount int
		var enrolledHostIDs []uint
		ctx := context.Background()

		ds := new(mock.Store)
		ds.NewHostFunc = func(ctx context.Context, host *fleet.Host) (*fleet.Host, error) {
			hostIDSeq++
			host.ID = hostIDSeq
			return host, nil
		}
		ds.CountEnrolledHostsFunc = func(ctx context.Context) (int, error) {
			return enrolledHostCount, nil
		}
		ds.EnrolledHostIDsFunc = func(ctx context.Context) ([]uint, error) {
			return enrolledHostIDs, nil
		}

		requireInvokedAndReset := func(flag *bool) {
			require.True(t, *flag)
			*flag = false
		}

		wrappedDS := New(ds, pool, WithEnforcedHostLimit(10)) // limit is irrelevant for this test

		// create a few hosts kept in sync
		h1, err := wrappedDS.NewHost(ctx, &fleet.Host{})
		require.NoError(t, err)
		h2, err := wrappedDS.NewHost(ctx, &fleet.Host{})
		require.NoError(t, err)
		h3, err := wrappedDS.NewHost(ctx, &fleet.Host{})
		require.NoError(t, err)

		conn := pool.Get()
		defer conn.Close()

		redisIDs, err := redigo.Strings(conn.Do("SMEMBERS", enrolledHostsSetKey))
		require.NoError(t, err)
		require.ElementsMatch(t, []string{fmt.Sprint(h1.ID), fmt.Sprint(h2.ID), fmt.Sprint(h3.ID)}, redisIDs)

		// syncing with the correct count does not trigger a sync
		enrolledHostCount = 3
		err = SyncEnrolledHostIDs(ctx, ds, pool)
		require.NoError(t, err)
		requireInvokedAndReset(&ds.CountEnrolledHostsFuncInvoked)
		require.False(t, ds.EnrolledHostIDsFuncInvoked)

		// syncing with a non-matching count triggers a sync
		enrolledHostCount = 2
		enrolledHostIDs = []uint{h1.ID, h3.ID} // will set the redis key to those values
		err = SyncEnrolledHostIDs(ctx, ds, pool)
		require.NoError(t, err)
		requireInvokedAndReset(&ds.CountEnrolledHostsFuncInvoked)
		requireInvokedAndReset(&ds.EnrolledHostIDsFuncInvoked)

		redisIDs, err = redigo.Strings(conn.Do("SMEMBERS", enrolledHostsSetKey))
		require.NoError(t, err)
		require.ElementsMatch(t, []string{fmt.Sprint(h1.ID), fmt.Sprint(h3.ID)}, redisIDs)
	}

	t.Run("standalone", func(t *testing.T) {
		pool := redistest.SetupRedis(t, enrolledHostsSetKey, false, false, false)
		runTest(t, pool)
	})

	t.Run("cluster", func(t *testing.T) {
		pool := redistest.SetupRedis(t, enrolledHostsSetKey, true, true, false)
		runTest(t, pool)
	})
}
