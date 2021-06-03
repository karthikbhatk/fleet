package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/xml"
	"net/url"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/fleetdm/fleet/server/contexts/viewer"
	"github.com/fleetdm/fleet/server/kolide"
	"github.com/fleetdm/fleet/server/sso"
	"github.com/pkg/errors"
)

func (svc *Service) SSOSettings(ctx context.Context) (*kolide.SSOSettings, error) {
	// skipauth: Basic SSO settings are available to unauthenticated users (so
	// that they have the necessary information to initiate SSO).
	svc.authz.SkipAuthorization(ctx)

	appConfig, err := svc.ds.AppConfig()
	if err != nil {
		return nil, errors.Wrap(err, "SSOSettings getting app config")
	}
	settings := &kolide.SSOSettings{
		IDPName:     appConfig.IDPName,
		IDPImageURL: appConfig.IDPImageURL,
		SSOEnabled:  appConfig.EnableSSO,
	}
	return settings, nil
}

func (svc *Service) InitiateSSO(ctx context.Context, redirectURL string) (string, error) {
	// skipauth: User context does not yet exist. Unauthenticated users may
	// initiate SSO.
	svc.authz.SkipAuthorization(ctx)

	appConfig, err := svc.ds.AppConfig()
	if err != nil {
		return "", errors.Wrap(err, "InitiateSSO getting app config")
	}

	metadata, err := svc.getMetadata(appConfig)
	if err != nil {
		return "", errors.Wrap(err, "InitiateSSO getting metadata")
	}

	settings := sso.Settings{
		Metadata: metadata,
		// Construct call back url to send to idp
		AssertionConsumerServiceURL: appConfig.KolideServerURL + svc.config.Server.URLPrefix + "/api/v1/fleet/sso/callback",
		SessionStore:                svc.ssoSessionStore,
		OriginalURL:                 redirectURL,
	}

	// If issuer is not explicitly set, default to host name.
	var issuer string
	if appConfig.EntityID == "" {
		u, err := url.Parse(appConfig.KolideServerURL)
		if err != nil {
			return "", errors.Wrap(err, "parsing kolide server url")
		}
		issuer = u.Hostname()
	} else {
		issuer = appConfig.EntityID
	}
	idpURL, err := sso.CreateAuthorizationRequest(&settings, issuer)
	if err != nil {
		return "", errors.Wrap(err, "InitiateSSO creating authorization")
	}

	return idpURL, nil
}

func (svc *Service) getMetadata(config *kolide.AppConfig) (*sso.Metadata, error) {
	if config.MetadataURL != "" {
		metadata, err := sso.GetMetadata(config.MetadataURL)
		if err != nil {
			return nil, err
		}
		return metadata, nil
	}
	if config.Metadata != "" {
		metadata, err := sso.ParseMetadata(config.Metadata)
		if err != nil {
			return nil, err
		}
		return metadata, nil
	}
	return nil, errors.Errorf("missing metadata for idp %s", config.IDPName)
}

func (svc *Service) CallbackSSO(ctx context.Context, auth kolide.Auth) (*kolide.SSOSession, error) {
	// skipauth: User context does not yet exist. Unauthenticated users may
	// hit the SSO callback.
	svc.authz.SkipAuthorization(ctx)

	appConfig, err := svc.ds.AppConfig()
	if err != nil {
		return nil, errors.Wrap(err, "get config for sso")
	}

	// Load the request metadata if available

	// localhost:9080/simplesaml/saml2/idp/SSOService.php?spentityid=https://localhost:8080
	var metadata *sso.Metadata
	var redirectURL string
	if appConfig.EnableSSOIdPLogin && auth.RequestID() == "" {
		// Missing request ID indicates this was IdP-initiated. Only allow if
		// configured to do so.
		metadata, err = svc.getMetadata(appConfig)
		if err != nil {
			return nil, errors.Wrap(err, "get sso metadata")
		}
		redirectURL = "/"
	} else {
		session, err := svc.ssoSessionStore.Get(auth.RequestID())
		if err != nil {
			return nil, errors.Wrap(err, "sso request invalid")
		}
		// Remove session to so that is can't be reused before it expires.
		err = svc.ssoSessionStore.Expire(auth.RequestID())
		if err != nil {
			return nil, errors.Wrap(err, "remove sso request")
		}
		if err := xml.Unmarshal([]byte(session.Metadata), &metadata); err != nil {
			return nil, errors.Wrap(err, "unmarshal metadata")
		}
		redirectURL = session.OriginalURL
	}

	// Validate response
	validator, err := sso.NewValidator(*metadata)
	if err != nil {
		return nil, errors.Wrap(err, "create validator from metadata")
	}
	// make sure the response hasn't been tampered with
	auth, err = validator.ValidateSignature(auth)
	if err != nil {
		return nil, errors.Wrap(err, "signature validation failed")
	}
	// make sure the response isn't stale
	err = validator.ValidateResponse(auth)
	if err != nil {
		return nil, errors.Wrap(err, "response validation failed")
	}

	// Get and log in user
	user, err := svc.userByEmailOrUsername(auth.UserID())
	if err != nil {
		return nil, errors.Wrap(err, "find user in sso callback")
	}
	// if the user is not sso enabled they are not authorized
	if !user.SSOEnabled {
		return nil, errors.New("user not configured to use sso")
	}
	token, err := svc.makeSession(user.ID)
	if err != nil {
		return nil, errors.Wrap(err, "make session in sso callback")
	}
	result := &kolide.SSOSession{
		Token:       token,
		RedirectURL: redirectURL,
	}
	return result, nil
}

