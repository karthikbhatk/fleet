package mdmlifecycle

import (
	"context"

	"github.com/fleetdm/fleet/v4/server/contexts/ctxerr"
	"github.com/fleetdm/fleet/v4/server/contexts/license"
	"github.com/fleetdm/fleet/v4/server/fleet"
	"github.com/fleetdm/fleet/v4/server/worker"
	kitlog "github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type HostAction string

const (
	HostActionTurnOn  HostAction = "turn-on"
	HostActionTurnOff HostAction = "turn-off"
	HostActionReset   HostAction = "reset"
	HostActionDelete  HostAction = "delete"
)

type HostOptions struct {
	Action          HostAction
	Platform        string
	UUID            string
	HardwareSerial  string
	HardwareModel   string
	EnrollReference string
	Host            *fleet.Host
}

// HostLifecycle manages MDM host lifecycle actions
type HostLifecycle struct {
	ds     fleet.Datastore
	logger kitlog.Logger
}

func New(ds fleet.Datastore, logger kitlog.Logger) *HostLifecycle {
	return &HostLifecycle{
		ds:     ds,
		logger: logger,
	}
}

// Do executes the provided MDMHostLifecycleAction based on the platform requested
func (t *HostLifecycle) Do(ctx context.Context, opts HostOptions) error {
	switch opts.Platform {
	case "darwin":
		err := t.doDarwin(ctx, opts)
		return ctxerr.Wrapf(ctx, err, "running darwin lifecycle action %s", opts.Action)
	case "windows":
		err := t.doWindows(ctx, opts)
		return ctxerr.Wrapf(ctx, err, "running windows lifecycle action %s", opts.Action)
	default:
		return ctxerr.Errorf(ctx, "unsupported platform %s", opts.Platform)
	}
}

func (t *HostLifecycle) doDarwin(ctx context.Context, opts HostOptions) error {
	switch opts.Action {
	case HostActionTurnOn:
		return t.darwinTurnOn(ctx, opts)

	case HostActionTurnOff:
		return t.uuidAction(ctx, t.ds.MDMAppleTurnOff, opts)

	case HostActionReset:
		return t.darwinReset(ctx, opts)

	case HostActionDelete:
		return t.darwinDelete(ctx, opts)

	default:
		return ctxerr.Errorf(ctx, "unknown action %s", opts.Action)

	}
}

func (t *HostLifecycle) doWindows(ctx context.Context, opts HostOptions) error {
	switch opts.Action {
	case HostActionReset, HostActionTurnOn:
		return t.uuidAction(ctx, t.ds.MDMWindowsResetEnrollment, opts)

	case HostActionDelete, HostActionTurnOff:
		return t.uuidAction(ctx, t.ds.MDMWindowsTurnOff, opts)

	default:
		return ctxerr.Errorf(ctx, "unknown action %s", opts.Action)
	}
}

type uuidFn func(context.Context, string) error

func (t *HostLifecycle) uuidAction(ctx context.Context, action uuidFn, opts HostOptions) error {
	if opts.UUID == "" {
		return ctxerr.New(ctx, "UUID option is required for this action")
	}

	return action(ctx, opts.UUID)
}

func (t *HostLifecycle) darwinReset(ctx context.Context, opts HostOptions) error {
	if opts.UUID == "" || opts.HardwareSerial == "" || opts.HardwareModel == "" {
		return ctxerr.New(ctx, "UUID, HardwareSerial and HardwareModel options are required for this action")
	}

	host := &fleet.Host{
		UUID:           opts.UUID,
		HardwareSerial: opts.HardwareSerial,
		HardwareModel:  opts.HardwareModel,
	}
	if err := t.ds.MDMAppleUpsertHost(ctx, host); err != nil {
		return ctxerr.Wrap(ctx, err, "upserting mdm host")
	}

	err := t.ds.MDMAppleResetEnrollment(ctx, opts.UUID)
	return ctxerr.Wrap(ctx, err, "reset mdm enrollment")
}

func (t *HostLifecycle) darwinTurnOn(ctx context.Context, opts HostOptions) error {
	if opts.UUID == "" {
		return ctxerr.New(ctx, "UUID option is required for this action")
	}

	nanoEnroll, err := t.ds.GetNanoMDMEnrollment(ctx, opts.UUID)
	if err != nil {
		return ctxerr.Wrap(ctx, err, "retrieving nano enrollment info")
	}

	if nanoEnroll == nil ||
		!nanoEnroll.Enabled ||
		nanoEnroll.Type != "Device" ||
		nanoEnroll.TokenUpdateTally != 1 {
		return nil
	}

	info, err := t.ds.GetHostMDMCheckinInfo(ctx, opts.UUID)
	if err != nil {
		return ctxerr.Wrap(ctx, err, "getting checkin info")
	}

	var tmID *uint
	if info.TeamID != 0 {
		tmID = &info.TeamID
	}

	// TODO: improve this to not enqueue the job if a host that is
	// assigned in ABM is manually enrolling for some reason.
	if info.DEPAssignedToFleet || info.InstalledFromDEP {
		t.logger.Log("info", "queueing post-enroll task for newly enrolled DEP device", "host_uuid", opts.UUID)
		if err := worker.QueueAppleMDMJob(
			ctx,
			t.ds,
			t.logger,
			worker.AppleMDMPostDEPEnrollmentTask,
			opts.UUID,
			tmID,
			opts.EnrollReference,
		); err != nil {
			return ctxerr.Wrap(ctx, err, "queue DEP post-enroll task")
		}
	}

	// manual MDM enrollments that are not fleet-enrolled yet
	if !info.InstalledFromDEP && !info.OsqueryEnrolled {
		if err := worker.QueueAppleMDMJob(
			ctx,
			t.ds,
			t.logger,
			worker.AppleMDMPostManualEnrollmentTask,
			opts.UUID,
			tmID,
			opts.EnrollReference,
		); err != nil {
			return ctxerr.Wrap(ctx, err, "queue manual post-enroll task")
		}
	}

	return nil
}

func (t *HostLifecycle) darwinDelete(ctx context.Context, opts HostOptions) error {
	if opts.Host == nil {
		return ctxerr.New(ctx, "a non-nil Host option is required to perform this action")
	}

	if !license.IsPremium(ctx) {
		// only premium tier supports DEP so nothing more to do
		return nil
	}

	ac, err := t.ds.AppConfig(ctx)
	if err != nil {
		return ctxerr.Wrap(ctx, err, "get app config")
	} else if !ac.MDM.AppleBMEnabledAndConfigured {
		// if ABM is not enabled and configured, nothing more to do
		return nil
	}

	dep, err := t.ds.GetHostDEPAssignment(ctx, opts.Host.ID)
	if err != nil && !fleet.IsNotFound(err) {
		return ctxerr.Wrap(ctx, err, "get host dep assignment")
	}

	if dep != nil && dep.DeletedAt == nil {
		return t.restorePendingDEPHost(ctx, opts.Host, ac)
	}

	// no DEP assignment was found or the DEP assignment was deleted in ABM
	// so nothing more to do
	return nil
}

func (t *HostLifecycle) restorePendingDEPHost(ctx context.Context, host *fleet.Host, appCfg *fleet.AppConfig) error {
	tmID, err := t.getConfigAppleBMDefaultTeamID(ctx, appCfg)
	if err != nil {
		return ctxerr.Wrap(ctx, err, "restore pending dep host")
	}
	host.TeamID = tmID

	if err := t.ds.RestoreMDMApplePendingDEPHost(ctx, host); err != nil {
		return ctxerr.Wrap(ctx, err, "restore pending dep host")
	}

	if _, err := worker.QueueMacosSetupAssistantJob(ctx, t.ds, t.logger,
		worker.MacosSetupAssistantHostsTransferred, tmID, host.HardwareSerial); err != nil {
		return ctxerr.Wrap(ctx, err, "queue macos setup assistant update profile job")
	}

	return nil
}

func (t *HostLifecycle) getConfigAppleBMDefaultTeamID(ctx context.Context, appCfg *fleet.AppConfig) (*uint, error) {
	var tmID *uint
	if name := appCfg.MDM.AppleBMDefaultTeam; name != "" {
		team, err := t.ds.TeamByName(ctx, name)
		switch {
		case fleet.IsNotFound(err):
			level.Debug(t.logger).Log(
				"msg",
				"unable to find default team assigned in config, mdm devices won't be assigned to a team",
				"team_name",
				name,
			)
			return nil, nil
		case err != nil:
			return nil, ctxerr.Wrap(ctx, err, "get default team for mdm devices")
		case team != nil:
			tmID = &team.ID
		}
	}

	return tmID, nil
}
