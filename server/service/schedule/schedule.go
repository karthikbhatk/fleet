// Package schedule allows periodic run of a list of jobs.
//
// Type Schedule allows grouping a set of Jobs to run at specific intervals.
// Each Job is executed serially in the order they were added to the Schedule.
package schedule

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/fleetdm/fleet/v4/server/contexts/ctxerr"
	"github.com/fleetdm/fleet/v4/server/fleet"
	"github.com/getsentry/sentry-go"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// ReloadInterval reloads and returns a new interval.
type ReloadInterval func(ctx context.Context) (time.Duration, error)

// Schedule runs a list of jobs serially at a given schedule.
//
// Each job is executed one after the other in the order they were added.
// If one of the job fails, an error is logged and the scheduler
// continues with the next.
type Schedule struct {
	ctx        context.Context
	name       string
	instanceID string
	logger     log.Logger

	schedIntervalMu sync.Mutex // protects schedInterval.
	schedInterval   time.Duration

	done chan struct{}

	configReloadInterval   time.Duration
	configReloadIntervalFn ReloadInterval

	locker Locker

	altLockName string

	jobs []Job

	statsStore CronStatsStore
}

// JobFn is the signature of a Job.
type JobFn func(context.Context) error

// Job represents a job that can be added to Scheduler.
type Job struct {
	// ID is the unique identifier for the job.
	ID string
	// Fn is the job itself.
	Fn JobFn
}

// Locker allows a Schedule to acquire a lock before running jobs.
type Locker interface {
	Lock(ctx context.Context, scheduleName string, scheduleInstanceID string, expiration time.Duration) (bool, error)
	Unlock(ctx context.Context, scheduleName string, scheduleInstanceID string) error
}

// CronStatsStore allows a Schedule to store and retrieve statistics pertaining to the Schedule
type CronStatsStore interface {
	// GetLatestCronStats returns the most recent cron stats for the named cron schedule. If no rows
	// are found, it returns an empty CronStats struct
	GetLatestCronStats(ctx context.Context, name string) (fleet.CronStats, error)
	// InsertCronStats inserts cron stats for the named cron schedule
	InsertCronStats(ctx context.Context, statsType fleet.CronStatsType, name string, instance string, status fleet.CronStatsStatus) (int, error)
	// UpdateCronStats updates the status of the identified cron stats record
	UpdateCronStats(ctx context.Context, id int, status fleet.CronStatsStatus) error
}

// Option allows configuring a Schedule.
type Option func(*Schedule)

// WithLogger sets a logger for the Schedule.
func WithLogger(l log.Logger) Option {
	return func(s *Schedule) {
		s.logger = log.With(l, "schedule", s.name)
	}
}

// WithConfigReloadInterval allows setting a reload interval function,
// that will allow updating the interval of a running schedule.
//
// If not set, then the schedule performs no interval reloading.
func WithConfigReloadInterval(interval time.Duration, fn ReloadInterval) Option {
	return func(s *Schedule) {
		s.configReloadInterval = interval
		s.configReloadIntervalFn = fn
	}
}

// WithAltLockID sets an alternative identifier to use when acquiring the lock.
//
// If not set, then the Schedule's name is used for acquiring the lock.
func WithAltLockID(name string) Option {
	return func(s *Schedule) {
		s.altLockName = name
	}
}

// WithJob adds a job to the Schedule.
//
// Each job is executed in the order they are added.
func WithJob(id string, fn JobFn) Option {
	return func(s *Schedule) {
		s.jobs = append(s.jobs, Job{
			ID: id,
			Fn: fn,
		})
	}
}

