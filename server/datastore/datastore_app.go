package datastore

import (
	"encoding/json"
	"sort"
	"testing"

	"github.com/fleetdm/fleet/server/kolide"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testOrgInfo(t *testing.T, ds kolide.Datastore) {
	info := &kolide.AppConfig{
		OrgName:    "Kolide",
		OrgLogoURL: "localhost:8080/logo.png",
	}

	info, err := ds.NewAppConfig(info)
	assert.Nil(t, err)
	require.NotNil(t, info)

	info2, err := ds.AppConfig()
	require.Nil(t, err)
	assert.Equal(t, info2.OrgName, info.OrgName)
	assert.False(t, info2.SMTPConfigured)

	info2.OrgName = "koolide"
	info2.SMTPDomain = "foo"
	info2.SMTPConfigured = true
	info2.SMTPSenderAddress = "123"
	info2.SMTPServer = "server"
	info2.SMTPPort = 100
	info2.SMTPAuthenticationType = kolide.AuthTypeUserNamePassword
	info2.SMTPUserName = "username"
	info2.SMTPPassword = "password"
	info2.SMTPEnableTLS = false
	info2.SMTPAuthenticationMethod = kolide.AuthMethodCramMD5
	info2.SMTPVerifySSLCerts = true
	info2.SMTPEnableStartTLS = true
	info2.EnableSSO = true
	info2.EntityID = "kolide"
	info2.MetadataURL = "https://idp.com/metadata.xml"
	info2.IssuerURI = "https://idp.issuer.com"
	info2.IDPName = "My IDP"

	err = ds.SaveAppConfig(info2)
	require.Nil(t, err)

	info3, err := ds.AppConfig()
	require.Nil(t, err)
	assert.Equal(t, info2, info3)

	info4, err := ds.NewAppConfig(info3)
	assert.Nil(t, err)
	assert.Equal(t, info3, info4)
}

func testAdditionalQueries(t *testing.T, ds kolide.Datastore) {
	additional := json.RawMessage("not valid json")
	info := &kolide.AppConfig{
		OrgName:           "Kolide",
		OrgLogoURL:        "localhost:8080/logo.png",
		AdditionalQueries: &additional,
	}

	_, err := ds.NewAppConfig(info)
	assert.NotNil(t, err)

	additional = json.RawMessage(`{}`)
	info, err = ds.NewAppConfig(info)
	assert.Nil(t, err)

	additional = json.RawMessage(`{"foo": "bar"}`)
	info, err = ds.NewAppConfig(info)
	assert.Nil(t, err)
	assert.JSONEq(t, `{"foo":"bar"}`, string(*info.AdditionalQueries))
}

func testEnrollSecrets(t *testing.T, ds kolide.Datastore) {
	team1, err := ds.NewTeam(&kolide.Team{Name: "team1"})
	require.NoError(t, err)

	secret, err := ds.VerifyEnrollSecret("missing")
	assert.Error(t, err)
	assert.Nil(t, secret)

	err = ds.ApplyEnrollSecrets(&team1.ID,
		[]*kolide.EnrollSecret{
			{Secret: "one_secret", TeamID: &team1.ID},
		},
	)
	assert.NoError(t, err)

	secret, err = ds.VerifyEnrollSecret("one")
	assert.Error(t, err, "secret should not match")
	assert.Nil(t, secret, "secret should be nil")
	secret, err = ds.VerifyEnrollSecret("one_secret")
	assert.NoError(t, err)
	assert.Equal(t, &team1.ID, secret.TeamID)
	secret, err = ds.VerifyEnrollSecret("two_secret")
	assert.Error(t, err)
	assert.Nil(t, secret)

	// Add global secret
	err = ds.ApplyEnrollSecrets(
		nil,
		[]*kolide.EnrollSecret{
			{Secret: "two_secret"},
		},
	)
	assert.NoError(t, err)

	secret, err = ds.VerifyEnrollSecret("one_secret")
	assert.NoError(t, err)
	assert.Equal(t, &team1.ID, secret.TeamID)
	secret, err = ds.VerifyEnrollSecret("two_secret")
	assert.NoError(t, err)
	assert.Equal(t, (*uint)(nil), secret.TeamID)

	// Remove team secret
	err = ds.ApplyEnrollSecrets(&team1.ID, []*kolide.EnrollSecret{})
	assert.NoError(t, err)
	secret, err = ds.VerifyEnrollSecret("one_secret")
	assert.Error(t, err)
	assert.Nil(t, secret)
	secret, err = ds.VerifyEnrollSecret("two_secret")
	assert.NoError(t, err)
	assert.Equal(t, (*uint)(nil), secret.TeamID)
}

func testEnrollSecretsCaseSensitive(t *testing.T, ds kolide.Datastore) {
	err := ds.ApplyEnrollSecrets(
		nil,
		[]*kolide.EnrollSecret{
			{Secret: "one_secret"},
		},
	)
	require.NoError(t, err)

	_, err = ds.VerifyEnrollSecret("one_secret")
	assert.NoError(t, err, "enroll secret should match with matching case")
	_, err = ds.VerifyEnrollSecret("One_Secret")
	assert.Error(t, err, "enroll secret with different case should not verify")
}

func testEnrollSecretRoundtrip(t *testing.T, ds kolide.Datastore) {
	team1, err := ds.NewTeam(&kolide.Team{Name: "team1"})
	require.NoError(t, err)

	secrets, err := ds.GetEnrollSecrets(nil)
	require.NoError(t, err)
	assert.Len(t, secrets, 0)

	secrets, err = ds.GetEnrollSecrets(&team1.ID)
	require.NoError(t, err)
	assert.Len(t, secrets, 0)

	expectedSecrets := []*kolide.EnrollSecret{
		{Secret: "one_secret"},
		{Secret: "two_secret"},
	}
	err = ds.ApplyEnrollSecrets(&team1.ID, expectedSecrets)
	require.NoError(t, err)

	secrets, err = ds.GetEnrollSecrets(&team1.ID)
	require.NoError(t, err)
	require.Len(t, secrets, 2)
	// sort secrets before equality checks to ensure proper order
	sort.Slice(secrets, func(i, j int) bool { return secrets[i].Secret < secrets[j].Secret })

	assert.Equal(t, "one_secret", secrets[0].Secret)
	assert.Equal(t, "two_secret", secrets[1].Secret)

	expectedSecrets[0].Secret += "_global"
	expectedSecrets[1].Secret += "_global"
	err = ds.ApplyEnrollSecrets(nil, expectedSecrets)
	require.NoError(t, err)

	secrets, err = ds.GetEnrollSecrets(nil)
	require.NoError(t, err)
	require.Len(t, secrets, 2)

}

func testEnrollSecretUniqueness(t *testing.T, ds kolide.Datastore) {
	team1, err := ds.NewTeam(&kolide.Team{Name: "team1"})
	require.NoError(t, err)

	expectedSecrets := []*kolide.EnrollSecret{
		{Secret: "one_secret"},
	}
	err = ds.ApplyEnrollSecrets(&team1.ID, expectedSecrets)
	require.NoError(t, err)

	// Same secret at global level should not be allowed
	err = ds.ApplyEnrollSecrets(nil, expectedSecrets)
	require.Error(t, err)
}
