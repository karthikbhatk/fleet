package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	eewebhooks "github.com/fleetdm/fleet/v4/ee/server/webhooks"
	"github.com/fleetdm/fleet/v4/pkg/fleethttp"
	"github.com/fleetdm/fleet/v4/server"
	"github.com/fleetdm/fleet/v4/server/config"
	"github.com/fleetdm/fleet/v4/server/contexts/ctxerr"
	"github.com/fleetdm/fleet/v4/server/fleet"
	"github.com/fleetdm/fleet/v4/server/policies"
	"github.com/fleetdm/fleet/v4/server/service/externalsvc"
	"github.com/fleetdm/fleet/v4/server/service/schedule"
	"github.com/fleetdm/fleet/v4/server/vulnerabilities"
	"github.com/fleetdm/fleet/v4/server/vulnerabilities/msrc"
	"github.com/fleetdm/fleet/v4/server/vulnerabilities/oval"
	"github.com/fleetdm/fleet/v4/server/vulnerabilities/utils"
	"github.com/fleetdm/fleet/v4/server/webhooks"
	"github.com/fleetdm/fleet/v4/server/worker"
	"github.com/getsentry/sentry-go"
	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

func errHandler(ctx context.Context, logger kitlog.Logger, msg string, err error) {
	level.Error(logger).Log("msg", msg, "err", err)
	sentry.CaptureException(err)
	ctxerr.Handle(ctx, err)
}

func startVulnerabilitiesSchedule(
	ctx context.Context,
	instanceID string,
	ds fleet.Datastore,
	logger kitlog.Logger,
	config *config.VulnerabilitiesConfig,
	license *fleet.LicenseInfo,
) *schedule.Schedule {
	interval := config.Periodicity
	vulnerabilitiesLogger := kitlog.With(logger, "cron", "vulnerabilities")
	s := schedule.New(
		ctx, "vulnerabilities", instanceID, interval, ds,
		schedule.WithLogger(vulnerabilitiesLogger),
		schedule.WithJob(
			"cron_vulnerabilities",
			func(ctx context.Context) error {
				// TODO(lucas): Decouple cronVulnerabilities into multiple jobs.
				return cronVulnerabilities(ctx, ds, vulnerabilitiesLogger, config, license)
			},
		),
		schedule.WithJob(
			"cron_sync_host_software",
			func(ctx context.Context) error {
				return ds.SyncHostsSoftware(ctx, time.Now())
			},
		),
	)
	s.Start()
	return s
}

func cronVulnerabilities(
	ctx context.Context,
	ds fleet.Datastore,
	logger kitlog.Logger,
	config *config.VulnerabilitiesConfig,
	license *fleet.LicenseInfo,
) error {
	if config.CurrentInstanceChecks == "no" || config.CurrentInstanceChecks == "0" {
		level.Info(logger).Log("msg", "host not configured to check for vulnerabilities")
		return nil
	}

	level.Info(logger).Log("periodicity", config.Periodicity)

	appConfig, err := ds.AppConfig(ctx)
	if err != nil {
		return fmt.Errorf("reading app config: %w", err)
	}

	if !appConfig.Features.EnableSoftwareInventory {
		level.Info(logger).Log("msg", "software inventory not configured")
		return nil
	}

	var vulnPath string
	switch {
	case config.DatabasesPath != "" && appConfig.VulnerabilitySettings.DatabasesPath != "":
		vulnPath = config.DatabasesPath
		level.Info(logger).Log(
			"msg", "fleet config takes precedence over app config when both are configured",
			"databases_path", vulnPath,
		)
	case config.DatabasesPath != "":
		vulnPath = config.DatabasesPath
	case appConfig.VulnerabilitySettings.DatabasesPath != "":
		vulnPath = appConfig.VulnerabilitySettings.DatabasesPath
	default:
		level.Info(logger).Log("msg", "vulnerability scanning not configured, vulnerabilities databases path is empty")
	}
	if vulnPath != "" {
		level.Info(logger).Log("msg", "scanning vulnerabilities")
		if err := scanVulnerabilities(ctx, ds, logger, config, appConfig, vulnPath, license); err != nil {
			return fmt.Errorf("scanning vulnerabilities: %w", err)
		}
	}

	return nil
}

