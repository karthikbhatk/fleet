package service

import (
	"context"
	"errors"

	"github.com/fleetdm/fleet/server/contexts/viewer"
	"github.com/fleetdm/fleet/server/fleet"
	"github.com/go-kit/kit/endpoint"
	"github.com/kolide/kit/version"
)

type appConfigRequest struct {
	Payload fleet.AppConfigPayload
}

type appConfigResponse struct {
	OrgInfo            *fleet.OrgInfo             `json:"org_info,omitempty"`
	ServerSettings     *fleet.ServerSettings      `json:"server_settings,omitempty"`
	SMTPSettings       *fleet.SMTPSettingsPayload `json:"smtp_settings,omitempty"`
	SSOSettings        *fleet.SSOSettingsPayload  `json:"sso_settings,omitempty"`
	HostExpirySettings *fleet.HostExpirySettings  `json:"host_expiry_settings,omitempty"`
	HostSettings       *fleet.HostSettings        `json:"host_settings,omitempty"`
	AgentOptions       *fleet.AgentOptions        `json:"agent_options,omitempty"`
	License            *fleet.LicenseInfo         `json:"license,omitempty"`
	Err                error                       `json:"error,omitempty"`
}

func (r appConfigResponse) error() error { return r.Err }

func makeGetAppConfigEndpoint(svc fleet.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		vc, ok := viewer.FromContext(ctx)
		if !ok {
			return nil, errors.New("could not fetch user")
		}
		config, err := svc.AppConfig(ctx)
		if err != nil {
			return nil, err
		}
		license, err := svc.License(ctx)
		if err != nil {
			return nil, err
		}
		var smtpSettings *fleet.SMTPSettingsPayload
		var ssoSettings *fleet.SSOSettingsPayload
		var hostExpirySettings *fleet.HostExpirySettings
		var agentOptions *fleet.AgentOptions
		// only admin can see smtp, sso, and host expiry settings
		if vc.CanPerformAdminActions() {
			smtpSettings = smtpSettingsFromAppConfig(config)
			if smtpSettings.SMTPPassword != nil {
				*smtpSettings.SMTPPassword = "********"
			}
			ssoSettings = &fleet.SSOSettingsPayload{
				EntityID:          &config.EntityID,
				IssuerURI:         &config.IssuerURI,
				IDPImageURL:       &config.IDPImageURL,
				Metadata:          &config.Metadata,
				MetadataURL:       &config.MetadataURL,
				IDPName:           &config.IDPName,
				EnableSSO:         &config.EnableSSO,
				EnableSSOIdPLogin: &config.EnableSSOIdPLogin,
			}
			hostExpirySettings = &fleet.HostExpirySettings{
				HostExpiryEnabled: &config.HostExpiryEnabled,
				HostExpiryWindow:  &config.HostExpiryWindow,
			}
			agentOptions = &fleet.AgentOptions{
				Config: config.AgentOptions,
			}
		}
		response := appConfigResponse{
			OrgInfo: &fleet.OrgInfo{
				OrgName:    &config.OrgName,
				OrgLogoURL: &config.OrgLogoURL,
			},
			ServerSettings: &fleet.ServerSettings{
				ServerURL:   &config.ServerURL,
				LiveQueryDisabled: &config.LiveQueryDisabled,
			},
			SMTPSettings:       smtpSettings,
			SSOSettings:        ssoSettings,
			HostExpirySettings: hostExpirySettings,
			HostSettings: &fleet.HostSettings{
				AdditionalQueries: config.AdditionalQueries,
			},
			License: license,
			AgentOptions: agentOptions,
		}
		return response, nil
	}
}