// New creates and returns a Schedule.
// Jobs are added with the WithJob Option.
//
// The jobs are executed serially in order at the provided interval.
//
// The provided locker is used to acquire/release a lock before running the jobs.
// The provided name and instanceID of the Schedule is used as the locking identifier.
func New(
	ctx context.Context,
	name string,
	instanceID string,
	interval time.Duration,
	locker Locker,
	statsStore CronStatsStore,
	opts ...Option,
) *Schedule {
	sch := &Schedule{
		ctx:                  ctx,
		name:                 name,
		instanceID:           instanceID,
		logger:               log.NewNopLogger(),
		done:                 make(chan struct{}),
		configReloadInterval: 1 * time.Hour, // by default we will check for updated config once per hour
		schedInterval:        truncateSecondsWithFloor(interval),
		locker:               locker,
		statsStore:           statsStore,
	}
	for _, fn := range opts {
		fn(sch)
	}
	return sch
}

// Start starts running the added jobs.
//
// All jobs must be added before calling Start.
func (s *Schedule) Start() {
	var intervalStartedAt time.Time // start time of the most recent run of the scheduled jobs
	var m sync.Mutex                // protects intervalStartedAt

	getIntervalStartedAt := func() (start time.Time) {
		m.Lock()
		defer m.Unlock()

		return intervalStartedAt
	}

	setIntervalStartedAt := func(start time.Time) {
		m.Lock()
		defer m.Unlock()

		intervalStartedAt = start.Truncate(time.Second)
	}

	stats, err := s.getStats()
	if err != nil {
		level.Error(s.logger).Log("err", "start schedule", "details", err)
		sentry.CaptureException(err)
		ctxerr.Handle(s.ctx, err)
	}
	setIntervalStartedAt(stats.CreatedAt)

	initialWait := 10 * time.Second
	if schedInterval := s.getSchedInterval(); schedInterval < initialWait {
		initialWait = schedInterval
	}
	schedTicker := time.NewTicker(initialWait)

	var g sync.WaitGroup
	g.Add(+1)
	go func() {
		defer func() {
			s.releaseLock()
			g.Done()
		}()

		for {
			level.Debug(s.logger).Log("waiting", "wait for next tick")

			select {
			case <-s.ctx.Done():
				schedTicker.Stop()
				return

			case <-schedTicker.C:
				level.Debug(s.logger).Log("waiting", "done")

				schedInterval := s.getSchedInterval()

				stats, err := s.getStats()
				if err != nil {
					level.Error(s.logger).Log("err", "get cron stats", "details", err)
					sentry.CaptureException(err)
					ctxerr.Handle(s.ctx, err)
					// skip ahead to the next interval
					schedTicker.Reset(schedInterval)
					continue
				}

				if stats.Status == fleet.CronStatsStatusPending {
					// skip ahead to the next interval
					schedTicker.Reset(schedInterval)
					continue
				}

				prevStart := getIntervalStartedAt()
				if stats.CreatedAt.After(prevStart) {
					// if there's a diff between the datastore and our local value, we use the
					// more recent timestamp and update our local value accordingly
					setIntervalStartedAt(stats.CreatedAt)
					prevStart = getIntervalStartedAt()
				}
				remainingInterval := getRemainingInterval(prevStart, schedInterval)

				if time.Since(prevStart) < schedInterval {
					// wait for the remaining interval plus a small buffer
					schedTicker.Reset(remainingInterval + 100*time.Millisecond)
					continue
				}

				prevFinish := stats.UpdatedAt.Truncate(time.Second)
				prevRuntime := prevFinish.Sub(prevStart)
				if prevRuntime > schedInterval {
					// if the previous run took longer than the schedule interval, we wait until the start of the next full interval
					newStart := prevStart.Add(time.Since(prevStart).Truncate(schedInterval)) // advances start time by the number of full interval elasped
					setIntervalStartedAt(newStart)
					schedTicker.Reset(getRemainingInterval(newStart, schedInterval))
					continue
				}

				ok, cancelHold := s.holdLock()
				if !ok {
					// failed to get a lock so skip ahead to the next interval
					schedTicker.Reset(schedInterval)
					continue
				}

				newStart := time.Now()
				setIntervalStartedAt(newStart)
				level.Info(s.logger).Log("status", "pending")

				statsID, err := s.insertStats(fleet.CronStatsTypeScheduled, fleet.CronStatsStatusPending)
				if err != nil {
					level.Error(s.logger).Log("err", fmt.Sprintf("insert cron stats %s", s.name), "details", err)
					sentry.CaptureException(err)
					ctxerr.Handle(s.ctx, err)
				}

				for _, job := range s.jobs {
					level.Debug(s.logger).Log("msg", "starting", "jobID", job.ID)
					if err := runJob(s.ctx, job.Fn); err != nil {
						level.Error(s.logger).Log("err", "running job", "details", err, "jobID", job.ID)
						sentry.CaptureException(err)
						ctxerr.Handle(s.ctx, err)
					}
				}
				level.Info(s.logger).Log("status", "completed")

				if err := s.updateStats(statsID, fleet.CronStatsStatusCompleted); err != nil {
					level.Error(s.logger).Log("err", fmt.Sprintf("update cron stats %s", s.name), "details", err)
					sentry.CaptureException(err)
					ctxerr.Handle(s.ctx, err)
				}

				// we need to re-synchronize this schedule instance so that the next scheduled run
				// starts at the beginning of the next full interval
				//
				// for example, if the interval is 1hr and the schedule takes 0.2 hrs to run
				// then we wait 0.8 hrs until the next time we run the schedule, or if the
				// the schedule takes 1.5 hrs to run then we wait 0.5 hrs (skipping the scheduled
				// tick that would have overlapped with the 1.5hrs running time)
				schedInterval = s.getSchedInterval()
				if time.Since(newStart) > schedInterval {
					level.Info(s.logger).Log("msg", fmt.Sprintf("total runtime (%v) exceeded schedule interval (%v)", time.Since(newStart), schedInterval))
					newStart = newStart.Add(time.Since(newStart).Truncate(schedInterval)) // advances start time by the number of full interval elasped
				}
				remainingInterval = getRemainingInterval(newStart, schedInterval)
				clearTickerChannel(schedTicker) // in case another tick arrived during this run
				schedTicker.Reset(remainingInterval)
				cancelHold()
			}
		}
	}()

	if s.configReloadIntervalFn != nil {
		// WithConfigReloadInterval option applies so we periodically check for config updates and
		// reset the schedInterval for the previous loop
		g.Add(+1)
		go func() {
			defer g.Done()

			configTicker := time.NewTicker(s.configReloadInterval)
			for {
				select {
				case <-s.ctx.Done():
					configTicker.Stop()
					return
				case <-configTicker.C:
					schedInterval := s.getSchedInterval()
					newInterval, err := s.configReloadIntervalFn(s.ctx)
					if err != nil {
						level.Error(s.logger).Log("err", "schedule interval config reload failed", "details", err)
						sentry.CaptureException(err)
						continue
					}

					newInterval = truncateSecondsWithFloor(newInterval)
					if newInterval <= 0 {
						level.Debug(s.logger).Log("msg", "config reload interval method returned invalid interval")
						continue
					}
					if schedInterval == newInterval {
						continue
					}
					s.setSchedInterval(newInterval)

					intervalStartedAt := getIntervalStartedAt()
					newWait := 10 * time.Millisecond
					if time.Since(intervalStartedAt) < schedInterval {
						newWait = schedInterval - time.Since(intervalStartedAt)
					}

					clearTickerChannel(schedTicker)
					schedTicker.Reset(newWait)

					level.Debug(s.logger).Log("msg", fmt.Sprintf("new schedule interval %v", schedInterval))
					level.Debug(s.logger).Log("msg", fmt.Sprintf("time until next schedule tick %v", newWait))
				}
			}
		}()
	}

	go func() {
		g.Wait()
		level.Debug(s.logger).Log("msg", "close schedule")
		close(s.done) // communicates that the scheduler has finished running its goroutines
		schedTicker.Stop()
	}()
}

