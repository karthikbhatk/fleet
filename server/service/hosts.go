package service

import (
	"context"

	"github.com/fleetdm/fleet/v4/server/fleet"
)

/////////////////////////////////////////////////////////////////////////////////
// Delete
/////////////////////////////////////////////////////////////////////////////////

type deleteHostsRequest struct {
	IDs     []uint `json:"ids"`
	Filters struct {
		MatchQuery string           `json:"query"`
		Status     fleet.HostStatus `json:"status"`
		LabelID    *uint            `json:"label_id"`
		TeamID     *uint            `json:"team_id"`
	} `json:"filters"`
}

type deleteHostsResponse struct {
	Err error `json:"error,omitempty"`
}

func (r deleteHostsResponse) error() error { return r.Err }

func deleteHostsEndpoint(ctx context.Context, request interface{}, svc fleet.Service) (interface{}, error) {
	req := request.(*deleteHostsRequest)
	listOpt := fleet.HostListOptions{
		ListOptions: fleet.ListOptions{
			MatchQuery: req.Filters.MatchQuery,
		},
		StatusFilter: req.Filters.Status,
		TeamFilter:   req.Filters.TeamID,
	}
	err := svc.DeleteHosts(ctx, req.IDs, listOpt, req.Filters.LabelID)
	if err != nil {
		return deleteHostsResponse{Err: err}, nil
	}
	return deleteHostsResponse{}, nil
}

func (svc Service) DeleteHosts(ctx context.Context, ids []uint, opts fleet.HostListOptions, lid *uint) error {
	if err := svc.authz.Authorize(ctx, &fleet.Host{}, fleet.ActionWrite); err != nil {
		return err
	}

	if len(ids) > 0 && (lid != nil || !opts.Empty()) {
		return &badRequestError{"Cannot specify a list of ids and filters at the same time"}
	}

	if len(ids) > 0 {
		return svc.ds.DeleteHosts(ctx, ids)
	}

	hostIDs, err := svc.hostIDsFromFilters(ctx, opts, lid)
	if err != nil {
		return err
	}

	if len(hostIDs) == 0 {
		return nil
	}
	return svc.ds.DeleteHosts(ctx, hostIDs)
}

/////////////////////////////////////////////////////////////////////////////////
// Count
/////////////////////////////////////////////////////////////////////////////////

type countHostsRequest struct {
	Opts fleet.HostListOptions `url:"host_options"`
}

type countHostsResponse struct {
	Count uint  `json:"count"`
	Err   error `json:"error,omitempty"`
}

func (r countHostsResponse) error() error { return r.Err }

func countHostsEndpoint(ctx context.Context, request interface{}, svc fleet.Service) (interface{}, error) {
	req := request.(*countHostsRequest)
	count, err := svc.CountHosts(ctx, req.Opts)
	if err != nil {
		return countHostsResponse{Err: err}, nil
	}
	return countHostsResponse{Count: count}, nil
}

func (svc Service) CountHosts(ctx context.Context, opts fleet.HostListOptions) (uint, error) {
	if err := svc.authz.Authorize(ctx, &fleet.Host{}, fleet.ActionRead); err != nil {
		return 0, err
	}

	return svc.countHostFromFilters(ctx, opts)
}

/////////////////////////////////////////////////////////////////////////////////
// Get host
/////////////////////////////////////////////////////////////////////////////////

func getHostEndpoint(ctx context.Context, request interface{}, svc fleet.Service) (interface{}, error) {
	req := request.(*getHostRequest)
	host, err := svc.GetHost(ctx, req.ID)
	if err != nil {
		return getHostResponse{Err: err}, nil
	}

	resp, err := hostDetailResponseForHost(ctx, svc, host)
	if err != nil {
		return getHostResponse{Err: err}, nil
	}

	return getHostResponse{Host: resp}, nil
}
