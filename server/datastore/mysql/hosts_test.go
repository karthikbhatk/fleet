package mysql

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/WatchBeam/clock"
	"github.com/fleetdm/fleet/v4/server/fleet"
	"github.com/fleetdm/fleet/v4/server/ptr"
	"github.com/fleetdm/fleet/v4/server/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var enrollTests = []struct {
	uuid, hostname, platform, nodeKey string
}{
	0: {uuid: "6D14C88F-8ECF-48D5-9197-777647BF6B26",
		hostname: "web.fleet.co",
		platform: "linux",
		nodeKey:  "key0",
	},
	1: {uuid: "B998C0EB-38CE-43B1-A743-FBD7A5C9513B",
		hostname: "mail.fleet.co",
		platform: "linux",
		nodeKey:  "key1",
	},
	2: {uuid: "008F0688-5311-4C59-86EE-00C2D6FC3EC2",
		hostname: "home.fleet.co",
		platform: "darwin",
		nodeKey:  "key2",
	},
	3: {uuid: "uuid123",
		hostname: "fakehostname",
		platform: "darwin",
		nodeKey:  "key3",
	},
}

func TestSaveHosts(t *testing.T) {
	ds := CreateMySQLDS(t)
	defer ds.Close()

	host, err := ds.NewHost(&fleet.Host{
		DetailUpdatedAt: time.Now(),
		LabelUpdatedAt:  time.Now(),
		SeenTime:        time.Now(),
		NodeKey:         "1",
		UUID:            "1",
		Hostname:        "foo.local",
		PrimaryIP:       "192.168.1.1",
		PrimaryMac:      "30-65-EC-6F-C4-58",
	})
	require.NoError(t, err)
	require.NotNil(t, host)

	host.Hostname = "bar.local"
	err = ds.SaveHost(host)
	require.Nil(t, err)

	host, err = ds.Host(host.ID)
	require.Nil(t, err)
	assert.Equal(t, "bar.local", host.Hostname)
	assert.Equal(t, "192.168.1.1", host.PrimaryIP)
	assert.Equal(t, "30-65-EC-6F-C4-58", host.PrimaryMac)

	additionalJSON := json.RawMessage(`{"foobar": "bim"}`)
	host.Additional = &additionalJSON

	require.NoError(t, ds.SaveHost(host))
	require.NoError(t, ds.SaveHostAdditional(host))

	host, err = ds.Host(host.ID)
	require.Nil(t, err)
	require.NotNil(t, host)
	require.NotNil(t, host.Additional)
	assert.Equal(t, additionalJSON, *host.Additional)

	err = ds.SaveHost(host)
	require.Nil(t, err)

	host, err = ds.Host(host.ID)
	require.Nil(t, err)
	require.NotNil(t, host)

	err = ds.DeleteHost(host.ID)
	assert.Nil(t, err)

	host, err = ds.Host(host.ID)
	assert.NotNil(t, err)
	assert.Nil(t, host)
}

func TestDeleteHostWithSoftware(t *testing.T) {
	ds := CreateMySQLDS(t)
	defer ds.Close()

	host, err := ds.NewHost(&fleet.Host{
		DetailUpdatedAt: time.Now(),
		LabelUpdatedAt:  time.Now(),
		SeenTime:        time.Now(),
		NodeKey:         "1",
		UUID:            "1",
		Hostname:        "foo.local",
		PrimaryIP:       "192.168.1.1",
		PrimaryMac:      "30-65-EC-6F-C4-58",
	})
	require.NoError(t, err)
	require.NotNil(t, host)

	soft := fleet.HostSoftware{
		Modified: true,
		Software: []fleet.Software{
			{Name: "foo", Version: "0.0.1", Source: "chrome_extensions"},
			{Name: "foo", Version: "0.0.3", Source: "chrome_extensions"},
		},
	}
	host.HostSoftware = soft
	err = ds.SaveHostSoftware(host)
	require.NoError(t, err)

	err = ds.DeleteHost(host.ID)
	require.NoError(t, err)

	host, err = ds.Host(host.ID)
	assert.NotNil(t, err)
	assert.Nil(t, host)
}

func TestSaveHostPackStats(t *testing.T) {
	ds := CreateMySQLDS(t)
	defer ds.Close()

	host, err := ds.NewHost(&fleet.Host{
		DetailUpdatedAt: time.Now(),
		LabelUpdatedAt:  time.Now(),
		SeenTime:        time.Now(),
		NodeKey:         "1",
		UUID:            "1",
		Hostname:        "foo.local",
		PrimaryIP:       "192.168.1.1",
		PrimaryMac:      "30-65-EC-6F-C4-58",
	})
	require.NoError(t, err)
	require.NotNil(t, host)

	// Pack and query must exist for stats to save successfully
	pack1 := test.NewPack(t, ds, "test1")
	query1 := test.NewQuery(t, ds, "time", "select * from time", 0, true)
	squery1 := test.NewScheduledQuery(t, ds, pack1.ID, query1.ID, 30, true, true, "time-scheduled")
	stats1 := []fleet.ScheduledQueryStats{
		{
			ScheduledQueryName: squery1.Name,
			ScheduledQueryID:   squery1.ID,
			QueryName:          query1.Name,
			PackName:           pack1.Name,
			PackID:             pack1.ID,
			AverageMemory:      8000,
			Denylisted:         false,
			Executions:         164,
			Interval:           30,
			LastExecuted:       time.Unix(1620325191, 0).UTC(),
			OutputSize:         1337,
			SystemTime:         150,
			UserTime:           180,
			WallTime:           0,
		},
	}

	pack2 := test.NewPack(t, ds, "test2")
	squery2 := test.NewScheduledQuery(t, ds, pack2.ID, query1.ID, 30, true, true, "time-scheduled")
	query2 := test.NewQuery(t, ds, "processes", "select * from processes", 0, true)
	squery3 := test.NewScheduledQuery(t, ds, pack2.ID, query2.ID, 30, true, true, "processes")
	stats2 := []fleet.ScheduledQueryStats{
		{
			ScheduledQueryName: squery2.Name,
			ScheduledQueryID:   squery2.ID,
			QueryName:          query1.Name,
			PackName:           pack2.Name,
			PackID:             pack2.ID,
			AverageMemory:      431,
			Denylisted:         true,
			Executions:         1,
			Interval:           30,
			LastExecuted:       time.Unix(980943843, 0).UTC(),
			OutputSize:         134,
			SystemTime:         1656,
			UserTime:           18453,
			WallTime:           10,
		},
		{
			ScheduledQueryName: squery3.Name,
			ScheduledQueryID:   squery3.ID,
			QueryName:          query2.Name,
			PackName:           pack2.Name,
			PackID:             pack2.ID,
			AverageMemory:      8000,
			Denylisted:         false,
			Executions:         164,
			Interval:           30,
			LastExecuted:       time.Unix(1620325191, 0).UTC(),
			OutputSize:         1337,
			SystemTime:         150,
			UserTime:           180,
			WallTime:           0,
		},
	}

	host.PackStats = []fleet.PackStats{
		{
			PackName: "test1",
			// Append an additional entry to be sure that receiving stats for a
			// now-deleted query doesn't break saving. This extra entry should
			// not be returned on loading the host.
			QueryStats: append(stats1, fleet.ScheduledQueryStats{PackName: "foo", ScheduledQueryName: "bar"}),
		},
		{
			PackName:   "test2",
			QueryStats: stats2,
		},
	}

	require.NoError(t, ds.SaveHost(host))

	host, err = ds.Host(host.ID)
	require.NoError(t, err)

	require.Len(t, host.PackStats, 2)
	sort.Slice(host.PackStats, func(i, j int) bool {
		return host.PackStats[i].PackName < host.PackStats[j].PackName
	})
	assert.Equal(t, host.PackStats[0].PackName, "test1")
	assert.ElementsMatch(t, host.PackStats[0].QueryStats, stats1)
	assert.Equal(t, host.PackStats[1].PackName, "test2")
	assert.ElementsMatch(t, host.PackStats[1].QueryStats, stats2)

	// Set to nil should not overwrite
	host.PackStats = nil
	require.NoError(t, ds.SaveHost(host))
	host, err = ds.Host(host.ID)
	require.NoError(t, err)
	require.Len(t, host.PackStats, 2)
}

