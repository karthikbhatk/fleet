package service

import (
	"context"

	"github.com/fleetdm/fleet/server/kolide"
	"github.com/go-kit/kit/endpoint"
)

////////////////////////////////////////////////////////////////////////////////
// Create Team
////////////////////////////////////////////////////////////////////////////////

type createTeamRequest struct {
	payload kolide.TeamPayload
}

type createTeamResponse struct {
	Team *kolide.Team `json:"team,omitempty"`
	Err  error        `json:"error,omitempty"`
}

func (r createTeamResponse) error() error { return r.Err }

func makeCreateTeamEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createTeamRequest)

		team, err := svc.NewTeam(ctx, req.payload)
		if err != nil {
			return createTeamResponse{Err: err}, nil
		}

		return createTeamResponse{Team: team}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Modify Team
////////////////////////////////////////////////////////////////////////////////

type modifyTeamRequest struct {
	ID      uint
	payload kolide.TeamPayload
}

type modifyTeamResponse struct {
	Team *kolide.Team `json:"team,omitempty"`
	Err  error        `json:"error,omitempty"`
}

func (r modifyTeamResponse) error() error { return r.Err }

func makeModifyTeamEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(modifyTeamRequest)
		team, err := svc.ModifyTeam(ctx, req.ID, req.payload)
		if err != nil {
			return modifyTeamResponse{Err: err}, nil
		}

		return modifyTeamResponse{Team: team}, err
	}
}
