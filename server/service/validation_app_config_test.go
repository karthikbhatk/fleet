package service

import (
	"testing"

	"github.com/fleetdm/fleet/server/kolide"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSSONotPresent(t *testing.T) {
	invalid := &kolide.InvalidArgumentError{}
	var p kolide.AppConfigPayload
	validateSSOSettings(p, &kolide.AppConfig{}, invalid)
	assert.False(t, invalid.HasErrors())

}

func TestNeedFieldsPresent(t *testing.T) {
	invalid := &kolide.InvalidArgumentError{}
	config := kolide.AppConfig{
		EnableSSO:   true,
		EntityID:    "kolide",
		IssuerURI:   "http://issuer.idp.com",
		MetadataURL: "http://isser.metadata.com",
		IDPName:     "onelogin",
	}
	p := appConfigPayloadFromAppConfig(&config)
	validateSSOSettings(*p, &kolide.AppConfig{}, invalid)
	assert.False(t, invalid.HasErrors())
}

func TestMissingMetadata(t *testing.T) {
	invalid := &kolide.InvalidArgumentError{}
	config := kolide.AppConfig{
		EnableSSO: true,
		EntityID:  "kolide",
		IssuerURI: "http://issuer.idp.com",
		IDPName:   "onelogin",
	}
	p := appConfigPayloadFromAppConfig(&config)
	validateSSOSettings(*p, &kolide.AppConfig{}, invalid)
	require.True(t, invalid.HasErrors())
	assert.Contains(t, invalid.Error(), "metadata")
	assert.Contains(t, invalid.Error(), "either metadata or metadata_url must be defined")
}

func appConfigPayloadFromAppConfig(config *kolide.AppConfig) *kolide.AppConfigPayload {
	return &kolide.AppConfigPayload{
		OrgInfo: &kolide.OrgInfo{
			OrgLogoURL: &config.OrgLogoURL,
			OrgName:    &config.OrgName,
		},
		ServerSettings: &kolide.ServerSettings{
			KolideServerURL:   &config.KolideServerURL,
			LiveQueryDisabled: &config.LiveQueryDisabled,
		},
		SMTPSettings: smtpSettingsFromAppConfig(config),
		SSOSettings: &kolide.SSOSettingsPayload{
			EnableSSO:   &config.EnableSSO,
			IDPName:     &config.IDPName,
			Metadata:    &config.Metadata,
			MetadataURL: &config.MetadataURL,
			IssuerURI:   &config.IssuerURI,
			EntityID:    &config.EntityID,
		},
		HostExpirySettings: &kolide.HostExpirySettings{
			HostExpiryEnabled: &config.HostExpiryEnabled,
			HostExpiryWindow:  &config.HostExpiryWindow,
		},
	}
}