func scanVulnerabilities(
	ctx context.Context,
	ds fleet.Datastore,
	logger kitlog.Logger,
	config *config.VulnerabilitiesConfig,
	appConfig *fleet.AppConfig,
	vulnPath string,
	license *fleet.LicenseInfo,
) error {
	level.Debug(logger).Log("msg", "creating vulnerabilities databases path", "databases_path", vulnPath)
	err := os.MkdirAll(vulnPath, 0o755)
	if err != nil {
		return fmt.Errorf("create vulnerabilities databases directory: %w", err)
	}

	var vulnAutomationEnabled string

	// only one vuln automation (i.e. webhook or integration) can be enabled at a
	// time, enforced when updating the appconfig.
	if appConfig.WebhookSettings.VulnerabilitiesWebhook.Enable {
		vulnAutomationEnabled = "webhook"
	}

	// check for jira integrations
	for _, j := range appConfig.Integrations.Jira {
		if j.EnableSoftwareVulnerabilities {
			if vulnAutomationEnabled != "" {
				err := ctxerr.New(ctx, "jira check")
				errHandler(ctx, logger, "more than one automation enabled", err)
			}
			vulnAutomationEnabled = "jira"
			break
		}
	}

	// check for Zendesk integrations
	for _, z := range appConfig.Integrations.Zendesk {
		if z.EnableSoftwareVulnerabilities {
			if vulnAutomationEnabled != "" {
				err := ctxerr.New(ctx, "zendesk check")
				errHandler(ctx, logger, "more than one automation enabled", err)
			}
			vulnAutomationEnabled = "zendesk"
			break
		}
	}

	level.Debug(logger).Log("vulnAutomationEnabled", vulnAutomationEnabled)

	nvdVulns := checkNVDVulnerabilities(ctx, ds, logger, vulnPath, config, vulnAutomationEnabled != "")
	ovalVulns := checkOvalVulnerabilities(ctx, ds, logger, vulnPath, config, vulnAutomationEnabled != "")
	checkWinVulnerabilities(ctx, ds, logger, vulnPath, config, vulnAutomationEnabled != "")

	// If no automations enabled, then there is nothing else to do...
	if vulnAutomationEnabled == "" {
		return nil
	}

	vulns := make([]fleet.SoftwareVulnerability, len(nvdVulns)+len(ovalVulns))
	vulns = append(vulns, nvdVulns...)
	vulns = append(vulns, ovalVulns...)

	meta, err := ds.ListCVEs(ctx, config.RecentVulnerabilityMaxAge)
	if err != nil {
		errHandler(ctx, logger, "could not fetch CVE meta", err)
		return nil
	}

	recentV, matchingMeta := utils.RecentVulns(vulns, meta)

	if len(recentV) > 0 {
		switch vulnAutomationEnabled {
		case "webhook":
			args := webhooks.VulnArgs{
				Vulnerablities: recentV,
				Meta:           matchingMeta,
				AppConfig:      appConfig,
				Time:           time.Now(),
			}
			mapper := webhooks.NewMapper()
			if license.IsPremium() {
				mapper = eewebhooks.NewMapper()
			}
			// send recent vulnerabilities via webhook
			if err := webhooks.TriggerVulnerabilitiesWebhook(
				ctx,
				ds,
				kitlog.With(logger, "webhook", "vulnerabilities"),
				args,
				mapper,
			); err != nil {
				errHandler(ctx, logger, "triggering vulnerabilities webhook", err)
			}

		case "jira":
			// queue job to create jira issues
			if err := worker.QueueJiraVulnJobs(
				ctx,
				ds,
				kitlog.With(logger, "jira", "vulnerabilities"),
				recentV,
			); err != nil {
				errHandler(ctx, logger, "queueing vulnerabilities to jira", err)
			}

		case "zendesk":
			// queue job to create zendesk ticket
			if err := worker.QueueZendeskVulnJobs(
				ctx,
				ds,
				kitlog.With(logger, "zendesk", "vulnerabilities"),
				recentV,
			); err != nil {
				errHandler(ctx, logger, "queueing vulnerabilities to Zendesk", err)
			}

		default:
			err = ctxerr.New(ctx, "no vuln automations enabled")
			errHandler(ctx, logger, "attempting to process vuln automations", err)
		}
	}

	return nil
}