// runJob executes the job function with panic recovery
func runJob(ctx context.Context, fn JobFn) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()

	if err := fn(ctx); err != nil {
		return err
	}
	return nil
}

// Done returns a channel that will be closed when the scheduler's context is done
// and it has finished running its goroutines.
func (s *Schedule) Done() <-chan struct{} {
	return s.done
}

// getScheduleInterval returns the schedule interval
func (s *Schedule) getSchedInterval() time.Duration {
	s.schedIntervalMu.Lock()
	defer s.schedIntervalMu.Unlock()

	return s.schedInterval
}

// setScheduleInterval sets the schedule interval after truncating the duration to seconds and
// applying a one second floor (e.g., 600ms becomes 1s, 1300ms becomes 2s, 1000ms becomes 2s)
func (s *Schedule) setSchedInterval(interval time.Duration) {
	s.schedIntervalMu.Lock()
	defer s.schedIntervalMu.Unlock()

	s.schedInterval = truncateSecondsWithFloor(interval)
}

func (s *Schedule) acquireLock() bool {
	locked, err := s.locker.Lock(s.ctx, s.getLockName(), s.instanceID, s.getSchedInterval())
	if err != nil {
		level.Error(s.logger).Log("msg", "lock failed", "err", err)
		sentry.CaptureException(err)
		return false
	}
	if locked {
		return true
	}
	level.Debug(s.logger).Log("msg", "not the lock leader, skipping")
	return false
}

