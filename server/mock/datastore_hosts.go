// Automatically generated by mockimpl. DO NOT EDIT!

package mock

import (
	"time"

	"github.com/fleetdm/fleet/server/kolide"
)

var _ kolide.HostStore = (*HostStore)(nil)

type NewHostFunc func(host *kolide.Host) (*kolide.Host, error)

type SaveHostFunc func(host *kolide.Host) error

type DeleteHostFunc func(hid uint) error

type HostFunc func(id uint) (*kolide.Host, error)

type HostByIdentifierFunc func(identifier string) (*kolide.Host, error)

type ListHostsFunc func(opt kolide.HostListOptions) ([]*kolide.Host, error)

type EnrollHostFunc func(osqueryHostId, nodeKey, secretName string, cooldown time.Duration) (*kolide.Host, error)

type AuthenticateHostFunc func(nodeKey string) (*kolide.Host, error)

type MarkHostSeenFunc func(host *kolide.Host, t time.Time) error

type MarkHostsSeenFunc func(hostIDs []uint, t time.Time) error

type CleanupIncomingHostsFunc func(t time.Time) error

type SearchHostsFunc func(query string, omit ...uint) ([]*kolide.Host, error)

type GenerateHostStatusStatisticsFunc func(now time.Time) (online uint, offline uint, mia uint, new uint, err error)

type DistributedQueriesForHostFunc func(host *kolide.Host) (map[uint]string, error)

type HostIDsByNameFunc func(hostnames []string) ([]uint, error)

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
}

func (s *HostStore) NewHost(host *kolide.Host) (*kolide.Host, error) {
	s.NewHostFuncInvoked = true
	return s.NewHostFunc(host)
}

func (s *HostStore) SaveHost(host *kolide.Host) error {
	s.SaveHostFuncInvoked = true
	return s.SaveHostFunc(host)
}

func (s *HostStore) DeleteHost(hid uint) error {
	s.DeleteHostFuncInvoked = true
	return s.DeleteHostFunc(hid)
}

func (s *HostStore) Host(id uint) (*kolide.Host, error) {
	s.HostFuncInvoked = true
	return s.HostFunc(id)
}

func (s *HostStore) HostByIdentifier(identifier string) (*kolide.Host, error) {
	s.HostByIdentifierFuncInvoked = true
	return s.HostByIdentifierFunc(identifier)
}

func (s *HostStore) ListHosts(opt kolide.HostListOptions) ([]*kolide.Host, error) {
	s.ListHostsFuncInvoked = true
	return s.ListHostsFunc(opt)
}

func (s *HostStore) EnrollHost(osqueryHostId, nodeKey, secretName string, cooldown time.Duration) (*kolide.Host, error) {
	s.EnrollHostFuncInvoked = true
	return s.EnrollHostFunc(osqueryHostId, nodeKey, secretName, cooldown)
}

func (s *HostStore) AuthenticateHost(nodeKey string) (*kolide.Host, error) {
	s.AuthenticateHostFuncInvoked = true
	return s.AuthenticateHostFunc(nodeKey)
}

func (s *HostStore) MarkHostSeen(host *kolide.Host, t time.Time) error {
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

func (s *HostStore) SearchHosts(query string, omit ...uint) ([]*kolide.Host, error) {
	s.SearchHostsFuncInvoked = true
	return s.SearchHostsFunc(query, omit...)
}

func (s *HostStore) GenerateHostStatusStatistics(now time.Time) (online uint, offline uint, mia uint, new uint, err error) {
	s.GenerateHostStatusStatisticsFuncInvoked = true
	return s.GenerateHostStatusStatisticsFunc(now)
}

func (s *HostStore) DistributedQueriesForHost(host *kolide.Host) (map[uint]string, error) {
	s.DistributedQueriesForHostFuncInvoked = true
	return s.DistributedQueriesForHostFunc(host)
}

func (s *HostStore) HostIDsByName(hostnames []string) ([]uint, error) {
	s.HostIDsByNameFuncInvoked = true
	return s.HostIDsByNameFunc(hostnames)
}