func TestIgnoresTeamPackStats(t *testing.T) {
	ds := CreateMySQLDS(t)
	defer ds.Close()

	host, err := ds.NewHost(&fleet.Host{
		DetailUpdatedAt: time.Now(),
		LabelUpdatedAt:  time.Now(),
		SeenTime:        time.Now(),
		NodeKey:         "1",
		UUID:            "1",
		Hostname:        "foo.local",
		PrimaryIP:       "192.168.1.1",
		PrimaryMac:      "30-65-EC-6F-C4-58",
	})
	require.NoError(t, err)
	require.NotNil(t, host)

	team, err := ds.NewTeam(&fleet.Team{Name: "team1"})
	require.NoError(t, err)
	require.NoError(t, ds.AddHostsToTeam(&team.ID, []uint{host.ID}))
	tp, err := ds.EnsureTeamPack(team.ID)
	require.NoError(t, err)

	tpQuery := test.NewQuery(t, ds, "tp-time", "select * from time", 0, true)
	tpSquery := test.NewScheduledQuery(t, ds, tp.ID, tpQuery.ID, 30, true, true, "time-scheduled")

	// Pack and query must exist for stats to save successfully
	pack1 := test.NewPack(t, ds, "test1")
	query1 := test.NewQuery(t, ds, "time", "select * from time", 0, true)
	squery1 := test.NewScheduledQuery(t, ds, pack1.ID, query1.ID, 30, true, true, "time-scheduled")
	stats1 := []fleet.ScheduledQueryStats{
		{
			ScheduledQueryName: squery1.Name,
			ScheduledQueryID:   squery1.ID,
			QueryName:          query1.Name,
			PackName:           pack1.Name,
			PackID:             pack1.ID,
			AverageMemory:      8000,
			Denylisted:         false,
			Executions:         164,
			Interval:           30,
			LastExecuted:       time.Unix(1620325191, 0).UTC(),
			OutputSize:         1337,
			SystemTime:         150,
			UserTime:           180,
			WallTime:           0,
		},
	}
	stats2 := []fleet.ScheduledQueryStats{
		{
			ScheduledQueryName: tpSquery.Name,
			ScheduledQueryID:   tpSquery.ID,
			QueryName:          tpQuery.Name,
			PackName:           tp.Name,
			PackID:             tp.ID,
			AverageMemory:      8000,
			Denylisted:         false,
			Executions:         164,
			Interval:           30,
			LastExecuted:       time.Unix(1620325191, 0).UTC(),
			OutputSize:         1337,
			SystemTime:         150,
			UserTime:           180,
			WallTime:           0,
		},
	}

	host.PackStats = []fleet.PackStats{
		{PackName: "test1", QueryStats: stats1},
		{PackName: teamScheduleName(team), QueryStats: stats2},
	}

	require.NoError(t, ds.SaveHost(host))

	host, err = ds.Host(host.ID)
	require.NoError(t, err)

	require.Len(t, host.PackStats, 1)
	assert.Equal(t, host.PackStats[0].PackName, "test1")
	assert.ElementsMatch(t, host.PackStats[0].QueryStats, stats1)
}

func TestDeleteHost(t *testing.T) {
	ds := CreateMySQLDS(t)
	defer ds.Close()

	host, err := ds.NewHost(&fleet.Host{
		DetailUpdatedAt: time.Now(),
		LabelUpdatedAt:  time.Now(),
		SeenTime:        time.Now(),
		NodeKey:         "1",
		UUID:            "1",
		Hostname:        "foo.local",
	})
	require.Nil(t, err)
	require.NotNil(t, host)

	err = ds.DeleteHost(host.ID)
	assert.Nil(t, err)

	host, err = ds.Host(host.ID)
	assert.NotNil(t, err)
}

func TestListHostsFilterAdditional(t *testing.T) {
	ds := CreateMySQLDS(t)
	defer ds.Close()

	h, err := ds.NewHost(&fleet.Host{
		DetailUpdatedAt: time.Now(),
		LabelUpdatedAt:  time.Now(),
		SeenTime:        time.Now(),
		OsqueryHostID:   "foobar",
		NodeKey:         "nodekey",
		UUID:            "uuid",
		Hostname:        "foobar.local",
	})
	require.Nil(t, err)

	filter := fleet.TeamFilter{User: test.UserAdmin}

	// Add additional
	additional := json.RawMessage(`{"field1": "v1", "field2": "v2"}`)
	h.Additional = &additional
	require.NoError(t, ds.SaveHostAdditional(h))

	hosts, err := ds.ListHosts(filter, fleet.HostListOptions{})
	require.Nil(t, err)
	assert.Nil(t, hosts[0].Additional)

	hosts, err = ds.ListHosts(filter, fleet.HostListOptions{AdditionalFilters: []string{"field1", "field2"}})
	require.Nil(t, err)
	assert.Equal(t, &additional, hosts[0].Additional)

	hosts, err = ds.ListHosts(filter, fleet.HostListOptions{AdditionalFilters: []string{"*"}})
	require.Nil(t, err)
	assert.Equal(t, &additional, hosts[0].Additional)

	additional = json.RawMessage(`{"field1": "v1", "missing": null}`)
	hosts, err = ds.ListHosts(filter, fleet.HostListOptions{AdditionalFilters: []string{"field1", "missing"}})
	require.Nil(t, err)
	assert.Equal(t, &additional, hosts[0].Additional)
}

