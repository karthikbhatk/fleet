// Automatically generated by mockimpl. DO NOT EDIT!

package mock

import (
	"time"

	"github.com/fleetdm/fleet/v4/server/fleet"
)

var _ fleet.HostStore = (*HostStore)(nil)

type NewHostFunc func(host *fleet.Host) (*fleet.Host, error)

type SaveHostFunc func(host *fleet.Host) error

type DeleteHostFunc func(hid uint) error

type HostFunc func(id uint) (*fleet.Host, error)

type HostByIdentifierFunc func(identifier string) (*fleet.Host, error)

type ListHostsFunc func(filter fleet.TeamFilter, opt fleet.HostListOptions) ([]*fleet.Host, error)

type EnrollHostFunc func(osqueryHostId, nodeKey string, teamID *uint, cooldown time.Duration) (*fleet.Host, error)

type AuthenticateHostFunc func(nodeKey string) (*fleet.Host, error)

type MarkHostSeenFunc func(host *fleet.Host, t time.Time) error

type MarkHostsSeenFunc func(hostIDs []uint, t time.Time) error

type CleanupIncomingHostsFunc func(t time.Time) error

type SearchHostsFunc func(filter fleet.TeamFilter, query string, omit ...uint) ([]*fleet.Host, error)

type GenerateHostStatusStatisticsFunc func(filter fleet.TeamFilter, now time.Time) (online uint, offline uint, mia uint, new uint, err error)

type DistributedQueriesForHostFunc func(host *fleet.Host) (map[uint]string, error)

type HostIDsByNameFunc func(filter fleet.TeamFilter, hostnames []string) ([]uint, error)

type AddHostsToTeamFunc func(teamID *uint, hostIDs []uint) error

type SaveHostAdditionalFunc func(host *fleet.Host) error

type HostStore struct {
	NewHostFunc        NewHostFunc
	NewHostFuncInvoked bool

	SaveHostFunc        SaveHostFunc
	SaveHostFuncInvoked bool

	DeleteHostFunc        DeleteHostFunc
	DeleteHostFuncInvoked bool

	HostFunc        HostFunc
	HostFuncInvoked bool

	HostByIdentifierFunc        HostByIdentifierFunc
	HostByIdentifierFuncInvoked bool

	ListHostsFunc        ListHostsFunc
	ListHostsFuncInvoked bool

	EnrollHostFunc        EnrollHostFunc
	EnrollHostFuncInvoked bool

	AuthenticateHostFunc        AuthenticateHostFunc
	AuthenticateHostFuncInvoked bool

	MarkHostSeenFunc        MarkHostSeenFunc
	MarkHostSeenFuncInvoked bool

	MarkHostsSeenFunc        MarkHostsSeenFunc
	MarkHostsSeenFuncInvoked bool

	CleanupIncomingHostsFunc        CleanupIncomingHostsFunc
	CleanupIncomingHostsFuncInvoked bool

	SearchHostsFunc        SearchHostsFunc
	SearchHostsFuncInvoked bool

	GenerateHostStatusStatisticsFunc        GenerateHostStatusStatisticsFunc
	GenerateHostStatusStatisticsFuncInvoked bool

	DistributedQueriesForHostFunc        DistributedQueriesForHostFunc
	DistributedQueriesForHostFuncInvoked bool

	HostIDsByNameFunc        HostIDsByNameFunc
	HostIDsByNameFuncInvoked bool

	AddHostsToTeamFunc        AddHostsToTeamFunc
	AddHostsToTeamFuncInvoked bool

	SaveHostAdditionalFunc        SaveHostAdditionalFunc
	SaveHostAdditionalFuncInvoked bool
}

func (s *HostStore) NewHost(host *fleet.Host) (*fleet.Host, error) {
	s.NewHostFuncInvoked = true
	return s.NewHostFunc(host)
}

func (s *HostStore) SaveHost(host *fleet.Host) error {
	s.SaveHostFuncInvoked = true
	return s.SaveHostFunc(host)
}

func (s *HostStore) DeleteHost(hid uint) error {
	s.DeleteHostFuncInvoked = true
	return s.DeleteHostFunc(hid)
}

func (s *HostStore) Host(id uint) (*fleet.Host, error) {
	s.HostFuncInvoked = true
	return s.HostFunc(id)
}

func (s *HostStore) HostByIdentifier(identifier string) (*fleet.Host, error) {
	s.HostByIdentifierFuncInvoked = true
	return s.HostByIdentifierFunc(identifier)
}

func (s *HostStore) ListHosts(filter fleet.TeamFilter, opt fleet.HostListOptions) ([]*fleet.Host, error) {
	s.ListHostsFuncInvoked = true
	return s.ListHostsFunc(filter, opt)
}

func (s *HostStore) EnrollHost(osqueryHostId, nodeKey string, teamID *uint, cooldown time.Duration) (*fleet.Host, error) {
	s.EnrollHostFuncInvoked = true
	return s.EnrollHostFunc(osqueryHostId, nodeKey, teamID, cooldown)
}

func (s *HostStore) AuthenticateHost(nodeKey string) (*fleet.Host, error) {
	s.AuthenticateHostFuncInvoked = true
	return s.AuthenticateHostFunc(nodeKey)
}

func (s *HostStore) MarkHostSeen(host *fleet.Host, t time.Time) error {
	s.MarkHostSeenFuncInvoked = true
	return s.MarkHostSeenFunc(host, t)
}

func (s *HostStore) MarkHostsSeen(hostIDs []uint, t time.Time) error {
	s.MarkHostsSeenFuncInvoked = true
	return s.MarkHostsSeenFunc(hostIDs, t)
}

func (s *HostStore) CleanupIncomingHosts(t time.Time) error {
	s.CleanupIncomingHostsFuncInvoked = true
	return s.CleanupIncomingHostsFunc(t)
}

func (s *HostStore) SearchHosts(filter fleet.TeamFilter, query string, omit ...uint) ([]*fleet.Host, error) {
	s.SearchHostsFuncInvoked = true
	return s.SearchHostsFunc(filter, query, omit...)
}

func (s *HostStore) GenerateHostStatusStatistics(filter fleet.TeamFilter, now time.Time) (online uint, offline uint, mia uint, new uint, err error) {
	s.GenerateHostStatusStatisticsFuncInvoked = true
	return s.GenerateHostStatusStatisticsFunc(filter, now)
}

func (s *HostStore) DistributedQueriesForHost(host *fleet.Host) (map[uint]string, error) {
	s.DistributedQueriesForHostFuncInvoked = true
	return s.DistributedQueriesForHostFunc(host)
}

func (s *HostStore) HostIDsByName(filter fleet.TeamFilter, hostnames []string) ([]uint, error) {
	s.HostIDsByNameFuncInvoked = true
	return s.HostIDsByNameFunc(filter, hostnames)
}

func (s *HostStore) AddHostsToTeam(teamID *uint, hostIDs []uint) error {
	s.AddHostsToTeamFuncInvoked = true
	return s.AddHostsToTeamFunc(teamID, hostIDs)
}

func (s *HostStore) SaveHostAdditional(host *fleet.Host) error {
	s.SaveHostAdditionalFuncInvoked = true
	return s.SaveHostAdditionalFunc(host)
}
