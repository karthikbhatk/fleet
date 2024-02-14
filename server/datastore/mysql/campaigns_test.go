package mysql

import (
	"context"
	"testing"
	"time"

	"github.com/WatchBeam/clock"
	"github.com/fleetdm/fleet/v4/server/fleet"
	"github.com/fleetdm/fleet/v4/server/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCampaigns(t *testing.T) {
	ds := CreateMySQLDS(t)

	cases := []struct {
		name string
		fn   func(t *testing.T, ds *Datastore)
	}{
		{"DistributedQuery", testCampaignsDistributedQuery},
		{"CleanupDistributedQuery", testCampaignsCleanupDistributedQuery},
		{"SaveDistributedQuery", testCampaignsSaveDistributedQuery},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			defer TruncateTables(t, ds)

			c.fn(t, ds)
		})
	}
}

func testCampaignsDistributedQuery(t *testing.T, ds *Datastore) {
	user := test.NewUser(t, ds, "Zach", "zwass@fleet.co", true)
	mockClock := clock.NewMockClock()
	query := test.NewQuery(t, ds, nil, "test", "select * from time", user.ID, false)
	campaign := test.NewCampaign(t, ds, query.ID, fleet.QueryRunning, mockClock.Now())

	{
		retrieved, err := ds.DistributedQueryCampaign(context.Background(), campaign.ID)
		require.Nil(t, err)
		assert.Equal(t, campaign.QueryID, retrieved.QueryID)
		assert.Equal(t, campaign.Status, retrieved.Status)
	}

	h1 := test.NewHost(t, ds, "foo.local", "192.168.1.10", "1", "1", mockClock.Now())
	h2 := test.NewHost(t, ds, "bar.local", "192.168.1.11", "2", "2", mockClock.Now().Add(-1*time.Hour))
	h3 := test.NewHost(t, ds, "baz.local", "192.168.1.12", "3", "3", mockClock.Now().Add(-13*time.Minute))

	l1 := fleet.LabelSpec{
		ID:    1,
		Name:  "label foo",
		Query: "query foo",
	}
	l2 := fleet.LabelSpec{
		ID:    2,
		Name:  "label bar",
		Query: "query bar",
	}
	err := ds.ApplyLabelSpecs(context.Background(), []*fleet.LabelSpec{&l1, &l2})
	require.Nil(t, err)

	checkTargets(t, ds, campaign.ID, fleet.HostTargets{})

	test.AddHostToCampaign(t, ds, campaign.ID, h1.ID)
	checkTargets(t, ds, campaign.ID, fleet.HostTargets{HostIDs: []uint{h1.ID}})

	test.AddLabelToCampaign(t, ds, campaign.ID, l1.ID)
	checkTargets(t, ds, campaign.ID, fleet.HostTargets{HostIDs: []uint{h1.ID}, LabelIDs: []uint{l1.ID}})

	test.AddLabelToCampaign(t, ds, campaign.ID, l2.ID)
	checkTargets(t, ds, campaign.ID, fleet.HostTargets{HostIDs: []uint{h1.ID}, LabelIDs: []uint{l1.ID, l2.ID}})

	test.AddHostToCampaign(t, ds, campaign.ID, h2.ID)
	test.AddHostToCampaign(t, ds, campaign.ID, h3.ID)

	checkTargets(t, ds, campaign.ID, fleet.HostTargets{HostIDs: []uint{h1.ID, h2.ID, h3.ID}, LabelIDs: []uint{l1.ID, l2.ID}})
}