func TestListHostsStatus(t *testing.T) {
	ds := CreateMySQLDS(t)
	defer ds.Close()

	for i := 0; i < 10; i++ {
		_, err := ds.NewHost(&fleet.Host{
			DetailUpdatedAt: time.Now(),
			LabelUpdatedAt:  time.Now(),
			SeenTime:        time.Now().Add(-time.Duration(i) * time.Minute),
			OsqueryHostID:   strconv.Itoa(i),
			NodeKey:         fmt.Sprintf("%d", i),
			UUID:            fmt.Sprintf("%d", i),
			Hostname:        fmt.Sprintf("foo.local%d", i),
		})
		assert.Nil(t, err)
		if err != nil {
			return
		}
	}

	filter := fleet.TeamFilter{User: test.UserAdmin}

	hosts, err := ds.ListHosts(filter, fleet.HostListOptions{StatusFilter: "online"})
	require.Nil(t, err)
	assert.Equal(t, 1, len(hosts))

	hosts, err = ds.ListHosts(filter, fleet.HostListOptions{StatusFilter: "offline"})
	require.Nil(t, err)
	assert.Equal(t, 9, len(hosts))

	hosts, err = ds.ListHosts(filter, fleet.HostListOptions{StatusFilter: "mia"})
	require.Nil(t, err)
	assert.Equal(t, 0, len(hosts))

	hosts, err = ds.ListHosts(filter, fleet.HostListOptions{StatusFilter: "new"})
	require.Nil(t, err)
	assert.Equal(t, 10, len(hosts))
}

func TestListHostsQuery(t *testing.T) {
	ds := CreateMySQLDS(t)
	defer ds.Close()

	hosts := []*fleet.Host{}
	for i := 0; i < 10; i++ {
		host, err := ds.NewHost(&fleet.Host{
			DetailUpdatedAt: time.Now(),
			LabelUpdatedAt:  time.Now(),
			SeenTime:        time.Now(),
			OsqueryHostID:   strconv.Itoa(i),
			NodeKey:         fmt.Sprintf("%d", i),
			UUID:            fmt.Sprintf("uuid_00%d", i),
			Hostname:        fmt.Sprintf("hostname%%00%d", i),
			HardwareSerial:  fmt.Sprintf("serial00%d", i),
		})
		require.NoError(t, err)
		host.PrimaryIP = fmt.Sprintf("192.168.1.%d", i)
		require.NoError(t, ds.SaveHost(host))
		hosts = append(hosts, host)
	}

	filter := fleet.TeamFilter{User: test.UserAdmin}

	team1, err := ds.NewTeam(&fleet.Team{Name: "team1"})
	require.NoError(t, err)
	team2, err := ds.NewTeam(&fleet.Team{Name: "team2"})
	require.NoError(t, err)

	for _, host := range hosts {
		require.NoError(t, ds.AddHostsToTeam(&team1.ID, []uint{host.ID}))
	}

	gotHosts, err := ds.ListHosts(filter, fleet.HostListOptions{})
	require.Nil(t, err)
	assert.Equal(t, len(hosts), len(gotHosts))

	gotHosts, err = ds.ListHosts(filter, fleet.HostListOptions{TeamFilter: &team1.ID})
	require.NoError(t, err)
	assert.Equal(t, len(hosts), len(gotHosts))

	gotHosts, err = ds.ListHosts(filter, fleet.HostListOptions{TeamFilter: &team2.ID})
	require.NoError(t, err)
	assert.Equal(t, 0, len(gotHosts))

	gotHosts, err = ds.ListHosts(filter, fleet.HostListOptions{TeamFilter: nil})
	require.NoError(t, err)
	assert.Equal(t, len(hosts), len(gotHosts))

	gotHosts, err = ds.ListHosts(filter, fleet.HostListOptions{ListOptions: fleet.ListOptions{MatchQuery: "00"}})
	require.Nil(t, err)
	assert.Equal(t, 10, len(gotHosts))

	gotHosts, err = ds.ListHosts(filter, fleet.HostListOptions{ListOptions: fleet.ListOptions{MatchQuery: "000"}})
	require.Nil(t, err)
	assert.Equal(t, 1, len(gotHosts))

	gotHosts, err = ds.ListHosts(filter, fleet.HostListOptions{ListOptions: fleet.ListOptions{MatchQuery: "192.168."}})
	require.Nil(t, err)
	assert.Equal(t, 10, len(gotHosts))

	gotHosts, err = ds.ListHosts(filter, fleet.HostListOptions{ListOptions: fleet.ListOptions{MatchQuery: "192.168.1.1"}})
	require.Nil(t, err)
	assert.Equal(t, 1, len(gotHosts))

	gotHosts, err = ds.ListHosts(filter, fleet.HostListOptions{ListOptions: fleet.ListOptions{MatchQuery: "hostname%00"}})
	require.Nil(t, err)
	assert.Equal(t, 10, len(gotHosts))

	gotHosts, err = ds.ListHosts(filter, fleet.HostListOptions{ListOptions: fleet.ListOptions{MatchQuery: "hostname%003"}})
	require.Nil(t, err)
	assert.Equal(t, 1, len(gotHosts))

	gotHosts, err = ds.ListHosts(filter, fleet.HostListOptions{ListOptions: fleet.ListOptions{MatchQuery: "uuid_"}})
	require.Nil(t, err)
	assert.Equal(t, 10, len(gotHosts))

	gotHosts, err = ds.ListHosts(filter, fleet.HostListOptions{ListOptions: fleet.ListOptions{MatchQuery: "uuid_006"}})
	require.Nil(t, err)
	assert.Equal(t, 1, len(gotHosts))

	gotHosts, err = ds.ListHosts(filter, fleet.HostListOptions{ListOptions: fleet.ListOptions{MatchQuery: "serial"}})
	require.Nil(t, err)
	assert.Equal(t, 10, len(gotHosts))

	gotHosts, err = ds.ListHosts(filter, fleet.HostListOptions{ListOptions: fleet.ListOptions{MatchQuery: "serial009"}})
	require.Nil(t, err)
	assert.Equal(t, 1, len(gotHosts))
}

