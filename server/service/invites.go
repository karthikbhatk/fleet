package service

import (
	"context"
	"encoding/base64"
	"errors"
	"html/template"
	"strings"

	"github.com/fleetdm/fleet/v4/server"
	"github.com/fleetdm/fleet/v4/server/contexts/viewer"
	"github.com/fleetdm/fleet/v4/server/fleet"
	"github.com/fleetdm/fleet/v4/server/mail"
)

////////////////////////////////////////////////////////////////////////////////
// Create invite
////////////////////////////////////////////////////////////////////////////////

type createInviteRequest struct {
	fleet.InvitePayload
}

type createInviteResponse struct {
	Invite *fleet.Invite `json:"invite,omitempty"`
	Err    error         `json:"error,omitempty"`
}

func (r createInviteResponse) error() error { return r.Err }

func createInviteEndpoint(ctx context.Context, request interface{}, svc fleet.Service) (interface{}, error) {
	req := request.(*createInviteRequest)
	invite, err := svc.InviteNewUser(ctx, req.InvitePayload)
	if err != nil {
		return createInviteResponse{Err: err}, nil
	}
	return createInviteResponse{invite, nil}, nil
}

func (svc *Service) InviteNewUser(ctx context.Context, payload fleet.InvitePayload) (*fleet.Invite, error) {
	if err := svc.authz.Authorize(ctx, &fleet.Invite{}, fleet.ActionWrite); err != nil {
		return nil, err
	}

	if payload.Email != nil {
		*payload.Email = strings.ToLower(*payload.Email)
	}
	// TODO: looks like this would panic if no email is present in the payload.

	// verify that the user with the given email does not already exist
	_, err := svc.ds.UserByEmail(ctx, *payload.Email)
	if err == nil {
		return nil, fleet.NewInvalidArgumentError("email", "a user with this account already exists")
	}
	var nfe fleet.NotFoundError
	if !errors.As(err, &nfe) {
		return nil, err
	}

	// find the user who created the invite
	v, ok := viewer.FromContext(ctx)
	if !ok {
		return nil, errors.New("missing viewer context for create invite")
	}
	inviter := v.User

	random, err := server.GenerateRandomText(svc.config.App.TokenKeySize)
	if err != nil {
		return nil, err
	}
	token := base64.URLEncoding.EncodeToString([]byte(random))

	invite := &fleet.Invite{
		Email:      *payload.Email,
		InvitedBy:  inviter.ID,
		Token:      token,
		GlobalRole: payload.GlobalRole,
		Teams:      payload.Teams,
	}
	if payload.Position != nil {
		invite.Position = *payload.Position
	}
	if payload.Name != nil {
		invite.Name = *payload.Name
	}
	if payload.SSOEnabled != nil {
		invite.SSOEnabled = *payload.SSOEnabled
	}

	invite, err = svc.ds.NewInvite(ctx, invite)
	if err != nil {
		return nil, err
	}

	config, err := svc.AppConfig(ctx)
	if err != nil {
		return nil, err
	}

	invitedBy := inviter.Name
	if invitedBy == "" {
		invitedBy = inviter.Email
	}
	inviteEmail := fleet.Email{
		Subject: "You are Invited to Fleet",
		To:      []string{invite.Email},
		Config:  config,
		Mailer: &mail.InviteMailer{
			Invite:    invite,
			BaseURL:   template.URL(config.ServerSettings.ServerURL + svc.config.Server.URLPrefix),
			AssetURL:  getAssetURL(),
			OrgName:   config.OrgInfo.OrgName,
			InvitedBy: invitedBy,
		},
	}

	err = svc.mailService.SendEmail(inviteEmail)
	if err != nil {
		return nil, err
	}
	return invite, nil
}

////////////////////////////////////////////////////////////////////////////////
// Update invite
////////////////////////////////////////////////////////////////////////////////

type updateInviteRequest struct {
	ID uint `url:"id"`
	fleet.InvitePayload
}

type updateInviteResponse struct {
	Invite *fleet.Invite `json:"invite"`
	Err    error         `json:"error,omitempty"`
}

func (r updateInviteResponse) error() error { return r.Err }

func updateInviteEndpoint(ctx context.Context, request interface{}, svc fleet.Service) (interface{}, error) {
	req := request.(*updateInviteRequest)
	invite, err := svc.UpdateInvite(ctx, req.ID, req.InvitePayload)
	if err != nil {
		return updateInviteResponse{Err: err}, nil
	}

	return updateInviteResponse{Invite: invite}, nil
}

func (svc *Service) UpdateInvite(ctx context.Context, id uint, payload fleet.InvitePayload) (*fleet.Invite, error) {
	if err := svc.authz.Authorize(ctx, &fleet.Invite{}, fleet.ActionWrite); err != nil {
		return nil, err
	}

	if err := fleet.ValidateRole(payload.GlobalRole.Ptr(), payload.Teams); err != nil {
		return nil, err
	}

	invite, err := svc.ds.Invite(ctx, id)
	if err != nil {
		return nil, err
	}

	if payload.Email != nil {
		invite.Email = *payload.Email
	}
	if payload.Name != nil {
		invite.Name = *payload.Name
	}
	if payload.Position != nil {
		invite.Position = *payload.Position
	}
	if payload.SSOEnabled != nil {
		invite.SSOEnabled = *payload.SSOEnabled
	}
	invite.GlobalRole = payload.GlobalRole
	invite.Teams = payload.Teams

	return svc.ds.UpdateInvite(ctx, id, invite)
}