func checkWinVulnerabilities(
	ctx context.Context,
	ds fleet.Datastore,
	logger kitlog.Logger,
	vulnPath string,
	config *config.VulnerabilitiesConfig,
	collectVulns bool,
) []fleet.OSVulnerability {
	var results []fleet.OSVulnerability

	// Get OS
	os, err := ds.ListOperatingSystems(ctx)
	if err != nil {
		errHandler(ctx, logger, "fetching list of operating systems", err)
		return nil
	}

	if !config.DisableDataSync {
		// Sync MSRC definitions
		client := fleethttp.NewClient()
		err = msrc.Sync(ctx, client, vulnPath, os)
		if err != nil {
			errHandler(ctx, logger, "updating msrc definitions", err)
		}
	}

	// Analyze all Win OS using the synched MSRC artifact.
	if !config.DisableWinOSVulnerabilities {
		for _, o := range os {
			start := time.Now()
			r, err := msrc.Analyze(ctx, ds, o, vulnPath, collectVulns)
			elapsed := time.Since(start)
			level.Debug(logger).Log(
				"msg", "msrc-analysis-done",
				"os name", o.Name,
				"os version", o.Version,
				"elapsed", elapsed,
				"found new", len(r))
			results = append(results, r...)
			if err != nil {
				errHandler(ctx, logger, "analyzing hosts for Windows vulnerabilities", err)
			}
		}
	}

	return results
}

func checkOvalVulnerabilities(
	ctx context.Context,
	ds fleet.Datastore,
	logger kitlog.Logger,
	vulnPath string,
	config *config.VulnerabilitiesConfig,
	collectVulns bool,
) []fleet.SoftwareVulnerability {
	if config.DisableDataSync {
		return nil
	}

	var results []fleet.SoftwareVulnerability

	// Get Platforms
	versions, err := ds.OSVersions(ctx, nil, nil, nil, nil)
	if err != nil {
		errHandler(ctx, logger, "updating oval definitions", err)
		return nil
	}

	// Sync on disk OVAL definitions with current OS Versions.
	client := fleethttp.NewClient()
	downloaded, err := oval.Refresh(ctx, client, versions, vulnPath)
	if err != nil {
		errHandler(ctx, logger, "updating oval definitions", err)
	}
	for _, d := range downloaded {
		level.Debug(logger).Log("oval-sync-downloaded", d)
	}

	// Analyze all supported os versions using the synched OVAL definitions.
	for _, version := range versions.OSVersions {
		start := time.Now()
		r, err := oval.Analyze(ctx, ds, version, vulnPath, collectVulns)
		elapsed := time.Since(start)
		level.Debug(logger).Log(
			"msg", "oval-analysis-done",
			"platform", version.Name,
			"elapsed", elapsed,
			"found new", len(r))
		results = append(results, r...)
		if err != nil {
			errHandler(ctx, logger, "analyzing oval definitions", err)
		}
	}

	return results
}

func checkNVDVulnerabilities(
	ctx context.Context,
	ds fleet.Datastore,
	logger kitlog.Logger,
	vulnPath string,
	config *config.VulnerabilitiesConfig,
	collectVulns bool,
) []fleet.SoftwareVulnerability {
	if !config.DisableDataSync {
		opts := vulnerabilities.SyncOptions{
			VulnPath:           config.DatabasesPath,
			CPEDBURL:           config.CPEDatabaseURL,
			CPETranslationsURL: config.CPETranslationsURL,
			CVEFeedPrefixURL:   config.CVEFeedPrefixURL,
		}
		err := vulnerabilities.Sync(opts)
		if err != nil {
			errHandler(ctx, logger, "syncing vulnerability database", err)
			return nil
		}
	}

	if err := vulnerabilities.LoadCVEMeta(logger, vulnPath, ds); err != nil {
		errHandler(ctx, logger, "load cve meta", err)
		// don't return, continue on ...
	}

	err := vulnerabilities.TranslateSoftwareToCPE(ctx, ds, vulnPath, logger)
	if err != nil {
		errHandler(ctx, logger, "analyzing vulnerable software: Software->CPE", err)
		return nil
	}

	vulns, err := vulnerabilities.TranslateCPEToCVE(ctx, ds, vulnPath, logger, collectVulns)
	if err != nil {
		errHandler(ctx, logger, "analyzing vulnerable software: CPE->CVE", err)
		return nil
	}

	return vulns
}