func TestEnrollHost(t *testing.T) {
	ds := CreateMySQLDS(t)
	defer ds.Close()

	test.AddAllHostsLabel(t, ds)

	team, err := ds.NewTeam(&fleet.Team{Name: "team1"})
	require.NoError(t, err)

	filter := fleet.TeamFilter{User: test.UserAdmin}
	hosts, err := ds.ListHosts(filter, fleet.HostListOptions{})
	require.Nil(t, err)
	for _, host := range hosts {
		assert.Zero(t, host.LastEnrolledAt)
	}

	for _, tt := range enrollTests {
		h, err := ds.EnrollHost(tt.uuid, tt.nodeKey, &team.ID, 0)
		require.Nil(t, err)

		assert.Equal(t, tt.uuid, h.OsqueryHostID)
		assert.Equal(t, tt.nodeKey, h.NodeKey)

		// This host should be allowed to re-enroll immediately if cooldown is disabled
		_, err = ds.EnrollHost(tt.uuid, tt.nodeKey+"new", nil, 0)
		require.NoError(t, err)

		// This host should not be allowed to re-enroll immediately if cooldown is enabled
		_, err = ds.EnrollHost(tt.uuid, tt.nodeKey+"new", nil, 10*time.Second)
		require.Error(t, err)
	}

	hosts, err = ds.ListHosts(filter, fleet.HostListOptions{})

	require.Nil(t, err)
	for _, host := range hosts {
		assert.NotZero(t, host.LastEnrolledAt)
	}
}

func TestAuthenticateHost(t *testing.T) {
	ds := CreateMySQLDS(t)
	defer ds.Close()

	test.AddAllHostsLabel(t, ds)
	for _, tt := range enrollTests {
		h, err := ds.EnrollHost(tt.uuid, tt.nodeKey, nil, 0)
		require.Nil(t, err)

		returned, err := ds.AuthenticateHost(h.NodeKey)
		require.NoError(t, err)
		assert.Equal(t, h.NodeKey, returned.NodeKey)
	}

	_, err := ds.AuthenticateHost("7B1A9DC9-B042-489F-8D5A-EEC2412C95AA")
	assert.Error(t, err)

	_, err = ds.AuthenticateHost("")
	assert.Error(t, err)
}

func TestAuthenticateHostCaseSensitive(t *testing.T) {
	ds := CreateMySQLDS(t)
	defer ds.Close()

	test.AddAllHostsLabel(t, ds)
	for _, tt := range enrollTests {
		h, err := ds.EnrollHost(tt.uuid, tt.nodeKey, nil, 0)
		require.Nil(t, err)

		_, err = ds.AuthenticateHost(strings.ToUpper(h.NodeKey))
		require.Error(t, err, "node key authentication should be case sensitive")
	}
}

func TestSearchHosts(t *testing.T) {
	ds := CreateMySQLDS(t)
	defer ds.Close()

	_, err := ds.NewHost(&fleet.Host{
		OsqueryHostID:   "1234",
		DetailUpdatedAt: time.Now(),
		LabelUpdatedAt:  time.Now(),
		SeenTime:        time.Now(),
		NodeKey:         "1",
		UUID:            "1",
		Hostname:        "foo.local",
	})
	require.Nil(t, err)

	h2, err := ds.NewHost(&fleet.Host{
		OsqueryHostID:   "5679",
		DetailUpdatedAt: time.Now(),
		LabelUpdatedAt:  time.Now(),
		SeenTime:        time.Now(),
		NodeKey:         "2",
		UUID:            "2",
		Hostname:        "bar.local",
	})
	require.Nil(t, err)

	h3, err := ds.NewHost(&fleet.Host{
		OsqueryHostID:   "99999",
		DetailUpdatedAt: time.Now(),
		LabelUpdatedAt:  time.Now(),
		SeenTime:        time.Now(),
		NodeKey:         "3",
		UUID:            "abc-def-ghi",
		Hostname:        "foo-bar.local",
	})
	require.Nil(t, err)

	user := &fleet.User{GlobalRole: ptr.String(fleet.RoleAdmin)}
	filter := fleet.TeamFilter{User: user}

	// We once threw errors when the search query was empty. Verify that we
	// don't error.
	_, err = ds.SearchHosts(filter, "")
	require.Nil(t, err)

	hosts, err := ds.SearchHosts(filter, "foo")
	assert.Nil(t, err)
	assert.Len(t, hosts, 2)

	host, err := ds.SearchHosts(filter, "foo", h3.ID)
	require.Nil(t, err)
	require.Len(t, host, 1)
	assert.Equal(t, "foo.local", host[0].Hostname)

	host, err = ds.SearchHosts(filter, "foo", h3.ID, h2.ID)
	require.Nil(t, err)
	require.Len(t, host, 1)
	assert.Equal(t, "foo.local", host[0].Hostname)

	host, err = ds.SearchHosts(filter, "abc")
	require.Nil(t, err)
	require.Len(t, host, 1)
	assert.Equal(t, "abc-def-ghi", host[0].UUID)

	none, err := ds.SearchHosts(filter, "xxx")
	assert.Nil(t, err)
	assert.Len(t, none, 0)

	// check to make sure search on ip address works
	h2.PrimaryIP = "99.100.101.103"
	err = ds.SaveHost(h2)
	require.Nil(t, err)

	hits, err := ds.SearchHosts(filter, "99.100.101")
	require.Nil(t, err)
	require.Equal(t, 1, len(hits))

	hits, err = ds.SearchHosts(filter, "99.100.111")
	require.Nil(t, err)
	assert.Equal(t, 0, len(hits))

	h3.PrimaryIP = "99.100.101.104"
	err = ds.SaveHost(h3)
	require.Nil(t, err)
	hits, err = ds.SearchHosts(filter, "99.100.101")
	require.Nil(t, err)
	assert.Equal(t, 2, len(hits))
	hits, err = ds.SearchHosts(filter, "99.100.101", h3.ID)
	require.Nil(t, err)
	assert.Equal(t, 1, len(hits))
}

