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

////////////////////////////////////////////////////////////////////////////////
// List Teams
////////////////////////////////////////////////////////////////////////////////

type listTeamsRequest struct {
	ListOptions kolide.ListOptions
}

type listTeamsResponse struct {
	Teams []kolide.Team `json:"teams"`
	Err   error         `json:"error,omitempty"`
}

func (r listTeamsResponse) error() error { return r.Err }

func makeListTeamsEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listTeamsRequest)
		teams, err := svc.ListTeams(ctx, req.ListOptions)
		if err != nil {
			return listTeamsResponse{Err: err}, nil
		}

		resp := listTeamsResponse{Teams: []kolide.Team{}}
		for _, team := range teams {
			resp.Teams = append(resp.Teams, *team)
		}
		return resp, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Delete Team
////////////////////////////////////////////////////////////////////////////////

type deleteTeamRequest struct {
	ID uint `json:"id"`
}

type deleteTeamResponse struct {
	Err error `json:"error,omitempty"`
}

func (r deleteTeamResponse) error() error { return r.Err }

func makeDeleteTeamEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteTeamRequest)
		err := svc.DeleteTeam(ctx, req.ID)
		if err != nil {
			return deleteTeamResponse{Err: err}, nil
		}
		return deleteTeamResponse{}, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// List Team Users
////////////////////////////////////////////////////////////////////////////////

type listTeamUsersRequest struct {
	TeamID      uint
	ListOptions kolide.ListOptions
}

func makeListTeamUsersEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(listTeamUsersRequest)
		users, err := svc.ListTeamUsers(ctx, req.TeamID, req.ListOptions)
		if err != nil {
			return listUsersResponse{Err: err}, nil
		}

		resp := listUsersResponse{Users: []kolide.User{}}
		for _, user := range users {
			resp.Users = append(resp.Users, *user)
		}
		return resp, nil
	}
}

////////////////////////////////////////////////////////////////////////////////
// Add / Delete Team Users
////////////////////////////////////////////////////////////////////////////////

type modifyTeamUsersRequest struct {
	TeamID uint // From request path
	// User ID and role must be specified for add users, user ID must be
	// specified for delete users.
	Users []kolide.TeamUser `json:"users"`
}

func makeAddTeamUsersEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(modifyTeamUsersRequest)
		team, err := svc.AddTeamUsers(ctx, req.TeamID, req.Users)
		if err != nil {
			return modifyTeamResponse{Err: err}, nil
		}

		return modifyTeamResponse{Team: team}, err
	}
}

func makeDeleteTeamUsersEndpoint(svc kolide.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(modifyTeamUsersRequest)
		team, err := svc.DeleteTeamUsers(ctx, req.TeamID, req.Users)
		if err != nil {
			return modifyTeamResponse{Err: err}, nil
		}

		return modifyTeamResponse{Team: team}, err
	}
}