func startAutomationsSchedule(
	ctx context.Context,
	instanceID string,
	ds fleet.Datastore,
	logger kitlog.Logger,
	intervalReload time.Duration,
	failingPoliciesSet fleet.FailingPolicySet,
) (*schedule.Schedule, error) {
	const defaultAutomationsInterval = 24 * time.Hour

	appConfig, err := ds.AppConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting app config: %w", err)
	}
	s := schedule.New(
		// TODO(sarah): Reconfigure settings so automations interval doesn't reside under webhook settings
		ctx, "automations", instanceID, appConfig.WebhookSettings.Interval.ValueOr(defaultAutomationsInterval), ds,
		schedule.WithLogger(kitlog.With(logger, "cron", "automations")),
		schedule.WithConfigReloadInterval(intervalReload, func(ctx context.Context) (time.Duration, error) {
			appConfig, err := ds.AppConfig(ctx)
			if err != nil {
				return 0, err
			}
			newInterval := appConfig.WebhookSettings.Interval.ValueOr(defaultAutomationsInterval)
			return newInterval, nil
		}),
		schedule.WithJob(
			"host_status_webhook",
			func(ctx context.Context) error {
				return webhooks.TriggerHostStatusWebhook(
					ctx, ds, kitlog.With(logger, "automation", "host_status"),
				)
			},
		),
		schedule.WithJob(
			"failing_policies_automation",
			func(ctx context.Context) error {
				return triggerFailingPoliciesAutomation(ctx, ds, kitlog.With(logger, "automation", "failing_policies"), failingPoliciesSet)
			},
		),
	)
	s.Start()
	return s, nil
}

