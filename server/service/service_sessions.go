package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"time"

	"github.com/fleetdm/fleet/v4/server/contexts/ctxerr"
	"github.com/fleetdm/fleet/v4/server/contexts/logging"
	"github.com/fleetdm/fleet/v4/server/fleet"
	"github.com/fleetdm/fleet/v4/server/sso"
	"github.com/go-kit/kit/log/level"
)

// SSOSettings returns a subset of the Single Sign-On settings as configured in
// the app config. Those can be exposed e.g. via the response to an HTTP request,
// and as such should not contain sensitive information.
func (svc *Service) SSOSettings(ctx context.Context) (*fleet.SessionSSOSettings, error) {
	// skipauth: Basic SSO settings are available to unauthenticated users (so
	// that they have the necessary information to initiate SSO).
	svc.authz.SkipAuthorization(ctx)

	logging.WithLevel(ctx, level.Info)

	appConfig, err := svc.ds.AppConfig(ctx)
	if err != nil {
		return nil, ctxerr.Wrap(ctx, err, "SessionSSOSettings getting app config")
	}

	settings := &fleet.SessionSSOSettings{
		IDPName:     appConfig.SSOSettings.IDPName,
		IDPImageURL: appConfig.SSOSettings.IDPImageURL,
		SSOEnabled:  appConfig.SSOSettings.EnableSSO,
	}
	return settings, nil
}

func (svc *Service) getMetadata(config *fleet.AppConfig) (*sso.Metadata, error) {
	if config.SSOSettings.MetadataURL != "" {
		metadata, err := sso.GetMetadata(config.SSOSettings.MetadataURL)
		if err != nil {
			return nil, err
		}
		return metadata, nil
	}

	if config.SSOSettings.Metadata != "" {
		metadata, err := sso.ParseMetadata(config.SSOSettings.Metadata)
		if err != nil {
			return nil, err
		}
		return metadata, nil
	}

	return nil, fmt.Errorf("missing metadata for idp %s", config.SSOSettings.IDPName)
}

func (svc *Service) CallbackSSO(ctx context.Context, auth fleet.Auth) (*fleet.SSOSession, error) {
	// skipauth: User context does not yet exist. Unauthenticated users may
	// hit the SSO callback.
	svc.authz.SkipAuthorization(ctx)

	logging.WithLevel(ctx, level.Info)

	appConfig, err := svc.ds.AppConfig(ctx)
	if err != nil {
		return nil, ctxerr.Wrap(ctx, err, "get config for sso")
	}

	// Load the request metadata if available

	// localhost:9080/simplesaml/saml2/idp/SSOService.php?spentityid=https://localhost:8080
	var metadata *sso.Metadata
	var redirectURL string

	if appConfig.SSOSettings.EnableSSOIdPLogin && auth.RequestID() == "" {
		// Missing request ID indicates this was IdP-initiated. Only allow if
		// configured to do so.
		metadata, err = svc.getMetadata(appConfig)
		if err != nil {
			return nil, ctxerr.Wrap(ctx, err, "get sso metadata")
		}
		redirectURL = "/"
	} else {
		session, err := svc.ssoSessionStore.Get(auth.RequestID())
		if err != nil {
			return nil, ctxerr.Wrap(ctx, err, "sso request invalid")
		}
		// Remove session to so that is can't be reused before it expires.
		err = svc.ssoSessionStore.Expire(auth.RequestID())
		if err != nil {
			return nil, ctxerr.Wrap(ctx, err, "remove sso request")
		}
		if err := xml.Unmarshal([]byte(session.Metadata), &metadata); err != nil {
			return nil, ctxerr.Wrap(ctx, err, "unmarshal metadata")
		}
		redirectURL = session.OriginalURL
	}

	// Validate response
	validator, err := sso.NewValidator(*metadata, sso.WithExpectedAudience(
		appConfig.SSOSettings.EntityID,
		appConfig.ServerSettings.ServerURL,
		appConfig.ServerSettings.ServerURL+svc.config.Server.URLPrefix+"/api/v1/fleet/sso/callback", // ACS
	))
	if err != nil {
		return nil, ctxerr.Wrap(ctx, err, "create validator from metadata")
	}
	// make sure the response hasn't been tampered with
	auth, err = validator.ValidateSignature(auth)
	if err != nil {
		return nil, ctxerr.Wrap(ctx, err, "signature validation failed")
	}
	// make sure the response isn't stale
	err = validator.ValidateResponse(auth)
	if err != nil {
		return nil, ctxerr.Wrap(ctx, err, "response validation failed")
	}

	// Get and log in user
	user, err := svc.ds.UserByEmail(ctx, auth.UserID())
	if err != nil {
		return nil, ctxerr.Wrap(ctx, err, "find user in sso callback")
	}
	// if the user is not sso enabled they are not authorized
	if !user.SSOEnabled {
		return nil, ctxerr.New(ctx, "user not configured to use sso")
	}
	token, err := svc.makeSession(ctx, user.ID)
	if err != nil {
		return nil, ctxerr.Wrap(ctx, err, "make session in sso callback")
	}
	result := &fleet.SSOSession{
		Token:       token,
		RedirectURL: redirectURL,
	}
	return result, nil
}

// makeSession is a helper that creates a new session after authentication
func (svc *Service) makeSession(ctx context.Context, id uint) (string, error) {
	sessionKeySize := svc.config.Session.KeySize
	key := make([]byte, sessionKeySize)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}

	sessionKey := base64.StdEncoding.EncodeToString(key)
	session := &fleet.Session{
		UserID:     id,
		Key:        sessionKey,
		AccessedAt: time.Now().UTC(),
	}

	_, err = svc.ds.NewSession(ctx, session)
	if err != nil {
		return "", ctxerr.Wrap(ctx, err, "creating new session")
	}

	return sessionKey, nil
}

func (svc *Service) GetSessionByKey(ctx context.Context, key string) (*fleet.Session, error) {
	session, err := svc.ds.SessionByKey(ctx, key)
	if err != nil {
		return nil, err
	}

	err = svc.validateSession(ctx, session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (svc *Service) validateSession(ctx context.Context, session *fleet.Session) error {
	if session == nil {
		return fleet.NewAuthRequiredError("active session not present")
	}

	sessionDuration := svc.config.Session.Duration
	if session.APIOnly != nil && *session.APIOnly {
		sessionDuration = 0 // make API-only tokens unlimited
	}

	// duration 0 = unlimited
	if sessionDuration != 0 && time.Since(session.AccessedAt) >= sessionDuration {
		err := svc.ds.DestroySession(ctx, session)
		if err != nil {
			return ctxerr.Wrap(ctx, err, "destroying session")
		}
		return fleet.NewAuthRequiredError("expired session")
	}

	return svc.ds.MarkSessionAccessed(ctx, session)
}
