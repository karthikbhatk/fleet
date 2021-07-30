package service

import (
	"context"
	"net/http"
	"time"

	"github.com/fleetdm/fleet/v4/server/fleet"
	"github.com/fleetdm/fleet/v4/server/ptr"
	kithttp "github.com/go-kit/kit/transport/http"
	"gopkg.in/guregu/null.v3"
)

type getTeamScheduleRequest struct {
	TeamID      uint              `url:"team_id"`
	ListOptions fleet.ListOptions `url:"list_options"`
}

type getTeamScheduleResponse struct {
	Scheduled []scheduledQueryResponse `json:"scheduled"`
	Err       error                    `json:"error,omitempty"`
}

func (r getTeamScheduleResponse) error() error { return r.Err }

func makeGetTeamScheduleEndpoint(svc fleet.Service, opts []kithttp.ServerOption) http.Handler {
	return newServer(
		makeAuthenticatedServiceEndpoint(svc, getTeamScheduleEndpoint),
		makeDecoderForOptionsAndIDs(getTeamScheduleRequest{}, "team_id"),
		opts,
	)
}

func getTeamScheduleEndpoint(ctx context.Context, request interface{}, svc fleet.Service) (interface{}, error) {
	req := request.(*getTeamScheduleRequest)
	resp := getTeamScheduleResponse{Scheduled: []scheduledQueryResponse{}}
	queries, err := svc.GetTeamScheduledQueries(ctx, req.TeamID, req.ListOptions)
	if err != nil {
		return getTeamScheduleResponse{Err: err}, nil
	}
	for _, q := range queries {
		resp.Scheduled = append(resp.Scheduled, scheduledQueryResponse{
			ScheduledQuery: *q,
		})
	}
	return resp, nil
}

func (svc Service) GetTeamScheduledQueries(ctx context.Context, teamID uint, opts fleet.ListOptions) ([]*fleet.ScheduledQuery, error) {
	if err := svc.authz.Authorize(ctx, &fleet.Pack{}, fleet.ActionRead); err != nil {
		return nil, err
	}

	gp, err := svc.ds.EnsureTeamPack(teamID)
	if err != nil {
		return nil, err
	}

	return svc.ds.ListScheduledQueriesInPack(gp.ID, opts)
}

func (mw loggingMiddleware) GetTeamScheduledQueries(ctx context.Context, teamID uint, opts fleet.ListOptions) ([]*fleet.ScheduledQuery, error) {
	var err error
	defer func(begin time.Time) {
		_ = mw.loggerDebug(err).Log("method", "GetTeamScheduledQueries", "err", err, "took", time.Since(begin))
	}(time.Now())
	return mw.Service.GetTeamScheduledQueries(ctx, teamID, opts)
}

/////////////////////////////////////////////////////////////////////////////////
// Add
/////////////////////////////////////////////////////////////////////////////////

type teamScheduleQueryRequest struct {
	TeamID uint `url:"team_id"`
	fleet.ScheduledQueryPayload
}

type teamScheduleQueryResponse struct {
	Scheduled *fleet.ScheduledQuery `json:"scheduled,omitempty"`
	Err       error                 `json:"error,omitempty"`
}

func (r teamScheduleQueryResponse) error() error { return r.Err }

func makeTeamScheduleQueryEndpoint(svc fleet.Service, opts []kithttp.ServerOption) http.Handler {
	return newServer(
		makeAuthenticatedServiceEndpoint(svc, teamScheduleQueryEndpoint),
		makeDecoderForTypeAndIDs(teamScheduleQueryRequest{}, "team_id"),
		opts,
	)
}

func uintValueOrZero(v *uint) uint {
	if v == nil {
		return 0
	}
	return *v
}

func nullIntToPtrUint(v *null.Int) *uint {
	if v == nil {
		return nil
	}
	return ptr.Uint(uint(v.ValueOrZero()))
}

func teamScheduleQueryEndpoint(ctx context.Context, request interface{}, svc fleet.Service) (interface{}, error) {
	req := request.(*teamScheduleQueryRequest)
	resp, err := svc.TeamScheduleQuery(ctx, req.TeamID, &fleet.ScheduledQuery{
		QueryID:  uintValueOrZero(req.QueryID),
		Interval: uintValueOrZero(req.Interval),
		Snapshot: req.Snapshot,
		Removed:  req.Removed,
		Platform: req.Platform,
		Version:  req.Version,
		Shard:    nullIntToPtrUint(req.Shard),
	})
	if err != nil {
		return teamScheduleQueryResponse{Err: err}, nil
	}
	_ = resp
	return teamScheduleQueryResponse{}, nil
}

func (svc Service) TeamScheduleQuery(ctx context.Context, teamID uint, q *fleet.ScheduledQuery) (*fleet.ScheduledQuery, error) {
	if err := svc.authz.Authorize(ctx, &fleet.Pack{}, fleet.ActionRead); err != nil {
		return nil, err
	}

	gp, err := svc.ds.EnsureTeamPack(teamID)
	if err != nil {
		return nil, err
	}
	q.PackID = gp.ID

	return svc.ScheduleQuery(ctx, q)
}