func triggerFailingPoliciesAutomation(
	ctx context.Context,
	ds fleet.Datastore,
	logger kitlog.Logger,
	failingPoliciesSet fleet.FailingPolicySet,
) error {
	appConfig, err := ds.AppConfig(ctx)
	if err != nil {
		return fmt.Errorf("getting app config: %w", err)
	}
	serverURL, err := url.Parse(appConfig.ServerSettings.ServerURL)
	if err != nil {
		return fmt.Errorf("parsing appConfig.ServerSettings.ServerURL: %w", err)
	}

	err = policies.TriggerFailingPoliciesAutomation(ctx, ds, logger, failingPoliciesSet, func(policy *fleet.Policy, cfg policies.FailingPolicyAutomationConfig) error {
		switch cfg.AutomationType {
		case policies.FailingPolicyWebhook:
			return webhooks.SendFailingPoliciesBatchedPOSTs(
				ctx, policy, failingPoliciesSet, cfg.HostBatchSize, serverURL, cfg.WebhookURL, time.Now(), logger)

		case policies.FailingPolicyJira:
			hosts, err := failingPoliciesSet.ListHosts(policy.ID)
			if err != nil {
				return ctxerr.Wrapf(ctx, err, "listing hosts for failing policies set %d", policy.ID)
			}
			if err := worker.QueueJiraFailingPolicyJob(ctx, ds, logger, policy, hosts); err != nil {
				return err
			}
			if err := failingPoliciesSet.RemoveHosts(policy.ID, hosts); err != nil {
				return ctxerr.Wrapf(ctx, err, "removing %d hosts from failing policies set %d", len(hosts), policy.ID)
			}

		case policies.FailingPolicyZendesk:
			hosts, err := failingPoliciesSet.ListHosts(policy.ID)
			if err != nil {
				return ctxerr.Wrapf(ctx, err, "listing hosts for failing policies set %d", policy.ID)
			}
			if err := worker.QueueZendeskFailingPolicyJob(ctx, ds, logger, policy, hosts); err != nil {
				return err
			}
			if err := failingPoliciesSet.RemoveHosts(policy.ID, hosts); err != nil {
				return ctxerr.Wrapf(ctx, err, "removing %d hosts from failing policies set %d", len(hosts), policy.ID)
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("triggering failing policies automation: %w", err)
	}

	return nil
}

func cronWorker(
	ctx context.Context,
	ds fleet.Datastore,
	logger kitlog.Logger,
	identifier string,
) {
	const (
		lockDuration        = 10 * time.Minute
		lockAttemptInterval = 10 * time.Minute
	)

	logger = kitlog.With(logger, "cron", lockKeyWorker)

	// create the worker and register the Jira and Zendesk jobs even if no
	// integration is enabled, as that config can change live (and if it's not
	// there won't be any records to process so it will mostly just sleep).
	w := worker.NewWorker(ds, logger)
	jira := &worker.Jira{
		Datastore:     ds,
		Log:           logger,
		NewClientFunc: newJiraClient,
	}
	zendesk := &worker.Zendesk{
		Datastore:     ds,
		Log:           logger,
		NewClientFunc: newZendeskClient,
	}
	// leave the url empty for now, will be filled when the lock is acquired with
	// the up-to-date config.
	w.Register(jira)
	w.Register(zendesk)

	// Read app config a first time before starting, to clear up any failer client
	// configuration if we're not on a fleet-owned server. Technically, the ServerURL
	// could change dynamically, but for the needs of forced client failures, this
	// is not a possible scenario.
	appConfig, err := ds.AppConfig(ctx)
	if err != nil {
		errHandler(ctx, logger, "couldn't read app config", err)
	}
	// we clear it even if we fail to load the app config, not a likely scenario
	// in our test environments for the needs of forced failures.
	if !strings.Contains(appConfig.ServerSettings.ServerURL, "fleetdm") {
		os.Unsetenv("FLEET_JIRA_CLIENT_FORCED_FAILURES")
		os.Unsetenv("FLEET_ZENDESK_CLIENT_FORCED_FAILURES")
	}

	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			level.Debug(logger).Log("waiting", "done")
			ticker.Reset(lockAttemptInterval)
		case <-ctx.Done():
			level.Debug(logger).Log("exit", "done with cron.")
			return
		}

		if locked, err := ds.Lock(ctx, lockKeyWorker, identifier, lockDuration); err != nil {
			level.Error(logger).Log("msg", "Error acquiring lock", "err", err)
			continue
		} else if !locked {
			level.Debug(logger).Log("msg", "Not the leader. Skipping...")
			continue
		}

		// Read app config to be able to use the latest configuration for the Jira
		// integration.
		appConfig, err := ds.AppConfig(ctx)
		if err != nil {
			errHandler(ctx, logger, "couldn't read app config", err)
			continue
		}

		jira.FleetURL = appConfig.ServerSettings.ServerURL
		zendesk.FleetURL = appConfig.ServerSettings.ServerURL

		workCtx, cancel := context.WithTimeout(ctx, lockDuration)
		if err := w.ProcessJobs(workCtx); err != nil {
			errHandler(ctx, logger, "Error processing jobs", err)
		}
		cancel() // don't use defer inside loop
	}
}

func newJiraClient(opts *externalsvc.JiraOptions) (worker.JiraClient, error) {
	client, err := externalsvc.NewJiraClient(opts)
	if err != nil {
		return nil, err
	}

	// create client wrappers to introduce forced failures if configured
	// to do so via the environment variable.
	// format is "<modulo number>;<cve1>,<cve2>,<cve3>,..."
	failerClient := newFailerClient(os.Getenv("FLEET_JIRA_CLIENT_FORCED_FAILURES"))
	if failerClient != nil {
		failerClient.JiraClient = client
		return failerClient, nil
	}
	return client, nil
}

func newZendeskClient(opts *externalsvc.ZendeskOptions) (worker.ZendeskClient, error) {
	client, err := externalsvc.NewZendeskClient(opts)
	if err != nil {
		return nil, err
	}

	// create client wrappers to introduce forced failures if configured
	// to do so via the environment variable.
	// format is "<modulo number>;<cve1>,<cve2>,<cve3>,..."
	failerClient := newFailerClient(os.Getenv("FLEET_ZENDESK_CLIENT_FORCED_FAILURES"))
	if failerClient != nil {
		failerClient.ZendeskClient = client
		return failerClient, nil
	}
	return client, nil
}

func newFailerClient(forcedFailures string) *worker.TestAutomationFailer {
	var failerClient *worker.TestAutomationFailer
	if forcedFailures != "" {

		parts := strings.Split(forcedFailures, ";")
		if len(parts) == 2 {
			mod, _ := strconv.Atoi(parts[0])
			cves := strings.Split(parts[1], ",")
			if mod > 0 || len(cves) > 0 {
				failerClient = &worker.TestAutomationFailer{
					FailCallCountModulo: mod,
					AlwaysFailCVEs:      cves,
				}
			}
		}
	}
	return failerClient
}

func startCleanupsAndAggregationSchedule(
	ctx context.Context, instanceID string, ds fleet.Datastore, logger kitlog.Logger, enrollHostLimiter fleet.EnrollHostLimiter,
) {
	schedule.New(
		ctx, "cleanups_then_aggregation", instanceID, 1*time.Hour, ds,
		// Using leader for the lock to be backwards compatilibity with old deployments.
		schedule.WithAltLockID("leader"),
		schedule.WithLogger(kitlog.With(logger, "cron", "cleanups_then_aggregation")),
		// Run cleanup jobs first.
		schedule.WithJob(
			"distributed_query_campaigns",
			func(ctx context.Context) error {
				_, err := ds.CleanupDistributedQueryCampaigns(ctx, time.Now())
				return err
			},
		),
		schedule.WithJob(
			"incoming_hosts",
			func(ctx context.Context) error {
				_, err := ds.CleanupIncomingHosts(ctx, time.Now())
				return err
			},
		),
		schedule.WithJob(
			"carves",
			func(ctx context.Context) error {
				_, err := ds.CleanupCarves(ctx, time.Now())
				return err
			},
		),
		schedule.WithJob(
			"expired_hosts",
			func(ctx context.Context) error {
				_, err := ds.CleanupExpiredHosts(ctx)
				return err
			},
		),
		schedule.WithJob(
			"policy_membership",
			func(ctx context.Context) error {
				return ds.CleanupPolicyMembership(ctx, time.Now())
			},
		),
		schedule.WithJob(
			"sync_enrolled_host_ids",
			func(ctx context.Context) error {
				return enrollHostLimiter.SyncEnrolledHostIDs(ctx)
			},
		),
		schedule.WithJob(
			"cleanup_host_operating_systems",
			func(ctx context.Context) error {
				return ds.CleanupHostOperatingSystems(ctx)
			},
		),
		// Run aggregation jobs after cleanups.
		schedule.WithJob(
			"query_aggregated_stats",
			func(ctx context.Context) error {
				return ds.UpdateQueryAggregatedStats(ctx)
			},
		),
		schedule.WithJob(
			"scheduled_query_aggregated_stats",
			func(ctx context.Context) error {
				return ds.UpdateScheduledQueryAggregatedStats(ctx)
			},
		),
		schedule.WithJob(
			"aggregated_munki_and_mdm",
			func(ctx context.Context) error {
				return ds.GenerateAggregatedMunkiAndMDM(ctx)
			},
		),
		schedule.WithJob(
			"update_os_versions",
			func(ctx context.Context) error {
				return ds.UpdateOSVersions(ctx)
			},
		),
	).Start()
}

func startSendStatsSchedule(ctx context.Context, instanceID string, ds fleet.Datastore, config config.FleetConfig, license *fleet.LicenseInfo, logger kitlog.Logger) {
	schedule.New(
		ctx, "stats", instanceID, 1*time.Hour, ds,
		schedule.WithLogger(kitlog.With(logger, "cron", "stats")),
		schedule.WithJob(
			"try_send_statistics",
			func(ctx context.Context) error {
				// NOTE(mna): this is not a route from the fleet server (not in server/service/handler.go) so it
				// will not automatically support the /latest/ versioning. Leaving it as /v1/ for that reason.
				return trySendStatistics(ctx, ds, fleet.StatisticsFrequency, "https://fleetdm.com/api/v1/webhooks/receive-usage-analytics", config, license)
			},
		),
	).Start()
}

func trySendStatistics(ctx context.Context, ds fleet.Datastore, frequency time.Duration, url string, config config.FleetConfig, license *fleet.LicenseInfo) error {
	ac, err := ds.AppConfig(ctx)
	if err != nil {
		return err
	}
	if !ac.ServerSettings.EnableAnalytics {
		return nil
	}

	stats, shouldSend, err := ds.ShouldSendStatistics(ctx, frequency, config, license)
	if err != nil {
		return err
	}
	if !shouldSend {
		return nil
	}

	err = server.PostJSONWithTimeout(ctx, url, stats)
	if err != nil {
		return err
	}
	return ds.RecordStatisticsSent(ctx)
}