func makeModifyAppConfigEndpoint(svc fleet.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(appConfigRequest)
		config, err := svc.ModifyAppConfig(ctx, req.Payload)
		if err != nil {
			return appConfigResponse{Err: err}, nil
		}
		response := appConfigResponse{
			OrgInfo: &fleet.OrgInfo{
				OrgName:    &config.OrgName,
				OrgLogoURL: &config.OrgLogoURL,
			},
			ServerSettings: &fleet.ServerSettings{
				ServerURL:   &config.ServerURL,
				LiveQueryDisabled: &config.LiveQueryDisabled,
			},
			SMTPSettings: smtpSettingsFromAppConfig(config),
			SSOSettings: &fleet.SSOSettingsPayload{
				EntityID:          &config.EntityID,
				IssuerURI:         &config.IssuerURI,
				IDPImageURL:       &config.IDPImageURL,
				Metadata:          &config.Metadata,
				MetadataURL:       &config.MetadataURL,
				IDPName:           &config.IDPName,
				EnableSSO:         &config.EnableSSO,
				EnableSSOIdPLogin: &config.EnableSSOIdPLogin,
			},
			HostExpirySettings: &fleet.HostExpirySettings{
				HostExpiryEnabled: &config.HostExpiryEnabled,
				HostExpiryWindow:  &config.HostExpiryWindow,
			},
			AgentOptions: &fleet.AgentOptions{
				Config: config.AgentOptions,
			},
		}
		if response.SMTPSettings.SMTPPassword != nil {
			*response.SMTPSettings.SMTPPassword = "********"
		}
		return response, nil
	}
}

func smtpSettingsFromAppConfig(config *fleet.AppConfig) *fleet.SMTPSettingsPayload {
	authType := config.SMTPAuthenticationType.String()
	authMethod := config.SMTPAuthenticationMethod.String()
	return &fleet.SMTPSettingsPayload{
		SMTPEnabled:              &config.SMTPConfigured,
		SMTPConfigured:           &config.SMTPConfigured,
		SMTPSenderAddress:        &config.SMTPSenderAddress,
		SMTPServer:               &config.SMTPServer,
		SMTPPort:                 &config.SMTPPort,
		SMTPAuthenticationType:   &authType,
		SMTPUserName:             &config.SMTPUserName,
		SMTPPassword:             &config.SMTPPassword,
		SMTPEnableTLS:            &config.SMTPEnableTLS,
		SMTPAuthenticationMethod: &authMethod,
		SMTPDomain:               &config.SMTPDomain,
		SMTPVerifySSLCerts:       &config.SMTPVerifySSLCerts,
		SMTPEnableStartTLS:       &config.SMTPEnableStartTLS,
	}
}

////////////////////////////////////////////////////////////////////////////////
// Apply Enroll Secret Spec
////////////////////////////////////////////////////////////////////////////////

type applyEnrollSecretSpecRequest struct {
	Spec *fleet.EnrollSecretSpec `json:"spec"`
}

type applyEnrollSecretSpecResponse struct {
	Err error `json:"error,omitempty"`
}

func (r applyEnrollSecretSpecResponse) error() error { return r.Err }

func makeApplyEnrollSecretSpecEndpoint(svc fleet.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(applyEnrollSecretSpecRequest)
		err := svc.ApplyEnrollSecretSpec(ctx, req.Spec)
		if err != nil {
			return applyEnrollSecretSpecResponse{Err: err}, nil
		}
		return applyEnrollSecretSpecResponse{}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Get Enroll Secret Spec
////////////////////////////////////////////////////////////////////////////////

type getEnrollSecretSpecResponse struct {
	Spec *fleet.EnrollSecretSpec `json:"specs"`
	Err  error                    `json:"error,omitempty"`
}

func (r getEnrollSecretSpecResponse) error() error { return r.Err }

func makeGetEnrollSecretSpecEndpoint(svc fleet.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		specs, err := svc.GetEnrollSecretSpec(ctx)
		if err != nil {
			return getEnrollSecretSpecResponse{Err: err}, nil
		}
		return getEnrollSecretSpecResponse{Spec: specs}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Version
////////////////////////////////////////////////////////////////////////////////

type versionResponse struct {
	*version.Info
	Err error `json:"error,omitempty"`
}

func (r versionResponse) error() error { return r.Err }

func makeVersionEndpoint(svc fleet.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		info, err := svc.Version(ctx)
		if err != nil {
			return versionResponse{Err: err}, nil
		}
		return versionResponse{Info: info}, nil
	}
}
