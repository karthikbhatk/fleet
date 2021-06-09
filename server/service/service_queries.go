package service

import (
	"context"

	"github.com/fleetdm/fleet/server/contexts/viewer"
	"github.com/fleetdm/fleet/server/fleet"
	"github.com/fleetdm/fleet/server/ptr"
	"github.com/pkg/errors"
)

func queryFromSpec(spec *fleet.QuerySpec) *fleet.Query {
	return &fleet.Query{
		Name:        spec.Name,
		Description: spec.Description,
		Query:       spec.Query,
	}
}

func specFromQuery(query *fleet.Query) *fleet.QuerySpec {
	return &fleet.QuerySpec{
		Name:        query.Name,
		Description: query.Description,
		Query:       query.Query,
	}
}

func (svc Service) ApplyQuerySpecs(ctx context.Context, specs []*fleet.QuerySpec) error {
	if err := svc.authz.Authorize(ctx, &fleet.Query{}, fleet.ActionWrite); err != nil {
		return err
	}

	vc, ok := viewer.FromContext(ctx)
	if !ok {
		return errors.New("user must be authenticated to apply queries")
	}

	queries := []*fleet.Query{}
	for _, spec := range specs {
		queries = append(queries, queryFromSpec(spec))
	}

	for _, query := range queries {
		if err := query.ValidateSQL(); err != nil {
			return err
		}
	}

	err := svc.ds.ApplyQueries(vc.UserID(), queries)
	return errors.Wrap(err, "applying queries")
}

func (svc Service) GetQuerySpecs(ctx context.Context) ([]*fleet.QuerySpec, error) {
	if err := svc.authz.Authorize(ctx, &fleet.Query{}, fleet.ActionRead); err != nil {
		return nil, err
	}

	queries, err := svc.ds.ListQueries(fleet.ListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "getting queries")
	}

	specs := []*fleet.QuerySpec{}
	for _, query := range queries {
		specs = append(specs, specFromQuery(query))
	}
	return specs, nil
}

func (svc Service) GetQuerySpec(ctx context.Context, name string) (*fleet.QuerySpec, error) {
	if err := svc.authz.Authorize(ctx, &fleet.Query{}, fleet.ActionRead); err != nil {
		return nil, err
	}

	query, err := svc.ds.QueryByName(name)
	if err != nil {
		return nil, err
	}
	return specFromQuery(query), nil
}

func (svc Service) ListQueries(ctx context.Context, opt fleet.ListOptions) ([]*fleet.Query, error) {
	if err := svc.authz.Authorize(ctx, &fleet.Query{}, fleet.ActionRead); err != nil {
		return nil, err
	}

	return svc.ds.ListQueries(opt)
}

func (svc *Service) GetQuery(ctx context.Context, id uint) (*fleet.Query, error) {
	if err := svc.authz.Authorize(ctx, &fleet.Query{}, fleet.ActionRead); err != nil {
		return nil, err
	}

	return svc.ds.Query(id)
}

func (svc *Service) NewQuery(ctx context.Context, p fleet.QueryPayload) (*fleet.Query, error) {
	if err := svc.authz.Authorize(ctx, &fleet.Query{}, fleet.ActionWrite); err != nil {
		return nil, err
	}

	query := &fleet.Query{Saved: true}

	if p.Name != nil {
		query.Name = *p.Name
	}

	if p.Description != nil {
		query.Description = *p.Description
	}

	if p.Query != nil {
		query.Query = *p.Query
	}

	if p.ObserverCanRun != nil {
		query.ObserverCanRun = *p.ObserverCanRun
	}

	vc, ok := viewer.FromContext(ctx)
	if ok {
		query.AuthorID = ptr.Uint(vc.UserID())
		query.AuthorName = vc.FullName()
	}

	if err := query.ValidateSQL(); err != nil {
		return nil, err
	}

	query, err := svc.ds.NewQuery(query)
	if err != nil {
		return nil, err
	}

	return query, nil
}

func (svc *Service) ModifyQuery(ctx context.Context, id uint, p fleet.QueryPayload) (*fleet.Query, error) {
	if err := svc.authz.Authorize(ctx, &fleet.Query{}, fleet.ActionWrite); err != nil {
		return nil, err
	}

	query, err := svc.ds.Query(id)
	if err != nil {
		return nil, err
	}

	if p.Name != nil {
		query.Name = *p.Name
	}

	if p.Description != nil {
		query.Description = *p.Description
	}

	if p.Query != nil {
		query.Query = *p.Query
	}

	if p.ObserverCanRun != nil {
		query.ObserverCanRun = *p.ObserverCanRun
	}

	if err := query.ValidateSQL(); err != nil {
		return nil, err
	}

	if err := svc.ds.SaveQuery(query); err != nil {
		return nil, err
	}

	return query, nil
}

func (svc *Service) DeleteQuery(ctx context.Context, name string) error {
	if err := svc.authz.Authorize(ctx, &fleet.Query{}, fleet.ActionWrite); err != nil {
		return err
	}

	return svc.ds.DeleteQuery(name)
}

func (svc *Service) DeleteQueryByID(ctx context.Context, id uint) error {
	if err := svc.authz.Authorize(ctx, &fleet.Query{}, fleet.ActionWrite); err != nil {
		return err
	}

	query, err := svc.ds.Query(id)
	if err != nil {
		return errors.Wrap(err, "lookup query by ID")
	}

	return errors.Wrap(svc.ds.DeleteQuery(query.Name), "delete query")
}

func (svc *Service) DeleteQueries(ctx context.Context, ids []uint) (uint, error) {
	if err := svc.authz.Authorize(ctx, &fleet.Query{}, fleet.ActionWrite); err != nil {
		return 0, err
	}

	return svc.ds.DeleteQueries(ids)
}