func (svc *Service) Login(ctx context.Context, username, password string) (*kolide.User, string, error) {
	// skipauth: No user context available yet to authorize against.
	svc.authz.SkipAuthorization(ctx)

	// If there is an error, sleep until the request has taken at least 1
	// second. This means that generally a login failure for any reason will
	// take ~1s and frustrate a timing attack.
	var err error
	defer func(start time.Time) {
		if err != nil {
			time.Sleep(time.Until(start.Add(1 * time.Second)))
		}
	}(time.Now())

	user, err := svc.userByEmailOrUsername(username)
	if _, ok := err.(kolide.NotFoundError); ok {
		return nil, "", kolide.NewAuthFailedError("user not found")
	}
	if err != nil {
		return nil, "", kolide.NewAuthFailedError(err.Error())
	}

	if err = user.ValidatePassword(password); err != nil {
		return nil, "", kolide.NewAuthFailedError("invalid password")
	}

	if user.SSOEnabled {
		return nil, "", kolide.NewAuthFailedError("password login disabled for sso users")
	}

	token, err := svc.makeSession(user.ID)
	if err != nil {
		return nil, "", kolide.NewAuthFailedError(err.Error())
	}

	return user, token, nil
}

func (svc *Service) userByEmailOrUsername(username string) (*kolide.User, error) {
	if strings.Contains(username, "@") {
		return svc.ds.UserByEmail(username)
	}
	return svc.ds.User(username)
}

// makeSession is a helper that creates a new session after authentication
func (svc *Service) makeSession(id uint) (string, error) {
	sessionKeySize := svc.config.Session.KeySize
	key := make([]byte, sessionKeySize)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}

	session := &kolide.Session{
		UserID:     id,
		Key:        base64.StdEncoding.EncodeToString(key),
		AccessedAt: time.Now().UTC(),
	}

	session, err = svc.ds.NewSession(session)
	if err != nil {
		return "", errors.Wrap(err, "creating new session")
	}

	tokenString, err := generateJWT(session.Key, svc.config.Auth.JwtKey)
	if err != nil {
		return "", errors.Wrap(err, "generating JWT token")
	}

	return tokenString, nil
}

func (svc *Service) Logout(ctx context.Context) error {
	// skipauth: Any user can always log out of their own session.
	svc.authz.SkipAuthorization(ctx)

	// TODO: this should not return an error if the user wasn't logged in
	return svc.DestroySession(ctx)
}

func (svc *Service) DestroySession(ctx context.Context) error {
	vc, ok := viewer.FromContext(ctx)
	if !ok {
		return kolide.ErrNoContext
	}

	session, err := svc.ds.SessionByID(vc.SessionID())
	if err != nil {
		return err
	}

	if err := svc.authz.Authorize(ctx, session, "write"); err != nil {
		return err
	}

	return svc.ds.DestroySession(session)
}

func (svc *Service) GetInfoAboutSessionsForUser(ctx context.Context, id uint) ([]*kolide.Session, error) {
	if err := svc.authz.Authorize(ctx, &kolide.Session{UserID: id}, "write"); err != nil {
		return nil, err
	}

	var validatedSessions []*kolide.Session

	sessions, err := svc.ds.ListSessionsForUser(id)
	if err != nil {
		return validatedSessions, err
	}

	for _, session := range sessions {
		if svc.validateSession(session) == nil {
			validatedSessions = append(validatedSessions, session)
		}
	}

	return validatedSessions, nil
}

func (svc *Service) DeleteSessionsForUser(ctx context.Context, id uint) error {
	if err := svc.authz.Authorize(ctx, &kolide.Session{UserID: id}, "write"); err != nil {
		return err
	}

	return svc.ds.DestroyAllSessionsForUser(id)
}

func (svc *Service) GetInfoAboutSession(ctx context.Context, id uint) (*kolide.Session, error) {
	session, err := svc.ds.SessionByID(id)
	if err != nil {
		return nil, err
	}

	if err := svc.authz.Authorize(ctx, &kolide.Session{UserID: id}, "read"); err != nil {
		return nil, err
	}

	err = svc.validateSession(session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (svc *Service) GetSessionByKey(ctx context.Context, key string) (*kolide.Session, error) {
	session, err := svc.ds.SessionByKey(key)
	if err != nil {
		return nil, err
	}

	err = svc.validateSession(session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (svc *Service) DeleteSession(ctx context.Context, id uint) error {
	session, err := svc.ds.SessionByID(id)
	if err != nil {
		return err
	}

	if err := svc.authz.Authorize(ctx, session, "write"); err != nil {
		return err
	}

	return svc.ds.DestroySession(session)
}

func (svc *Service) validateSession(session *kolide.Session) error {
	if session == nil {
		return kolide.NewAuthRequiredError("active session not present")
	}

	sessionDuration := svc.config.Session.Duration
	// duration 0 = unlimited
	if sessionDuration != 0 && time.Since(session.AccessedAt) >= sessionDuration {
		err := svc.ds.DestroySession(session)
		if err != nil {
			return errors.Wrap(err, "destroying session")
		}
		return kolide.NewAuthRequiredError("expired session")
	}

	return svc.ds.MarkSessionAccessed(session)
}

// Given a session key create a JWT to be delivered to the client
func generateJWT(sessionKey, jwtKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"session_key": sessionKey,
	})

	return token.SignedString([]byte(jwtKey))
}
