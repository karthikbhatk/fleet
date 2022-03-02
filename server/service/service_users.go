package service

import (
	"context"
	"encoding/base64"
	"html/template"
	"time"

	"github.com/fleetdm/fleet/v4/server"

	"github.com/fleetdm/fleet/v4/server/contexts/ctxerr"
	"github.com/fleetdm/fleet/v4/server/fleet"
	"github.com/fleetdm/fleet/v4/server/mail"
	"github.com/fleetdm/fleet/v4/server/ptr"
)

func (svc *Service) CreateUserFromInvite(ctx context.Context, p fleet.UserPayload) (*fleet.User, error) {
	// skipauth: There is no viewer context at this point. We rely on verifying
	// the invite for authNZ.
	svc.authz.SkipAuthorization(ctx)

	if err := p.VerifyInviteCreate(); err != nil {
		return nil, ctxerr.Wrap(ctx, err, "verify user payload")
	}

	invite, err := svc.VerifyInvite(ctx, *p.InviteToken)
	if err != nil {
		return nil, err
	}

	// set the payload role property based on an existing invite.
	p.GlobalRole = invite.GlobalRole.Ptr()
	p.Teams = &invite.Teams

	user, err := svc.newUser(ctx, p)
	if err != nil {
		return nil, err
	}

	err = svc.ds.DeleteInvite(ctx, invite.ID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (svc *Service) CreateInitialUser(ctx context.Context, p fleet.UserPayload) (*fleet.User, error) {
	// skipauth: Only the initial user creation should be allowed to skip
	// authorization (because there is not yet a user context to check against).
	svc.authz.SkipAuthorization(ctx)

	setupRequired, err := svc.SetupRequired(ctx)
	if err != nil {
		return nil, err
	}
	if !setupRequired {
		return nil, ctxerr.New(ctx, "a user already exists")
	}

	// Initial user should be global admin with no explicit teams
	p.GlobalRole = ptr.String(fleet.RoleAdmin)
	p.Teams = nil

	return svc.newUser(ctx, p)
}

func (svc *Service) newUser(ctx context.Context, p fleet.UserPayload) (*fleet.User, error) {
	var ssoEnabled bool
	// if user is SSO generate a fake password
	if (p.SSOInvite != nil && *p.SSOInvite) || (p.SSOEnabled != nil && *p.SSOEnabled) {
		fakePassword, err := server.GenerateRandomText(14)
		if err != nil {
			return nil, ctxerr.Wrap(ctx, err, "generate stand-in password")
		}
		p.Password = &fakePassword
		ssoEnabled = true
	}
	user, err := p.User(svc.config.Auth.SaltKeySize, svc.config.Auth.BcryptCost)
	if err != nil {
		return nil, err
	}
	user.SSOEnabled = ssoEnabled
	user, err = svc.ds.NewUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (svc *Service) UserUnauthorized(ctx context.Context, id uint) (*fleet.User, error) {
	// Explicitly no authorization check. Should only be used by middleware.
	return svc.ds.UserByID(ctx, id)
}

func (svc *Service) ResetPassword(ctx context.Context, token, password string) error {
	// skipauth: No viewer context available. The user is locked out of their
	// account and authNZ is performed entirely by providing a valid password
	// reset token.
	svc.authz.SkipAuthorization(ctx)

	if token == "" {
		return ctxerr.Wrap(ctx, fleet.NewInvalidArgumentError("token", "Token cannot be empty field"))
	}
	if password == "" {
		return ctxerr.Wrap(ctx, fleet.NewInvalidArgumentError("new_password", "New password cannot be empty field"))
	}
	if err := fleet.ValidatePasswordRequirements(password); err != nil {
		return ctxerr.Wrap(ctx, fleet.NewInvalidArgumentError("new_password", err.Error()))
	}

	reset, err := svc.ds.FindPassswordResetByToken(ctx, token)
	if err != nil {
		return ctxerr.Wrap(ctx, err, "looking up reset by token")
	}
	user, err := svc.ds.UserByID(ctx, reset.UserID)
	if err != nil {
		return ctxerr.Wrap(ctx, err, "retrieving user")
	}

	if user.SSOEnabled {
		return ctxerr.New(ctx, "password reset for single sign on user not allowed")
	}

	// prevent setting the same password
	if err := user.ValidatePassword(password); err == nil {
		return fleet.NewInvalidArgumentError("new_password", "cannot reuse old password")
	}

	err = svc.setNewPassword(ctx, user, password)
	if err != nil {
		return ctxerr.Wrap(ctx, err, "setting new password")
	}

	// delete password reset tokens for user
	if err := svc.ds.DeletePasswordResetRequestsForUser(ctx, user.ID); err != nil {
		return ctxerr.Wrap(ctx, err, "delete password reset requests")
	}

	// Clear sessions so that any other browsers will have to log in with
	// the new password
	if err := svc.ds.DestroyAllSessionsForUser(ctx, user.ID); err != nil {
		return ctxerr.Wrap(ctx, err, "delete user sessions")
	}

	return nil
}

func (svc *Service) RequestPasswordReset(ctx context.Context, email string) error {
	// skipauth: No viewer context available. The user is locked out of their
	// account and trying to reset their password.
	svc.authz.SkipAuthorization(ctx)

	// Regardless of error, sleep until the request has taken at least 1 second.
	// This means that any request to this method will take ~1s and frustrate a timing attack.
	defer func(start time.Time) {
		time.Sleep(time.Until(start.Add(1 * time.Second)))
	}(time.Now())

	user, err := svc.ds.UserByEmail(ctx, email)
	if err != nil {
		return err
	}
	if user.SSOEnabled {
		return ctxerr.New(ctx, "password reset for single sign on user not allowed")
	}

	random, err := server.GenerateRandomText(svc.config.App.TokenKeySize)
	if err != nil {
		return err
	}
	token := base64.URLEncoding.EncodeToString([]byte(random))

	request := &fleet.PasswordResetRequest{
		ExpiresAt: time.Now().Add(time.Hour * 24),
		UserID:    user.ID,
		Token:     token,
	}
	_, err = svc.ds.NewPasswordResetRequest(ctx, request)
	if err != nil {
		return err
	}

	config, err := svc.ds.AppConfig(ctx)
	if err != nil {
		return err
	}

	resetEmail := fleet.Email{
		Subject: "Reset Your Fleet Password",
		To:      []string{user.Email},
		Config:  config,
		Mailer: &mail.PasswordResetMailer{
			BaseURL:  template.URL(config.ServerSettings.ServerURL + svc.config.Server.URLPrefix),
			AssetURL: getAssetURL(),
			Token:    token,
		},
	}

	return svc.mailService.SendEmail(resetEmail)
}
