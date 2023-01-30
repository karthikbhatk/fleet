package update

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

// RunnerOptions is options provided for the update runner.
type RunnerOptions struct {
	// CheckInterval is the interval to check for updates.
	CheckInterval time.Duration
	// Targets is the names of the artifacts to watch for updates.
	Targets []string
}

// Runner is a specialized runner for an Updater. It is designed with Execute and
// Interrupt functions to be compatible with oklog/run.
//
// It uses an Updater and makes sure to keep its targets up-to-date.
type Runner struct {
	updater     *Updater
	opt         RunnerOptions
	cancel      chan struct{}
	localHashes map[string][]byte
	mu          sync.Mutex
}

// UpdateRunnerOptTargets updates the RunnerOptions.Targets with the given target
func (r *Runner) UpdateRunnerOptTargets(target string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	// check if target already exists
	for _, t := range r.opt.Targets {
		if t == target {
			return
		}
	}
	r.opt.Targets = append(r.opt.Targets, target)
}

// NewRunner creates a new runner with the provided options. The runner must be
// started with Execute.
func NewRunner(updater *Updater, opt RunnerOptions) (*Runner, error) {
	if opt.CheckInterval <= 0 {
		return nil, errors.New("runner must be configured with interval greater than 0")
	}
	if len(opt.Targets) == 0 {
		return nil, errors.New("runner must have nonempty subscriptions")
	}

	// Initialize the hashes of the local files for all tracked targets.
	//
	// This is an optimization to not compute the hash of the local files every opt.CheckInterval
	// (knowing that they are not expected to change during the execution of the runner).
	localHashes := make(map[string][]byte)
	for _, target := range opt.Targets {
		meta, err := updater.Lookup(target)
		if err != nil {
			return nil, fmt.Errorf("target %s lookup: %w", target, err)
		}
		localTarget, err := updater.localTarget(target)
		if err != nil {
			return nil, fmt.Errorf("get local path for %s: %w", target, err)
		}
		switch _, localHash, err := fileHashes(meta, localTarget.Path); {
		case err == nil:
			localHashes[target] = localHash
			log.Info().Msgf("hash(%s)=%x", target, localHash)
		case errors.Is(err, os.ErrNotExist):
			// This is expected to happen if the target is not yet downloaded,
			// or if the user manually changed the target channel.
		default:
			return nil, fmt.Errorf("%s file hash: %w", target, err)
		}
	}

	return &Runner{
		updater: updater,
		opt:     opt,
		// chan gets capacity of 1 so we don't end up hung if Interrupt is
		// called after Execute has already returned.
		cancel:      make(chan struct{}, 1),
		localHashes: localHashes,
	}, nil
}

// Execute begins a loop checking for updates.
func (r *Runner) Execute() error {
	log.Debug().Msg("start updater")

	ticker := time.NewTicker(r.opt.CheckInterval)
	defer ticker.Stop()

	// Run until cancel or returning an error
	for {
		select {
		case <-r.cancel:
			return nil
		case <-ticker.C:
			didUpdate, err := r.UpdateAction()
			if err != nil {
				log.Info().Err(err).Msg("update failed")
			}
			if didUpdate {
				log.Info().Msg("exiting due to successful update")
				return nil
			}
		}
	}
}

// UpdateAction checks for updates on all targets.
// Returns true if one of the targets has been updated.
//
// NOTE: If it returns (true, non-nil error) then it means some target/s
// were successfully upgraded and some failed to upgrade.
func (r *Runner) UpdateAction() (bool, error) {
	if err := r.updater.UpdateMetadata(); err != nil {
		// Consider this a non-fatal error since it will be common to be offline
		// or otherwise unable to retrieve the metadata.
		return false, fmt.Errorf("update metadata: %w", err)
	}

	var didUpdate bool
	for _, target := range r.opt.Targets {
		meta, err := r.updater.Lookup(target)
		if err != nil {
			return didUpdate, fmt.Errorf("lookup failed: %w", err)
		}
		_, metaHash, err := selectHashFunction(meta)
		if err != nil {
			return didUpdate, fmt.Errorf("select hash for cache: %w", err)
		}

		// Check if we need to update the orbit symlink (e.g. if channel changed)
		needsSymlinkUpdate := false
		if target == "orbit" {
			var err error
			needsSymlinkUpdate, err = r.needsOrbitSymlinkUpdate()
			if err != nil {
				return false, fmt.Errorf("check symlink failed: %w", err)
			}
		}

		// Check whether the hash of the repository is different than
		// that of the target local file.
		if !bytes.Equal(r.localHashes[target], metaHash) || needsSymlinkUpdate {
			// Update detected
			log.Info().Str("target", target).Msg("update detected")
			if err := r.updateTarget(target); err != nil {
				return didUpdate, fmt.Errorf("update %s: %w", target, err)
			}
			log.Info().Str("target", target).Msg("update completed")
			didUpdate = true
			if target == "nudge" {
				// TODO handle this appropriately
				didUpdate = false
			}
		} else {
			log.Debug().Str("target", target).Msg("no update")
		}
	}

	return didUpdate, nil
}

func (r *Runner) needsOrbitSymlinkUpdate() (bool, error) {
	localTarget, err := r.updater.Get("orbit")
	if err != nil {
		return false, fmt.Errorf("get binary: %w", err)
	}
	path := localTarget.ExecPath

	// Symlink Orbit binary
	linkPath := filepath.Join(r.updater.opt.RootDirectory, "bin", "orbit", filepath.Base(path))

	existingPath, err := os.Readlink(linkPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return true, nil
		}
		return false, fmt.Errorf("read existing symlink: %w", err)
	}

	return existingPath != path, nil
}

func (r *Runner) updateTarget(target string) error {
	localTarget, err := r.updater.Get(target)
	if err != nil {
		return fmt.Errorf("get binary: %w", err)
	}
	path := localTarget.ExecPath

	if target != "orbit" {
		return nil
	}

	// Symlink Orbit binary
	linkPath := filepath.Join(r.updater.opt.RootDirectory, "bin", "orbit", filepath.Base(path))
	// Rename the old file otherwise overwrite fails
	if err := os.Rename(linkPath, linkPath+".old"); err != nil {
		return fmt.Errorf("move old symlink current: %w", err)
	}
	if err := os.Symlink(path, linkPath); err != nil {
		return fmt.Errorf("symlink current: %w", err)
	}

	return nil
}

func (r *Runner) Interrupt(err error) {
	r.cancel <- struct{}{}
	log.Debug().Err(err).Msg("interrupt updater")
}