func TestSearchHostsLimit(t *testing.T) {
	ds := CreateMySQLDS(t)
	defer ds.Close()

	filter := fleet.TeamFilter{User: test.UserAdmin}

	for i := 0; i < 15; i++ {
		_, err := ds.NewHost(&fleet.Host{
			DetailUpdatedAt: time.Now(),
			LabelUpdatedAt:  time.Now(),
			SeenTime:        time.Now(),
			OsqueryHostID:   fmt.Sprintf("host%d", i),
			NodeKey:         fmt.Sprintf("%d", i),
			UUID:            fmt.Sprintf("%d", i),
			Hostname:        fmt.Sprintf("foo.%d.local", i),
		})
		require.Nil(t, err)
	}

	hosts, err := ds.SearchHosts(filter, "foo")
	require.Nil(t, err)
	assert.Len(t, hosts, 10)
}

func TestGenerateHostStatusStatistics(t *testing.T) {
	ds := CreateMySQLDS(t)
	defer ds.Close()

	filter := fleet.TeamFilter{User: test.UserAdmin}
	mockClock := clock.NewMockClock()

	online, offline, mia, new, err := ds.GenerateHostStatusStatistics(filter, mockClock.Now())
	assert.Nil(t, err)
	assert.Equal(t, uint(0), online)
	assert.Equal(t, uint(0), offline)
	assert.Equal(t, uint(0), mia)
	assert.Equal(t, uint(0), new)

	// Online
	h, err := ds.NewHost(&fleet.Host{
		ID:              1,
		OsqueryHostID:   "1",
		NodeKey:         "1",
		DetailUpdatedAt: mockClock.Now().Add(-30 * time.Second),
		LabelUpdatedAt:  mockClock.Now().Add(-30 * time.Second),
		SeenTime:        mockClock.Now().Add(-30 * time.Second),
	})
	require.Nil(t, err)
	h.DistributedInterval = 15
	h.ConfigTLSRefresh = 30
	require.Nil(t, ds.SaveHost(h))

	// Online
	h, err = ds.NewHost(&fleet.Host{
		ID:              2,
		OsqueryHostID:   "2",
		NodeKey:         "2",
		DetailUpdatedAt: mockClock.Now().Add(-1 * time.Minute),
		LabelUpdatedAt:  mockClock.Now().Add(-1 * time.Minute),
		SeenTime:        mockClock.Now().Add(-1 * time.Minute),
	})
	require.Nil(t, err)
	h.DistributedInterval = 60
	h.ConfigTLSRefresh = 3600
	require.Nil(t, ds.SaveHost(h))

	// Offline
	h, err = ds.NewHost(&fleet.Host{
		ID:              3,
		OsqueryHostID:   "3",
		NodeKey:         "3",
		DetailUpdatedAt: mockClock.Now().Add(-1 * time.Hour),
		LabelUpdatedAt:  mockClock.Now().Add(-1 * time.Hour),
		SeenTime:        mockClock.Now().Add(-1 * time.Hour),
	})
	require.Nil(t, err)
	h.DistributedInterval = 300
	h.ConfigTLSRefresh = 300
	require.Nil(t, ds.SaveHost(h))

	// MIA
	h, err = ds.NewHost(&fleet.Host{
		ID:              4,
		OsqueryHostID:   "4",
		NodeKey:         "4",
		DetailUpdatedAt: mockClock.Now().Add(-35 * (24 * time.Hour)),
		LabelUpdatedAt:  mockClock.Now().Add(-35 * (24 * time.Hour)),
		SeenTime:        mockClock.Now().Add(-35 * (24 * time.Hour)),
	})
	require.Nil(t, err)

	online, offline, mia, new, err = ds.GenerateHostStatusStatistics(filter, mockClock.Now())
	assert.Nil(t, err)
	assert.Equal(t, uint(2), online)
	assert.Equal(t, uint(1), offline)
	assert.Equal(t, uint(1), mia)
	assert.Equal(t, uint(4), new)

	online, offline, mia, new, err = ds.GenerateHostStatusStatistics(filter, mockClock.Now().Add(1*time.Hour))
	assert.Nil(t, err)
	assert.Equal(t, uint(0), online)
	assert.Equal(t, uint(3), offline)
	assert.Equal(t, uint(1), mia)
	assert.Equal(t, uint(4), new)
}

func TestMarkHostSeen(t *testing.T) {
	ds := CreateMySQLDS(t)
	defer ds.Close()

	mockClock := clock.NewMockClock()

	anHourAgo := mockClock.Now().Add(-1 * time.Hour).UTC()
	aDayAgo := mockClock.Now().Add(-24 * time.Hour).UTC()

	h1, err := ds.NewHost(&fleet.Host{
		ID:              1,
		OsqueryHostID:   "1",
		UUID:            "1",
		NodeKey:         "1",
		DetailUpdatedAt: aDayAgo,
		LabelUpdatedAt:  aDayAgo,
		SeenTime:        aDayAgo,
	})
	assert.Nil(t, err)

	{
		h1Verify, err := ds.Host(1)
		assert.Nil(t, err)
		require.NotNil(t, h1Verify)
		assert.WithinDuration(t, aDayAgo, h1Verify.SeenTime, time.Second)
	}

	err = ds.MarkHostSeen(h1, anHourAgo)
	assert.Nil(t, err)

	{
		h1Verify, err := ds.Host(1)
		assert.Nil(t, err)
		require.NotNil(t, h1Verify)
		assert.WithinDuration(t, anHourAgo, h1Verify.SeenTime, time.Second)
	}
}

