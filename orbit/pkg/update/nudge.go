package update

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/fleetdm/fleet/v4/orbit/pkg/constant"
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

	// UpdateRunner is the wrapped Runner where Nudge will be set as a target. It is responsible for
	// actually ensuring that Nudge is installed and updated via the designated TUF server.
	UpdateRunner *Runner

	// rootDir is where the Nudge configuration will be stored
	rootDir string

	// frequency is the minimum amount of time that must pass to launch
	// Nudge
	frequency time.Duration

	// ensures only one command runs at a time, protects access to lastRun
	cmdMu   sync.Mutex
	lastRun time.Time
}

func ApplyNudgeConfigFetcherMiddleware(
	f OrbitConfigFetcher,
	u *Runner,
	rootDir string,
	frequency time.Duration,
) OrbitConfigFetcher {
	return &NudgeConfigFetcher{Fetcher: f, UpdateRunner: u, rootDir: rootDir, frequency: frequency}
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
		log.Info().Msg("empty nudge config, removing nudge as target")
		n.UpdateRunner.RemoveRunnerOptTarget("nudge")
		return cfg, nil
	}

	if !n.UpdateRunner.HasRunnerOptTarget("nudge") {
		log.Info().Msg("adding nudge as target")
		n.UpdateRunner.AddRunnerOptTarget("nudge")
		return cfg, n.UpdateRunner.StoreLocalHash("nudge")
	}

	if err := n.manageNudgeConfig(*cfg.NudgeConfig); err != nil {
		log.Info().Err(err).Msg("nudge configuration")
		return cfg, err
	}

	if err := n.manageNudgeLaunch(); err != nil {
		log.Info().Err(err).Msg("nudge launch")
		return cfg, err
	}

	return cfg, nil
}

func (n *NudgeConfigFetcher) manageNudgeConfig(nudgeCfg fleet.NudgeConfig) error {
	jsonCfg, err := json.Marshal(nudgeCfg)
	if err != nil {
		return err
	}

	cfgFile := filepath.Join(n.rootDir, nudgeConfigFile)
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
	// has been tampered and and contains very large contents, we don't
	// want to load them into memory.
	if fileInfo.Size() != int64(len(jsonCfg)) {
		return writeConfig()
	}

	fileBytes, err := os.ReadFile(cfgFile)
	if err != nil {
		return err
	}

	if !bytes.Equal(fileBytes, jsonCfg) {
		return writeConfig()
	}

	return nil
}

func (n *NudgeConfigFetcher) manageNudgeLaunch() error {
	cfgFile := filepath.Join(n.rootDir, nudgeConfigFile)

	if n.cmdMu.TryLock() {
		defer n.cmdMu.Unlock()

		if time.Since(n.lastRun) > n.frequency {
			nudgePath, err := n.UpdateRunner.updater.ExecutableLocalPath("nudge")
			if err != nil {
				return err
			}

			cmd := exec.Command(nudgePath, "-json-url", fmt.Sprintf("file://%s", cfgFile))
			cmd.Stderr = os.Stderr
			cmd.Stdout = os.Stdout
			log.Info().Str("cmd", cmd.String()).Msg("start Nudge")
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("open path %q: %w", cfgFile, err)
			}

			n.lastRun = time.Now()
		}
	}

	return nil
}
