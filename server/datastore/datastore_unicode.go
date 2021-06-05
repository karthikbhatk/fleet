package datastore

import (
	"testing"
	"time"

	"github.com/fleetdm/fleet/server/fleet"
	"github.com/fleetdm/fleet/server/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testUnicode(t *testing.T, ds fleet.Datastore) {
	if ds.Name() == "inmem" {
		t.Skip("inmem is being deprecated, test skipped")
	}

	l1 := fleet.LabelSpec{
		ID:    1,
		Name:  "測試",
		Query: "query foo",
	}
	err := ds.ApplyLabelSpecs([]*fleet.LabelSpec{&l1})
	require.Nil(t, err)

	label, err := ds.Label(l1.ID)
	require.Nil(t, err)
	assert.Equal(t, "測試", label.Name)

	host, err := ds.NewHost(&fleet.Host{
		HostName:         "🍌",
		DetailUpdateTime: time.Now(),
		LabelUpdateTime:  time.Now(),
		SeenTime:         time.Now(),
	})
	require.Nil(t, err)

	host, err = ds.Host(host.ID)
	require.Nil(t, err)
	assert.Equal(t, "🍌", host.HostName)

	user, err := ds.NewUser(&fleet.User{Username: "🍱", Password: []byte{}})
	require.Nil(t, err)

	user, err = ds.User(user.Username)
	require.Nil(t, err)
	assert.Equal(t, "🍱", user.Username)

	pack := test.NewPack(t, ds, "👨🏾‍🚒")

	pack, err = ds.Pack(pack.ID)
	require.Nil(t, err)
	assert.Equal(t, "👨🏾‍🚒", pack.Name)
}