func TestMarkHostsSeen(t *testing.T) {
	ds := CreateMySQLDS(t)
	defer ds.Close()

	mockClock := clock.NewMockClock()

	aSecondAgo := mockClock.Now().Add(-1 * time.Second).UTC()
	anHourAgo := mockClock.Now().Add(-1 * time.Hour).UTC()
	aDayAgo := mockClock.Now().Add(-24 * time.Hour).UTC()

	h1, err := ds.NewHost(&fleet.Host{
		ID:              1,
		OsqueryHostID:   "1",
		UUID:            "1",
		NodeKey:         "1",
		DetailUpdatedAt: aDayAgo,
		LabelUpdatedAt:  aDayAgo,
		SeenTime:        aDayAgo,
	})
	require.Nil(t, err)

	h2, err := ds.NewHost(&fleet.Host{
		ID:              2,
		OsqueryHostID:   "2",
		UUID:            "2",
		NodeKey:         "2",
		DetailUpdatedAt: aDayAgo,
		LabelUpdatedAt:  aDayAgo,
		SeenTime:        aDayAgo,
	})
	require.Nil(t, err)

	err = ds.MarkHostsSeen([]uint{h1.ID}, anHourAgo)
	assert.Nil(t, err)

	{
		h1Verify, err := ds.Host(h1.ID)
		assert.Nil(t, err)
		require.NotNil(t, h1Verify)
		assert.WithinDuration(t, anHourAgo, h1Verify.SeenTime, time.Second)

		h2Verify, err := ds.Host(h2.ID)
		assert.Nil(t, err)
		require.NotNil(t, h2Verify)
		assert.WithinDuration(t, aDayAgo, h2Verify.SeenTime, time.Second)
	}

	err = ds.MarkHostsSeen([]uint{h1.ID, h2.ID}, aSecondAgo)
	assert.Nil(t, err)

	{
		h1Verify, err := ds.Host(h1.ID)
		assert.Nil(t, err)
		require.NotNil(t, h1Verify)
		assert.WithinDuration(t, aSecondAgo, h1Verify.SeenTime, time.Second)

		h2Verify, err := ds.Host(h2.ID)
		assert.Nil(t, err)
		require.NotNil(t, h2Verify)
		assert.WithinDuration(t, aSecondAgo, h2Verify.SeenTime, time.Second)
	}

}

func TestCleanupIncomingHosts(t *testing.T) {
	ds := CreateMySQLDS(t)
	defer ds.Close()

	mockClock := clock.NewMockClock()

	h1, err := ds.NewHost(&fleet.Host{
		ID:              1,
		OsqueryHostID:   "1",
		UUID:            "1",
		NodeKey:         "1",
		DetailUpdatedAt: mockClock.Now(),
		LabelUpdatedAt:  mockClock.Now(),
		SeenTime:        mockClock.Now(),
	})
	require.Nil(t, err)

	h2, err := ds.NewHost(&fleet.Host{
		ID:              2,
		OsqueryHostID:   "2",
		UUID:            "2",
		NodeKey:         "2",
		Hostname:        "foobar",
		OsqueryVersion:  "3.2.3",
		DetailUpdatedAt: mockClock.Now(),
		LabelUpdatedAt:  mockClock.Now(),
		SeenTime:        mockClock.Now(),
	})
	require.Nil(t, err)

	err = ds.CleanupIncomingHosts(mockClock.Now().UTC())
	assert.Nil(t, err)

	// Both hosts should still exist because they are new
	_, err = ds.Host(h1.ID)
	assert.Nil(t, err)
	_, err = ds.Host(h2.ID)
	assert.Nil(t, err)

	err = ds.CleanupIncomingHosts(mockClock.Now().Add(6 * time.Minute).UTC())
	assert.Nil(t, err)

	// Now only the host with details should exist
	_, err = ds.Host(h1.ID)
	assert.NotNil(t, err)
	_, err = ds.Host(h2.ID)
	assert.Nil(t, err)
}

func TestHostIDsByName(t *testing.T) {
	ds := CreateMySQLDS(t)
	defer ds.Close()

	for i := 0; i < 10; i++ {
		_, err := ds.NewHost(&fleet.Host{
			DetailUpdatedAt: time.Now(),
			LabelUpdatedAt:  time.Now(),
			SeenTime:        time.Now(),
			OsqueryHostID:   fmt.Sprintf("host%d", i),
			NodeKey:         fmt.Sprintf("%d", i),
			UUID:            fmt.Sprintf("%d", i),
			Hostname:        fmt.Sprintf("foo.%d.local", i),
		})
		require.Nil(t, err)
	}

	filter := fleet.TeamFilter{User: test.UserAdmin}
	hosts, err := ds.HostIDsByName(filter, []string{"foo.2.local", "foo.1.local", "foo.5.local"})
	require.Nil(t, err)
	sort.Slice(hosts, func(i, j int) bool { return hosts[i] < hosts[j] })
	assert.Equal(t, hosts, []uint{2, 3, 6})
}

func TestHostAdditional(t *testing.T) {
	ds := CreateMySQLDS(t)
	defer ds.Close()

	_, err := ds.NewHost(&fleet.Host{
		DetailUpdatedAt: time.Now(),
		LabelUpdatedAt:  time.Now(),
		SeenTime:        time.Now(),
		OsqueryHostID:   "foobar",
		NodeKey:         "nodekey",
		UUID:            "uuid",
		Hostname:        "foobar.local",
	})
	require.Nil(t, err)

	h, err := ds.AuthenticateHost("nodekey")
	require.Nil(t, err)
	assert.Equal(t, "foobar.local", h.Hostname)
	assert.Nil(t, h.Additional)

	// Additional not yet set
	h, err = ds.Host(h.ID)
	require.Nil(t, err)
	assert.Nil(t, h.Additional)

	// Add additional
	additional := json.RawMessage(`{"additional": "result"}`)
	h.Additional = &additional
	require.NoError(t, ds.SaveHostAdditional(h))

	// Additional should not be loaded for authenticatehost
	h, err = ds.AuthenticateHost("nodekey")
	require.Nil(t, err)
	assert.Equal(t, "foobar.local", h.Hostname)
	assert.Nil(t, h.Additional)

	h, err = ds.Host(h.ID)
	require.Nil(t, err)
	assert.Equal(t, &additional, h.Additional)

	// Update besides additional. Additional should be unchanged.
	h, err = ds.AuthenticateHost("nodekey")
	require.Nil(t, err)
	h.Hostname = "baz.local"
	err = ds.SaveHost(h)
	require.Nil(t, err)

	h, err = ds.AuthenticateHost("nodekey")
	require.Nil(t, err)
	assert.Equal(t, "baz.local", h.Hostname)
	assert.Nil(t, h.Additional)

	h, err = ds.Host(h.ID)
	require.Nil(t, err)
	assert.Equal(t, &additional, h.Additional)

	// Update additional
	additional = json.RawMessage(`{"other": "additional"}`)
	h, err = ds.AuthenticateHost("nodekey")
	require.Nil(t, err)
	h.Additional = &additional
	err = ds.SaveHostAdditional(h)
	require.Nil(t, err)

	h, err = ds.AuthenticateHost("nodekey")
	require.Nil(t, err)
	assert.Equal(t, "baz.local", h.Hostname)
	assert.Nil(t, h.Additional)

	h, err = ds.Host(h.ID)
	require.Nil(t, err)
	assert.Equal(t, &additional, h.Additional)
}

