package update

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fleetdm/fleet/v4/orbit/pkg/constant"
	"github.com/fleetdm/fleet/v4/orbit/pkg/execuser"
	"github.com/fleetdm/fleet/v4/server/fleet"
	"github.com/rs/zerolog/log"
)

const nudgeConfigFile = "nudge-config.json"

// NudgeConfigFetcher is a kind of middleware that wraps an OrbitConfigFetcher and a Runner.
// It checks the config supplied by the wrapped OrbitConfigFetcher to detects whether the Fleet
// server has supplied a Nudge config. If so, it sets Nudge as a target on the wrapped Runner.
type NudgeConfigFetcher struct {
	// Fetcher is the OrbitConfigFetcher that will be wrapped. It is responsible
	// for actually returning the orbit configuration or an error.
	Fetcher OrbitConfigFetcher
	opt     NudgeConfigFetcherOptions
	// ensures only one command runs at a time, protects access to lastRun
	cmdMu   sync.Mutex
	lastRun time.Time
}

type NudgeConfigFetcherOptions struct {
	// UpdateRunner is the wrapped Runner where Nudge will be set as a target. It is responsible for
	// actually ensuring that Nudge is installed and updated via the designated TUF server.
	UpdateRunner *Runner
	// RootDir is where the Nudge configuration will be stored
	RootDir string
	// Interval is the minimum amount of time that must pass to launch
	// Nudge
	Interval time.Duration
	// runNudgeFn can be set in tests to mock the command executed to
	// run Nudge.
	runNudgeFn func(execPath, configPath string) error
}

func ApplyNudgeConfigFetcherMiddleware(f OrbitConfigFetcher, opt NudgeConfigFetcherOptions) OrbitConfigFetcher {
	return &NudgeConfigFetcher{Fetcher: f, opt: opt}
}

// GetConfig calls the wrapped Fetcher's GetConfig method, and detects if the
// Fleet server has supplied a Nudge config.
//
// If a Nudge config is supplied, it:
//
// - ensures that Nudge is installed and updated via the designated TUF server.
// - ensures that Nudge is opened at an interval given by n.frequency with the
// provided config.
func (n *NudgeConfigFetcher) GetConfig() (*fleet.OrbitConfig, error) {
	log.Debug().Msg("running nudge config fetcher middleware")
	cfg, err := n.Fetcher.GetConfig()
	if err != nil {
		log.Info().Err(err).Msg("calling GetConfig from NudgeConfigFetcher")
		return nil, err
	}

	if cfg == nil {
		log.Debug().Msg("NudgeConfigFetcher received nil config")
		return nil, nil
	}

	if cfg.NudgeConfig == nil {
		log.Debug().Msg("empty nudge config, removing nudge as target")
		// TODO(roberto): by early returning and removing the target from the
		// runner/updater we ensure Nudge won't be opened/updated again
		// but we don't actually remove the file from disk. We
		// knowingly decided to do this as a post MVP optimization.
		n.opt.UpdateRunner.RemoveRunnerOptTarget("nudge")
		n.opt.UpdateRunner.updater.RemoveTargetInfo("nudge")
		return cfg, nil
	}

	updaterHasTarget := n.opt.UpdateRunner.HasRunnerOptTarget("nudge")
	runnerHasLocalHash := n.opt.UpdateRunner.HasLocalHash("nudge")
	if !updaterHasTarget || !runnerHasLocalHash {
		log.Info().Msg("refreshing the update runner config with Nudge targets and hashes")
		log.Debug().Msgf("updater has target: %t, runner has local hash: %t", updaterHasTarget, runnerHasLocalHash)
		return cfg, n.setTargetsAndHashes()
	}

	if err := n.configure(*cfg.NudgeConfig); err != nil {
		log.Info().Err(err).Msg("nudge configuration")
		return cfg, err
	}

	if err := n.launch(); err != nil {
		log.Info().Err(err).Msg("nudge launch")
		return cfg, err
	}

	return cfg, nil
}

