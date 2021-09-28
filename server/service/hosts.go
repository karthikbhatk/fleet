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
	Deleted []uint `json:"deleted,omitempty"`
	Err     error  `json:"error,omitempty"`
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
	resp, err := svc.DeleteHosts(ctx, req.IDs, listOpt, req.Filters.LabelID)
	if err != nil {
		return deleteHostsResponse{Err: err}, nil
	}
	return deleteHostsResponse{Deleted: resp}, nil
}

func (svc Service) DeleteHosts(ctx context.Context, ids []uint, opts fleet.HostListOptions, lid *uint) ([]uint, error) {
	if err := svc.authz.Authorize(ctx, &fleet.Host{}, fleet.ActionWrite); err != nil {
		return nil, err
	}

	if len(ids) > 0 {
		return svc.ds.DeleteHosts(ctx, ids)
	}

	hostIDs, err := svc.hostIDsFromFilters(ctx, opts, lid)
	if err != nil {
		return nil, err
	}

	if len(hostIDs) == 0 {
		return nil, nil
	}
	return svc.ds.DeleteHosts(ctx, hostIDs)
}