func TestHostByIdentifier(t *testing.T) {
	ds := CreateMySQLDS(t)
	defer ds.Close()

	for i := 1; i <= 10; i++ {
		_, err := ds.NewHost(&fleet.Host{
			DetailUpdatedAt: time.Now(),
			LabelUpdatedAt:  time.Now(),
			SeenTime:        time.Now(),
			OsqueryHostID:   fmt.Sprintf("osquery_host_id_%d", i),
			NodeKey:         fmt.Sprintf("node_key_%d", i),
			UUID:            fmt.Sprintf("uuid_%d", i),
			Hostname:        fmt.Sprintf("hostname_%d", i),
		})
		require.Nil(t, err)
	}

	var (
		h   *fleet.Host
		err error
	)
	h, err = ds.HostByIdentifier("uuid_1")
	require.NoError(t, err)
	assert.Equal(t, uint(1), h.ID)

	h, err = ds.HostByIdentifier("osquery_host_id_2")
	require.NoError(t, err)
	assert.Equal(t, uint(2), h.ID)

	h, err = ds.HostByIdentifier("node_key_4")
	require.NoError(t, err)
	assert.Equal(t, uint(4), h.ID)

	h, err = ds.HostByIdentifier("hostname_7")
	require.NoError(t, err)
	assert.Equal(t, uint(7), h.ID)

	h, err = ds.HostByIdentifier("foobar")
	require.Error(t, err)
}

func TestAddHostsToTeam(t *testing.T) {
	ds := CreateMySQLDS(t)
	defer ds.Close()

	team1, err := ds.NewTeam(&fleet.Team{Name: "team1"})
	require.NoError(t, err)
	team2, err := ds.NewTeam(&fleet.Team{Name: "team2"})
	require.NoError(t, err)

	for i := 0; i < 10; i++ {
		test.NewHost(t, ds, fmt.Sprint(i), "", "key"+fmt.Sprint(i), "uuid"+fmt.Sprint(i), time.Now())
	}

	for i := 1; i <= 10; i++ {
		host, err := ds.Host(uint(i))
		require.NoError(t, err)
		assert.Nil(t, host.TeamID)
	}

	require.NoError(t, ds.AddHostsToTeam(&team1.ID, []uint{1, 2, 3}))
	require.NoError(t, ds.AddHostsToTeam(&team2.ID, []uint{3, 4, 5}))

	for i := 1; i <= 10; i++ {
		host, err := ds.Host(uint(i))
		require.NoError(t, err)
		var expectedID *uint
		switch {
		case i <= 2:
			expectedID = &team1.ID
		case i <= 5:
			expectedID = &team2.ID
		}
		assert.Equal(t, expectedID, host.TeamID)
	}

	require.NoError(t, ds.AddHostsToTeam(nil, []uint{1, 2, 3, 4}))
	require.NoError(t, ds.AddHostsToTeam(&team1.ID, []uint{5, 6, 7, 8, 9, 10}))

	for i := 1; i <= 10; i++ {
		host, err := ds.Host(uint(i))
		require.NoError(t, err)
		var expectedID *uint
		switch {
		case i >= 5:
			expectedID = &team1.ID
		}
		assert.Equal(t, expectedID, host.TeamID)
	}
}

func TestSaveUsers(t *testing.T) {
	ds := CreateMySQLDS(t)
	defer ds.Close()

	host, err := ds.NewHost(&fleet.Host{
		DetailUpdatedAt: time.Now(),
		LabelUpdatedAt:  time.Now(),
		SeenTime:        time.Now(),
		NodeKey:         "1",
		UUID:            "1",
		Hostname:        "foo.local",
		PrimaryIP:       "192.168.1.1",
		PrimaryMac:      "30-65-EC-6F-C4-58",
	})
	require.NoError(t, err)
	require.NotNil(t, host)

	err = ds.SaveHost(host)
	require.Nil(t, err)

	host, err = ds.Host(host.ID)
	require.Nil(t, err)
	assert.Len(t, host.Users, 0)

	u1 := fleet.HostUser{
		Uid:       42,
		Username:  "user",
		Type:      "aaa",
		GroupName: "group",
	}
	u2 := fleet.HostUser{
		Uid:       43,
		Username:  "user2",
		Type:      "aaa",
		GroupName: "group",
	}
	host.Users = []fleet.HostUser{u1, u2}
	host.Modified = true

	err = ds.SaveHost(host)
	require.Nil(t, err)

	host, err = ds.Host(host.ID)
	require.Nil(t, err)
	require.Len(t, host.Users, 2)
	test.ElementsMatchSkipID(t, host.Users, []fleet.HostUser{u1, u2})

	// remove u1 user
	host.Users = []fleet.HostUser{u2}
	host.Modified = true

	err = ds.SaveHost(host)
	require.Nil(t, err)

	host, err = ds.Host(host.ID)
	require.Nil(t, err)
	require.Len(t, host.Users, 1)
	assert.Equal(t, host.Users[0].Uid, u2.Uid)
}

func addHostSeenLast(t *testing.T, ds fleet.Datastore, i, days int) {
	host, err := ds.NewHost(&fleet.Host{
		DetailUpdatedAt: time.Now(),
		LabelUpdatedAt:  time.Now(),
		SeenTime:        time.Now().Add(-1 * time.Duration(days) * 24 * time.Hour),
		OsqueryHostID:   fmt.Sprintf("%d", i),
		NodeKey:         fmt.Sprintf("%d", i),
		UUID:            fmt.Sprintf("%d", i),
		Hostname:        fmt.Sprintf("foo.local%d", i),
		PrimaryIP:       fmt.Sprintf("192.168.1.%d", i),
		PrimaryMac:      fmt.Sprintf("30-65-EC-6F-C4-5%d", i),
	})
	require.NoError(t, err)
	require.NotNil(t, host)
}

func TestTotalAndUnseenHostsSince(t *testing.T) {
	ds := CreateMySQLDS(t)
	defer ds.Close()

	addHostSeenLast(t, ds, 1, 0)

	total, unseen, err := ds.TotalAndUnseenHostsSince(1)
	require.NoError(t, err)
	assert.Equal(t, 1, total)
	assert.Equal(t, 0, unseen)

	addHostSeenLast(t, ds, 2, 2)
	addHostSeenLast(t, ds, 3, 4)

	total, unseen, err = ds.TotalAndUnseenHostsSince(1)
	require.NoError(t, err)
	assert.Equal(t, 3, total)
	assert.Equal(t, 2, unseen)
}