func testCampaignsCleanupDistributedQuery(t *testing.T, ds *Datastore) {
	ctx := context.Background()
	user := test.NewUser(t, ds, "Zach", "zwass@fleet.co", true)
	mockClock := clock.NewMockClock()
	query := test.NewQuery(t, ds, nil, "test", "select * from time", user.ID, false)

	c1 := test.NewCampaign(t, ds, query.ID, fleet.QueryWaiting, mockClock.Now())
	c2 := test.NewCampaign(t, ds, query.ID, fleet.QueryRunning, mockClock.Now())

	// Cleanup and verify that nothing changed (because time has not
	// advanced)
	expired, recentInactive, err := ds.CleanupDistributedQueryCampaigns(ctx, mockClock.Now())
	require.Nil(t, err)
	assert.Equal(t, uint(0), expired)
	assert.Empty(t, recentInactive)

	{
		retrieved, err := ds.DistributedQueryCampaign(ctx, c1.ID)
		require.Nil(t, err)
		assert.Equal(t, c1.QueryID, retrieved.QueryID)
		assert.Equal(t, c1.Status, retrieved.Status)
	}
	{
		retrieved, err := ds.DistributedQueryCampaign(ctx, c2.ID)
		require.Nil(t, err)
		assert.Equal(t, c2.QueryID, retrieved.QueryID)
		assert.Equal(t, c2.Status, retrieved.Status)
	}

	// Add some executions

	mockClock.AddTime(1*time.Minute + 1*time.Second)

	// Cleanup and verify that the campaign was expired and executions
	// deleted appropriately
	expired, recentInactive, err = ds.CleanupDistributedQueryCampaigns(ctx, mockClock.Now())
	require.Nil(t, err)
	assert.Equal(t, uint(1), expired)
	assert.Equal(t, []uint{c1.ID}, recentInactive)
	{
		// c1 should now be complete
		retrieved, err := ds.DistributedQueryCampaign(ctx, c1.ID)
		require.Nil(t, err)
		assert.Equal(t, c1.QueryID, retrieved.QueryID)
		assert.Equal(t, fleet.QueryComplete, retrieved.Status)
	}
	{
		retrieved, err := ds.DistributedQueryCampaign(ctx, c2.ID)
		require.Nil(t, err)
		assert.Equal(t, c2.QueryID, retrieved.QueryID)
		assert.Equal(t, c2.Status, retrieved.Status)
	}

	mockClock.AddTime(24*time.Hour + 1*time.Second)

	// Cleanup and verify that the campaign was expired and executions
	// deleted appropriately
	expired, recentInactive, err = ds.CleanupDistributedQueryCampaigns(ctx, mockClock.Now())
	require.Nil(t, err)
	assert.Equal(t, uint(1), expired)
	assert.ElementsMatch(t, []uint{c1.ID, c2.ID}, recentInactive)
	{
		retrieved, err := ds.DistributedQueryCampaign(ctx, c1.ID)
		require.Nil(t, err)
		assert.Equal(t, c1.QueryID, retrieved.QueryID)
		assert.Equal(t, fleet.QueryComplete, retrieved.Status)
	}
	{
		// c2 should now be complete
		retrieved, err := ds.DistributedQueryCampaign(ctx, c2.ID)
		require.Nil(t, err)
		assert.Equal(t, c2.QueryID, retrieved.QueryID)
		assert.Equal(t, fleet.QueryComplete, retrieved.Status)
	}

	// simulate another old campaign created > 7 days ago
	c3 := test.NewCampaign(t, ds, query.ID, fleet.QueryWaiting, mockClock.Now().AddDate(0, 0, -8))
	{
		retrieved, err := ds.DistributedQueryCampaign(ctx, c3.ID)
		require.Nil(t, err)
		assert.Equal(t, c3.QueryID, retrieved.QueryID)
		assert.Equal(t, fleet.QueryWaiting, retrieved.Status)
	}

	// cleanup will mark c3 as completed because it was waiting for > 1 minute,
	// but it won't return it as recently inactive because it's too old a query.
	expired, recentInactive, err = ds.CleanupDistributedQueryCampaigns(ctx, mockClock.Now())
	require.Nil(t, err)
	assert.Equal(t, uint(1), expired)
	assert.ElementsMatch(t, []uint{c1.ID, c2.ID}, recentInactive)

	// cleanup again does not expire any new campaign and still returns the same
	// recently inactive campaigns
	expired, recentInactive, err = ds.CleanupDistributedQueryCampaigns(ctx, mockClock.Now())
	require.Nil(t, err)
	assert.Equal(t, uint(0), expired)
	assert.ElementsMatch(t, []uint{c1.ID, c2.ID}, recentInactive)

	// move time forward 7 days and cleanup again, this time it returns no recent
	// inactive campaigns
	mockClock.AddTime(7*24*time.Hour + 1*time.Second)

	expired, recentInactive, err = ds.CleanupDistributedQueryCampaigns(ctx, mockClock.Now())
	require.Nil(t, err)
	assert.Equal(t, uint(0), expired)
	assert.Empty(t, recentInactive)
}

func testCampaignsSaveDistributedQuery(t *testing.T, ds *Datastore) {
	user := test.NewUser(t, ds, t.Name(), t.Name()+"zwass@fleet.co", true)

	mockClock := clock.NewMockClock()

	query := test.NewQuery(t, ds, nil, t.Name()+"test", "select * from time", user.ID, false)

	c1 := test.NewCampaign(t, ds, query.ID, fleet.QueryWaiting, mockClock.Now())
	gotC, err := ds.DistributedQueryCampaign(context.Background(), c1.ID)
	require.NoError(t, err)
	require.Equal(t, fleet.QueryWaiting, gotC.Status)

	c1.Status = fleet.QueryComplete
	require.NoError(t, ds.SaveDistributedQueryCampaign(context.Background(), c1))

	gotC, err = ds.DistributedQueryCampaign(context.Background(), c1.ID)
	require.NoError(t, err)
	require.Equal(t, fleet.QueryComplete, gotC.Status)
}

func checkTargets(t *testing.T, ds fleet.Datastore, campaignID uint, expectedTargets fleet.HostTargets) {
	targets, err := ds.DistributedQueryCampaignTargetIDs(context.Background(), campaignID)
	require.Nil(t, err)
	assert.ElementsMatch(t, expectedTargets.HostIDs, targets.HostIDs)
	assert.ElementsMatch(t, expectedTargets.LabelIDs, targets.LabelIDs)
	assert.ElementsMatch(t, expectedTargets.TeamIDs, targets.TeamIDs)
}