func (mw loggingMiddleware) TeamScheduleQuery(ctx context.Context, teamID uint, q *fleet.ScheduledQuery) (*fleet.ScheduledQuery, error) {
	var err error
	defer func(begin time.Time) {
		_ = mw.loggerDebug(err).Log("method", "TeamScheduleQuery", "err", err, "took", time.Since(begin))
	}(time.Now())
	return mw.Service.TeamScheduleQuery(ctx, teamID, q)
}

/////////////////////////////////////////////////////////////////////////////////
// Modify
/////////////////////////////////////////////////////////////////////////////////

type modifyTeamScheduleRequest struct {
	TeamID           uint `url:"team_id"`
	ScheduledQueryID uint `url:"scheduled_query_id"`
	fleet.ScheduledQueryPayload
}

type modifyTeamScheduleResponse struct {
	Scheduled *fleet.ScheduledQuery `json:"scheduled,omitempty"`
	Err       error                 `json:"error,omitempty"`
}

func (r modifyTeamScheduleResponse) error() error { return r.Err }

func makeModifyTeamScheduleEndpoint(svc fleet.Service, opts []kithttp.ServerOption) http.Handler {
	return newServer(
		makeAuthenticatedServiceEndpoint(svc, modifyTeamScheduleEndpoint),
		makeDecoderForTypeAndIDs(modifyTeamScheduleRequest{}, "team_id", "scheduled_query_id"),
		opts,
	)
}

func modifyTeamScheduleEndpoint(ctx context.Context, request interface{}, svc fleet.Service) (interface{}, error) {
	req := request.(*modifyTeamScheduleRequest)
	resp, err := svc.ModifyTeamScheduledQueries(ctx, req.TeamID, req.ScheduledQueryID, req.ScheduledQueryPayload)
	if err != nil {
		return modifyTeamScheduleResponse{Err: err}, nil
	}
	_ = resp
	return modifyTeamScheduleResponse{}, nil
}

func (svc Service) ModifyTeamScheduledQueries(ctx context.Context, teamID uint, scheduledQueryID uint, query fleet.ScheduledQueryPayload) (*fleet.ScheduledQuery, error) {
	if err := svc.authz.Authorize(ctx, &fleet.Pack{}, fleet.ActionWrite); err != nil {
		return nil, err
	}

	gp, err := svc.ds.EnsureTeamPack(teamID)
	if err != nil {
		return nil, err
	}

	query.PackID = ptr.Uint(gp.ID)

	return svc.ModifyScheduledQuery(ctx, scheduledQueryID, query)
}

func (mw loggingMiddleware) ModifyTeamScheduledQueries(ctx context.Context, teamID uint, scheduledQueryID uint, q fleet.ScheduledQueryPayload) (*fleet.ScheduledQuery, error) {
	var err error
	defer func(begin time.Time) {
		_ = mw.loggerDebug(err).Log("method", "ModifyTeamScheduledQueries", "err", err, "took", time.Since(begin))
	}(time.Now())
	return mw.Service.ModifyTeamScheduledQueries(ctx, teamID, scheduledQueryID, q)
}

/////////////////////////////////////////////////////////////////////////////////
// Delete
/////////////////////////////////////////////////////////////////////////////////

type deleteTeamScheduleRequest struct {
	TeamID           uint `url:"team_id"`
	ScheduledQueryID uint `url:"scheduled_query_id"`
}

type deleteTeamScheduleResponse struct {
	Scheduled *fleet.ScheduledQuery `json:"scheduled,omitempty"`
	Err       error                 `json:"error,omitempty"`
}

func (r deleteTeamScheduleResponse) error() error { return r.Err }

func makeDeleteTeamScheduleEndpoint(svc fleet.Service, opts []kithttp.ServerOption) http.Handler {
	return newServer(
		makeAuthenticatedServiceEndpoint(svc, deleteTeamScheduleEndpoint),
		makeDecoderForIDs(deleteTeamScheduleRequest{}, "team_id", "scheduled_query_id"),
		opts,
	)
}

func deleteTeamScheduleEndpoint(ctx context.Context, request interface{}, svc fleet.Service) (interface{}, error) {
	req := request.(*deleteTeamScheduleRequest)
	err := svc.DeleteTeamScheduledQueries(ctx, req.TeamID, req.ScheduledQueryID)
	if err != nil {
		return deleteTeamScheduleResponse{Err: err}, nil
	}
	return deleteTeamScheduleResponse{}, nil
}

func (svc Service) DeleteTeamScheduledQueries(ctx context.Context, teamID uint, scheduledQueryID uint) error {
	if err := svc.authz.Authorize(ctx, &fleet.Pack{}, fleet.ActionWrite); err != nil {
		return err
	}
	_ = teamID
	return svc.DeleteScheduledQuery(ctx, scheduledQueryID)
}

func (mw loggingMiddleware) DeleteTeamScheduledQueries(ctx context.Context, teamID uint, scheduledQueryID uint) error {
	var err error
	defer func(begin time.Time) {
		_ = mw.loggerDebug(err).Log("method", "DeleteTeamScheduledQueries", "err", err, "took", time.Since(begin))
	}(time.Now())
	return mw.Service.DeleteTeamScheduledQueries(ctx, teamID, scheduledQueryID)
}
