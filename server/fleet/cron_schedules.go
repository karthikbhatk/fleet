package fleet

import (
	"fmt"
	"net/http"
	"time"
)

type CronSchedulesService interface {
	// TriggerCronSchedule attempts to trigger an ad-hoc run of the named cron schedule.
	TriggerCronSchedule(name string) error
	// // GetCronScheduleNames returns a list of the names of all cron schedules registered with the service.
	// GetCronScheduleNames() []string
}

type CronSchedule interface {
	Trigger() (bool, *CronStats, error)
	Name() string
}

type CronSchedules struct {
	Schedules map[string]CronSchedule
}

// AuthzType implements authz.AuthzTyper.
func (cs *CronSchedules) AuthzType() string {
	return "cron_schedules"
}

// AddCronSchedules registers a new cron schedule with the service.
func (cs *CronSchedules) AddCronSchedule(sched CronSchedule, err error) error {
	if err != nil {
		return err
	}
	cs.Schedules[sched.Name()] = sched
	return nil
}

// TriggerCronSchedule attempts to trigger an ad-hoc run of the named cron schedule.
func (cs *CronSchedules) TriggerCronSchedule(name string) error {
	sched, ok := cs.Schedules[name]
	if !ok {
		return triggerNotFoundError{}
	}
	ok, stats, err := sched.Trigger()
	switch {
	case err != nil:
		return err
	case !ok:
		if stats == nil || string(stats.Status) == "" {
			return triggerConflictError{name: name}
		}
		return triggerConflictError{name: name, stats: stats}
	default:
		return nil
	}
}

// GetCronScheduleNames returns a list of the names of all cron schedules registered with the service.
func (cs *CronSchedules) GetCronScheduleNames() []string {
	var res []string
	for _, sched := range cs.Schedules {
		res = append(res, sched.Name())
	}
	return res
}

type triggerConflictError struct {
	name  string
	stats *CronStats
}

func (e triggerConflictError) Error() string {
	msg := "conflicts with current status"
	if e.stats != nil {
		msg += fmt.Sprintf(": %s %s %s run started %v ago", e.stats.Status, e.stats.Name, e.stats.StatsType, time.Since(e.stats.CreatedAt))
	}
	return msg
}

func (a triggerConflictError) StatusCode() int {
	return http.StatusConflict
}

type triggerNotFoundError struct {
	name string
}

func (e triggerNotFoundError) Error() string {
	msg := "unrecognized name"
	if e.name != "" {
		msg += fmt.Sprintf(": %s", e.name)
	}
	return msg
}

func (e triggerNotFoundError) IsNotFound() bool {
	return true
}

// CronStats represents statistics recorded in connection with a named set of jobs (sometimes
// referred to as a "cron" or "schedule"). Each record represents a separate "run" of the named job set.
type CronStats struct {
	ID int `db:"id"`
	// StatsType denotes whether the stats are associated with a run of jobs that was "triggered"
	// (i.e. run on an ad-hoc basis) or "scheduled" (i.e. run on a regularly scheduled interval).
	StatsType CronStatsType `db:"stats_type"`
	// Name is the name of the set of jobs (i.e. the schedule name).
	Name string `db:"name"`
	// Instance is the unique id of the Fleet instance that performed the run of jobs represented by
	// the stats.
	Instance string `db:"instance"`
	// CreatedAt is the time the stats record was created. It is assumed to be the start of the run.
	CreatedAt time.Time `db:"created_at"`
	// UpdatedAt is the time the stats record was last updated. For a "completed" run, this assumed
	// to be the end of the run.
	UpdatedAt time.Time `db:"updated_at"`
	// Status is the current status of the run. Recognized statuses are "pending", "completed", and
	// "expired".
	Status CronStatsStatus `db:"status"`
}

// CronStatsType is one of two recognized types of cron stats (i.e. "scheduled" or "triggered")
type CronStatsType string

// List of recognized cron stats types.
const (
	CronStatsTypeScheduled CronStatsType = "scheduled"
	CronStatsTypeTriggered CronStatsType = "triggered"
)

// CronStatsStatus is one of four recognized statuses of cron stats (i.e. "pending", "expired", "canceled", or "completed")
type CronStatsStatus string

// List of recognized cron stats statuses.
const (
	CronStatsStatusPending   CronStatsStatus = "pending"
	CronStatsStatusExpired   CronStatsStatus = "expired"
	CronStatsStatusCompleted CronStatsStatus = "completed"
	CronStatsStatusCanceled  CronStatsStatus = "canceled"
)
