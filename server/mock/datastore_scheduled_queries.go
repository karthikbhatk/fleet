// Automatically generated by mockimpl. DO NOT EDIT!

package mock

import "github.com/fleetdm/fleet/v4/server/fleet"

var _ fleet.ScheduledQueryStore = (*ScheduledQueryStore)(nil)

type ListScheduledQueriesInPackFunc func(id uint, opts fleet.ListOptions) ([]*fleet.ScheduledQuery, error)

type NewScheduledQueryFunc func(sq *fleet.ScheduledQuery, opts ...fleet.OptionalArg) (*fleet.ScheduledQuery, error)

type SaveScheduledQueryFunc func(sq *fleet.ScheduledQuery) (*fleet.ScheduledQuery, error)

type DeleteScheduledQueryFunc func(id uint) error

type ScheduledQueryFunc func(id uint) (*fleet.ScheduledQuery, error)

type CleanupOrphanScheduledQueryStatsFunc func() error

type ScheduledQueryStore struct {
	ListScheduledQueriesInPackFunc        ListScheduledQueriesInPackFunc
	ListScheduledQueriesInPackFuncInvoked bool

	NewScheduledQueryFunc        NewScheduledQueryFunc
	NewScheduledQueryFuncInvoked bool

	SaveScheduledQueryFunc        SaveScheduledQueryFunc
	SaveScheduledQueryFuncInvoked bool

	DeleteScheduledQueryFunc        DeleteScheduledQueryFunc
	DeleteScheduledQueryFuncInvoked bool

	ScheduledQueryFunc        ScheduledQueryFunc
	ScheduledQueryFuncInvoked bool

	CleanupOrphanScheduledQueryStatsFunc        CleanupOrphanScheduledQueryStatsFunc
	CleanupOrphanScheduledQueryStatsFuncInvoked bool
}

func (s *ScheduledQueryStore) ListScheduledQueriesInPack(id uint, opts fleet.ListOptions) ([]*fleet.ScheduledQuery, error) {
	s.ListScheduledQueriesInPackFuncInvoked = true
	return s.ListScheduledQueriesInPackFunc(id, opts)
}

func (s *ScheduledQueryStore) NewScheduledQuery(sq *fleet.ScheduledQuery, opts ...fleet.OptionalArg) (*fleet.ScheduledQuery, error) {
	s.NewScheduledQueryFuncInvoked = true
	return s.NewScheduledQueryFunc(sq, opts...)
}

func (s *ScheduledQueryStore) SaveScheduledQuery(sq *fleet.ScheduledQuery) (*fleet.ScheduledQuery, error) {
	s.SaveScheduledQueryFuncInvoked = true
	return s.SaveScheduledQueryFunc(sq)
}

func (s *ScheduledQueryStore) DeleteScheduledQuery(id uint) error {
	s.DeleteScheduledQueryFuncInvoked = true
	return s.DeleteScheduledQueryFunc(id)
}

func (s *ScheduledQueryStore) ScheduledQuery(id uint) (*fleet.ScheduledQuery, error) {
	s.ScheduledQueryFuncInvoked = true
	return s.ScheduledQueryFunc(id)
}

func (s *ScheduledQueryStore) CleanupOrphanScheduledQueryStats() error {
	s.CleanupOrphanScheduledQueryStatsFuncInvoked = true
	return s.CleanupOrphanScheduledQueryStatsFunc()
}
