// Automatically generated by mockimpl. DO NOT EDIT!

package mock

import (
	"time"

	"github.com/fleetdm/fleet/v4/server/fleet"
)

var _ fleet.TargetStore = (*TargetStore)(nil)

type CountHostsInTargetsFunc func(filter fleet.TeamFilter, targets fleet.HostTargets, now time.Time) (fleet.TargetMetrics, error)

type HostIDsInTargetsFunc func(filter fleet.TeamFilter, targets fleet.HostTargets) ([]uint, error)

type TargetStore struct {
	CountHostsInTargetsFunc        CountHostsInTargetsFunc
	CountHostsInTargetsFuncInvoked bool

	HostIDsInTargetsFunc        HostIDsInTargetsFunc
	HostIDsInTargetsFuncInvoked bool
}

func (s *TargetStore) CountHostsInTargets(filter fleet.TeamFilter, targets fleet.HostTargets, now time.Time) (fleet.TargetMetrics, error) {
	s.CountHostsInTargetsFuncInvoked = true
	return s.CountHostsInTargetsFunc(filter, targets, now)
}

func (s *TargetStore) HostIDsInTargets(filter fleet.TeamFilter, targets fleet.HostTargets) ([]uint, error) {
	s.HostIDsInTargetsFuncInvoked = true
	return s.HostIDsInTargetsFunc(filter, targets)
}