func (n *NudgeConfigFetcher) setTargetsAndHashes() error {
	n.opt.UpdateRunner.AddRunnerOptTarget("nudge")
	n.opt.UpdateRunner.updater.SetTargetInfo("nudge", NudgeMacOSTarget)
	// we don't want to keep nudge as a target if we failed to update the
	// cached hashes in the runner.
	if err := n.opt.UpdateRunner.StoreLocalHash("nudge"); err != nil {
		log.Debug().Msgf("removing nudge from target options, error updating local hashes: %e", err)
		n.opt.UpdateRunner.RemoveRunnerOptTarget("nudge")
		n.opt.UpdateRunner.updater.RemoveTargetInfo("nudge")
		return err
	}
	return nil
}

func (n *NudgeConfigFetcher) configure(nudgeCfg fleet.NudgeConfig) error {
	jsonCfg, err := json.Marshal(nudgeCfg)
	if err != nil {
		return err
	}

	cfgFile := filepath.Join(n.opt.RootDir, nudgeConfigFile)
	writeConfig := func() error {
		return os.WriteFile(cfgFile, jsonCfg, constant.DefaultWorldReadableFileMode)
	}

	fileInfo, err := os.Stat(cfgFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return writeConfig()
		}
		return err
	}

	// this not only an optimization, but mostly a safeguard: if the file
	// has been tampered and contains very large contents, we don't
	// want to load them into memory.
	if fileInfo.Size() != int64(len(jsonCfg)) {
		log.Debug().Msg("configuring nudge: local file has different size than remote, writing remote config")
		return writeConfig()
	}

	fileBytes, err := os.ReadFile(cfgFile)
	if err != nil {
		return err
	}

	if !bytes.Equal(fileBytes, jsonCfg) {
		log.Debug().Msg("configuring nudge: local file is different than remote, writing remote config")
		return writeConfig()
	}

	return nil
}

func (n *NudgeConfigFetcher) launch() error {
	cfgFile := filepath.Join(n.opt.RootDir, nudgeConfigFile)

	if n.cmdMu.TryLock() {
		defer n.cmdMu.Unlock()

		if time.Since(n.lastRun) > n.opt.Interval {
			nudge, err := n.opt.UpdateRunner.updater.localTarget("nudge")
			if err != nil {
				return err
			}

			// before moving forward, check that the file at the
			// path is the file we're about to open hasn't been
			// tampered with.
			meta, err := n.opt.UpdateRunner.updater.Lookup("nudge")
			if err != nil {
				return err
			}
			// if we can't find the file, or the hash doesn't match
			// make sure nudge is added as a target and the hashes
			// are refreshed
			if err := checkFileHash(meta, nudge.Path); err != nil {
				return n.setTargetsAndHashes()
			}

			fn := n.opt.runNudgeFn
			if fn == nil {
				fn = func(appPath, configPath string) error {
					// TODO(roberto): when an user selects "Later" from the
					// Nudge defer menu, the Nudge UI will be shown the
					// next time Nudge is launched. If for some reason orbit
					// restarts (eg: an update) and the user has a pending
					// OS update, we might show Nudge more than one time
					// every n.frequency.
					//
					// Note that this only happens for the "Later" option,
					// all other options behave as expected and Nudge will
					// respect the time chosen (eg: next day) and it won't
					// show up even if it's opened multiple times in that
					// interval.
					log.Info().Msg("running Nudge")
					return execuser.Run(
						appPath,
						execuser.WithArg("-json-url", configPath),
					)
				}
			}

			if err := fn(nudge.DirPath, fmt.Sprintf("file://%s", cfgFile)); err != nil {
				return fmt.Errorf("opening Nudge with config %q: %w", cfgFile, err)
			}

			n.lastRun = time.Now()
		}
	}

	return nil
}