func (s *Schedule) releaseLock() {
	err := s.locker.Unlock(s.ctx, s.getLockName(), s.instanceID)
	if err != nil {
		level.Error(s.logger).Log("msg", "unlock failed", "err", err)
		sentry.CaptureException(err)
	}
}

// holdLock attempts to acquire a schedule lock. If it successfully acquires the lock, it starts a
// goroutine that periodically extends the lock, and it returns `true` along with a
// context.CancelFunc that will end the goroutine and release the lock. If it is unable to initially
// acquire a lock, it returns `false, nil`. The maximum duration of the hold is two hours.
func (s *Schedule) holdLock() (bool, context.CancelFunc) {
	if ok := s.acquireLock(); !ok {
		return false, nil
	}

	ctx, cancelFn := context.WithCancel(s.ctx)

	go func() {
		t := time.NewTimer(s.getSchedInterval() * 8 / 10) // hold timer is 80% of schedule interval
		for {
			select {
			case <-ctx.Done():
				if !t.Stop() {
					<-t.C
				}
				s.releaseLock()
				return
			case <-t.C:
				s.acquireLock()
				t.Reset(s.getSchedInterval() * 8 / 10)
			}
		}
	}()

	return true, cancelFn
}

func (s *Schedule) getStats() (fleet.CronStats, error) {
	return s.statsStore.GetLatestCronStats(s.ctx, s.name)
}

func (s *Schedule) insertStats(statsType fleet.CronStatsType, status fleet.CronStatsStatus) (int, error) {
	return s.statsStore.InsertCronStats(s.ctx, statsType, s.name, s.instanceID, status)
}

func (s *Schedule) updateStats(id int, status fleet.CronStatsStatus) error {
	return s.statsStore.UpdateCronStats(s.ctx, id, status)
}

func (s *Schedule) getLockName() string {
	name := s.name
	if s.altLockName != "" {
		name = s.altLockName
	}
	return name
}

// getRemainingInterval returns the interval minus the remainder of dividing the time since state by
// the interval
func getRemainingInterval(start time.Time, interval time.Duration) time.Duration {
	if interval == 0 {
		return 0
	}
	return interval - (time.Since(start) % interval)
}

// clearTickerChannel performs a non-blocking select on the ticker channel
func clearTickerChannel(ticker *time.Ticker) {
	select {
	case <-ticker.C:
		// pull from ticker channel
	default:
		// ok
	}
}

// truncateSecondsWithFloor returns the result of truncating the duration to seconds and
// and applying a one second floor (e.g., 600ms becomes 1s, 1300ms becomes 2s, 1000ms becomes 2s)
func truncateSecondsWithFloor(d time.Duration) time.Duration {
	if d <= 1*time.Second {
		return 1 * time.Second
	}
	return d.Truncate(time.Second)
}