func TestListHostsByPolicy(t *testing.T) {
	ds := CreateMySQLDS(t)
	defer ds.Close()

	for i := 0; i < 10; i++ {
		_, err := ds.NewHost(&fleet.Host{
			DetailUpdatedAt: time.Now(),
			LabelUpdatedAt:  time.Now(),
			SeenTime:        time.Now().Add(-time.Duration(i) * time.Minute),
			OsqueryHostID:   strconv.Itoa(i),
			NodeKey:         fmt.Sprintf("%d", i),
			UUID:            fmt.Sprintf("%d", i),
			Hostname:        fmt.Sprintf("foo.local%d", i),
		})
		require.NoError(t, err)
	}

	filter := fleet.TeamFilter{User: test.UserAdmin}

	q := test.NewQuery(t, ds, "query1", "select 1", 0, true)
	p, err := ds.NewGlobalPolicy(q.ID)
	require.NoError(t, err)

	// When policy response is null, we list all hosts that haven't reported at all for the policy, or errored out
	hosts, err := ds.ListHosts(filter, fleet.HostListOptions{PolicyIDFilter: &p.ID})
	require.NoError(t, err)
	require.Len(t, hosts, 10)

	h1 := hosts[0]
	h2 := hosts[1]

	hosts, err = ds.ListHosts(filter, fleet.HostListOptions{PolicyIDFilter: &p.ID, PolicyResponseFilter: ptr.Bool(true)})
	require.NoError(t, err)
	require.Len(t, hosts, 0)

	hosts, err = ds.ListHosts(filter, fleet.HostListOptions{PolicyIDFilter: &p.ID, PolicyResponseFilter: ptr.Bool(false)})
	require.NoError(t, err)
	require.Len(t, hosts, 0)

	// Make one host pass the policy and another not pass
	require.NoError(t, ds.RecordPolicyQueryExecutions(h1, map[uint]*bool{1: ptr.Bool(true)}, time.Now()))
	require.NoError(t, ds.RecordPolicyQueryExecutions(h2, map[uint]*bool{1: ptr.Bool(false)}, time.Now()))

	hosts, err = ds.ListHosts(filter, fleet.HostListOptions{PolicyIDFilter: &p.ID, PolicyResponseFilter: ptr.Bool(true)})
	require.NoError(t, err)
	require.Len(t, hosts, 1)
	assert.Equal(t, h1.ID, hosts[0].ID)

	hosts, err = ds.ListHosts(filter, fleet.HostListOptions{PolicyIDFilter: &p.ID, PolicyResponseFilter: ptr.Bool(false)})
	require.NoError(t, err)
	require.Len(t, hosts, 1)
	assert.Equal(t, h2.ID, hosts[0].ID)

	hosts, err = ds.ListHosts(filter, fleet.HostListOptions{PolicyIDFilter: &p.ID})
	require.NoError(t, err)
	require.Len(t, hosts, 8)
}

func TestSaveTonsOfUsers(t *testing.T) {
	ds := CreateMySQLDS(t)
	defer ds.Close()

	host1, err := ds.NewHost(&fleet.Host{
		DetailUpdatedAt: time.Now(),
		LabelUpdatedAt:  time.Now(),
		SeenTime:        time.Now(),
		NodeKey:         "1",
		UUID:            "1",
		Hostname:        "foo.local",
		PrimaryIP:       "192.168.1.1",
		PrimaryMac:      "30-65-EC-6F-C4-58",
		OsqueryHostID:   "1",
	})
	require.NoError(t, err)
	require.NotNil(t, host1)

	host2, err := ds.NewHost(&fleet.Host{
		DetailUpdatedAt: time.Now(),
		LabelUpdatedAt:  time.Now(),
		SeenTime:        time.Now(),
		NodeKey:         "2",
		UUID:            "2",
		Hostname:        "foo2.local",
		PrimaryIP:       "192.168.1.2",
		PrimaryMac:      "30-65-EC-6F-C4-58",
		OsqueryHostID:   "2",
	})
	require.NoError(t, err)
	require.NotNil(t, host2)

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	errCh := make(chan error)
	var count1 int32
	var count2 int32

	go func() {
		for {
			host1, err := ds.Host(host1.ID)
			if err != nil {
				errCh <- err
				return
			}

			u1 := fleet.HostUser{
				Uid:       42,
				Username:  "user",
				Type:      "aaa",
				GroupName: "group",
			}
			u2 := fleet.HostUser{
				Uid:       43,
				Username:  "user2",
				Type:      "aaa",
				GroupName: "group",
			}
			host1.Users = []fleet.HostUser{u1, u2}
			host1.SeenTime = time.Now()
			host1.Modified = true
			soft := fleet.HostSoftware{
				Modified: true,
				Software: []fleet.Software{
					{Name: "foo", Version: "0.0.1", Source: "chrome_extensions"},
					{Name: "foo", Version: "0.0.3", Source: "chrome_extensions"},
				},
			}
			host1.HostSoftware = soft
			additional := json.RawMessage(`{"some":"thing"}`)
			host1.Additional = &additional

			err = ds.SaveHost(host1)
			if err != nil {
				errCh <- err
				return
			}
			atomic.AddInt32(&count1, 1)

			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}()

	go func() {
		for {
			host2, err := ds.Host(host2.ID)
			if err != nil {
				errCh <- err
				return
			}

			u1 := fleet.HostUser{
				Uid:       99,
				Username:  "user",
				Type:      "aaa",
				GroupName: "group",
			}
			u2 := fleet.HostUser{
				Uid:       98,
				Username:  "user2",
				Type:      "aaa",
				GroupName: "group",
			}
			host2.Users = []fleet.HostUser{u1, u2}
			host2.SeenTime = time.Now()
			host2.Modified = true
			soft := fleet.HostSoftware{
				Modified: true,
				Software: []fleet.Software{
					{Name: "foo", Version: "0.0.1", Source: "chrome_extensions"},
					{Name: "foo4", Version: "0.0.3", Source: "chrome_extensions"},
				},
			}
			host2.HostSoftware = soft
			additional := json.RawMessage(`{"some":"thing"}`)
			host2.Additional = &additional

			err = ds.SaveHost(host2)
			if err != nil {
				errCh <- err
				return
			}
			atomic.AddInt32(&count2, 1)

			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}()

	ticker := time.NewTicker(10 * time.Second)

	select {
	case err := <-errCh:
		require.NoError(t, err)
		cancelFunc()
	case <-ticker.C:
	}
	fmt.Println("Count1", atomic.LoadInt32(&count1))
	fmt.Println("Count2", atomic.LoadInt32(&count2))
}
