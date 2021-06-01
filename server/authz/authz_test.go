package authz

import (
	"encoding/json"
	"testing"

	"github.com/fleetdm/fleet/server/kolide"
	"github.com/fleetdm/fleet/server/ptr"
	"github.com/fleetdm/fleet/server/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	write = kolide.ActionWrite
	read  = kolide.ActionRead
)

var auth *Authorizer

func init() {
	var err error
	auth, err = NewAuthorizer()
	if err != nil {
		panic(err)
	}
}

type authTestCase struct {
	user   *kolide.User
	object interface{}
	action interface{}
	allow  bool
}

func TestAuthorizeAppConfig(t *testing.T) {
	t.Parallel()

	obj := &kolide.AppConfig{}
	runTestCases(t, []authTestCase{
		{user: nil, object: obj, action: read, allow: false},
		{user: nil, object: obj, action: write, allow: false},

		{user: test.UserNoRoles, object: obj, action: read, allow: true},
		{user: test.UserNoRoles, object: obj, action: write, allow: false},

		{user: test.UserAdmin, object: obj, action: read, allow: true},
		{user: test.UserAdmin, object: obj, action: write, allow: true},

		{user: test.UserMaintainer, object: obj, action: read, allow: true},
		{user: test.UserMaintainer, object: obj, action: write, allow: false},

		{user: test.UserObserver, object: obj, action: read, allow: true},
		{user: test.UserObserver, object: obj, action: write, allow: false},
	})
}

func TestAuthorizeSession(t *testing.T) {
	t.Parallel()

	obj := &kolide.Session{UserID: 42}
	runTestCases(t, []authTestCase{
		{user: nil, object: obj, action: read, allow: false},
		{user: nil, object: obj, action: write, allow: false},

		// Admin can read/write all
		{user: test.UserAdmin, object: obj, action: read, allow: true},
		{user: test.UserAdmin, object: obj, action: write, allow: true},

		// Regular users can read self
		{user: test.UserMaintainer, object: obj, action: read, allow: false},
		{user: test.UserMaintainer, object: obj, action: write, allow: false},
		{user: test.UserMaintainer, object: &kolide.Session{UserID: test.UserMaintainer.ID}, action: read, allow: true},
		{user: test.UserMaintainer, object: &kolide.Session{UserID: test.UserMaintainer.ID}, action: write, allow: true},

		{user: test.UserNoRoles, object: obj, action: read, allow: false},
		{user: test.UserNoRoles, object: obj, action: write, allow: false},
		{user: test.UserNoRoles, object: &kolide.Session{UserID: test.UserNoRoles.ID}, action: read, allow: true},
		{user: test.UserNoRoles, object: &kolide.Session{UserID: test.UserNoRoles.ID}, action: write, allow: true},

		{user: test.UserObserver, object: obj, action: read, allow: false},
		{user: test.UserObserver, object: obj, action: write, allow: false},
		{user: test.UserObserver, object: &kolide.Session{UserID: test.UserObserver.ID}, action: read, allow: true},
		{user: test.UserObserver, object: &kolide.Session{UserID: test.UserObserver.ID}, action: write, allow: true},
	})
}

func TestAuthorizeUser(t *testing.T) {
	t.Parallel()

	obj := &kolide.User{ID: 42}
	runTestCases(t, []authTestCase{
		{user: nil, object: obj, action: read, allow: false},
		{user: nil, object: obj, action: write, allow: false},

		// Admin can read/write all
		{user: test.UserAdmin, object: obj, action: read, allow: true},
		{user: test.UserAdmin, object: obj, action: write, allow: true},

		// Regular users can read all users and write self
		{user: test.UserMaintainer, object: obj, action: read, allow: true},
		{user: test.UserMaintainer, object: obj, action: write, allow: false},
		{user: test.UserMaintainer, object: test.UserMaintainer, action: read, allow: true},
		{user: test.UserMaintainer, object: test.UserMaintainer, action: write, allow: true},

		{user: test.UserNoRoles, object: obj, action: read, allow: true},
		{user: test.UserNoRoles, object: obj, action: write, allow: false},
		{user: test.UserNoRoles, object: test.UserNoRoles, action: read, allow: true},
		{user: test.UserNoRoles, object: test.UserNoRoles, action: write, allow: true},

		{user: test.UserObserver, object: obj, action: read, allow: true},
		{user: test.UserObserver, object: obj, action: write, allow: false},
		{user: test.UserObserver, object: test.UserObserver, action: read, allow: true},
		{user: test.UserObserver, object: test.UserObserver, action: write, allow: true},
	})
}

func assertAuthorized(t *testing.T, user *kolide.User, object, action interface{}) {
	t.Helper()

	assert.NoError(t, auth.Authorize(test.UserContext(user), object, action), "should be authorized\n%v\n%v\n%v", user, object, action)
}

func assertUnauthorized(t *testing.T, user *kolide.User, object, action interface{}) {
	t.Helper()

	assert.Error(t, auth.Authorize(test.UserContext(user), object, action), "should be unauthorized\n%v\n%v\n%v", user, object, action)
}

func runTestCases(t *testing.T, testCases []authTestCase) {
	t.Helper()

	for _, tt := range testCases {
		tt := tt
		t.Run("", func(t *testing.T) {
			t.Parallel()
			if tt.allow {
				assertAuthorized(t, tt.user, tt.object, tt.action)
			} else {
				assertUnauthorized(t, tt.user, tt.object, tt.action)
			}
		})
	}

}

func TestJSONToInterfaceUser(t *testing.T) {
	t.Parallel()

	subject, err := jsonToInterface(&kolide.User{GlobalRole: ptr.String(kolide.RoleAdmin)})
	require.NoError(t, err)
	{
		subject := subject.(map[string]interface{})
		assert.Equal(t, kolide.RoleAdmin, subject["global_role"])
		assert.Nil(t, subject["teams"])
	}

	subject, err = jsonToInterface(&kolide.User{
		Teams: []kolide.UserTeam{
			{Team: kolide.Team{ID: 3}, Role: kolide.RoleObserver},
			{Team: kolide.Team{ID: 42}, Role: kolide.RoleMaintainer},
		},
	})
	require.NoError(t, err)
	{
		subject := subject.(map[string]interface{})
		assert.Equal(t, nil, subject["global_role"])
		assert.Len(t, subject["teams"], 2)
		assert.Equal(t, kolide.RoleObserver, subject["teams"].([]interface{})[0].(map[string]interface{})["role"])
		assert.Equal(t, json.Number("3"), subject["teams"].([]interface{})[0].(map[string]interface{})["id"])
		assert.Equal(t, kolide.RoleMaintainer, subject["teams"].([]interface{})[1].(map[string]interface{})["role"])
		assert.Equal(t, json.Number("42"), subject["teams"].([]interface{})[1].(map[string]interface{})["id"])
	}
}
