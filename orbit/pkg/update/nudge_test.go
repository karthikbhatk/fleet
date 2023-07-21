package update

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/fleetdm/fleet/v4/orbit/pkg/constant"
	"github.com/fleetdm/fleet/v4/pkg/optjson"
	"github.com/fleetdm/fleet/v4/server/fleet"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestNudge(t *testing.T) {
	testingSuite := new(nudgeTestSuite)
	testingSuite.s = &testingSuite.Suite
	suite.Run(t, testingSuite)
}

type nudgeTestSuite struct {
	suite.Suite
	withTUF
}

func (s *nudgeTestSuite) TestUpdatesDisabled() {
	t := s.T()
	var err error
	cfg := &fleet.OrbitConfig{}
	cfg.NudgeConfig, err = fleet.NewNudgeConfig(fleet.MacOSUpdates{MinimumVersion: optjson.SetString("11"), Deadline: optjson.SetString("2022-01-04")})
	require.NoError(t, err)
	runNudgeFn := func(execPath, configPath string) error {
		return nil
	}
	var f OrbitConfigFetcher = &dummyConfigFetcher{cfg: cfg}
	f = ApplyNudgeConfigFetcherMiddleware(f, NudgeConfigFetcherOptions{
		UpdateRunner: nil,
		RootDir:      t.TempDir(),
		Interval:     time.Minute,
		runNudgeFn:   runNudgeFn,
	})

	// we used to get a panic if updates were disabled (see #11980)
	gotCfg, err := f.GetConfig()
	require.NoError(t, err)
	require.Equal(t, cfg, gotCfg)
}

func (s *nudgeTestSuite) TestNudgeConfigFetcherAddNudge() {
	t := s.T()
	tmpDir := t.TempDir()
	updater := &Updater{
		client: s.client,
		opt:    Options{Targets: make(map[string]TargetInfo), RootDirectory: tmpDir},
	}
	runner := &Runner{updater: updater, localHashes: make(map[string][]byte)}
	interval := time.Second
	cfg := &fleet.OrbitConfig{}
	nudgePath := "nudge/macos/stable/nudge.app.tar.gz"

	// set up mock runNudgeFn to capture exec command
	var execCmd func(command string, args ...string) *exec.Cmd
	var execOut string
	runNudgeFnInvoked := false
	runNudgeFn := func(execPath, configPath string) error {
		runNudgeFnInvoked = true
		if execCmd != nil {
			cmd := execCmd(execPath, configPath)
			out, err := cmd.Output()
			if err != nil {
				return err
			}
			execOut = string(out)
		}
		return nil
	}

	var f OrbitConfigFetcher = &dummyConfigFetcher{cfg: cfg}
	f = ApplyNudgeConfigFetcherMiddleware(f, NudgeConfigFetcherOptions{
		UpdateRunner: runner,
		RootDir:      tmpDir,
		Interval:     interval,
		runNudgeFn:   runNudgeFn,
	})
	configPath := filepath.Join(tmpDir, nudgeConfigFile)

	// nudge is not added to targets if nudge config is not present
	cfg.NudgeConfig = nil
	gotCfg, err := f.GetConfig()
	require.NoError(t, err)
	require.Equal(t, cfg, gotCfg)
	targets := runner.updater.opt.Targets
	require.Len(t, targets, 0)

	// set the config
	cfg.NudgeConfig, err = fleet.NewNudgeConfig(fleet.MacOSUpdates{MinimumVersion: optjson.SetString("11"), Deadline: optjson.SetString("2022-01-04")})
	require.NoError(t, err)

	// there's an error when the remote repo doesn't have the target yet
	gotCfg, err = f.GetConfig()
	require.ErrorContains(t, err, "tuf: file not found")
	require.Equal(t, cfg, gotCfg)

	// add nuge to the remote
	s.addRemoteTarget(nudgePath)

	// nothing happens if a nil runner is provided

	// nudge is added to targets when nudge config is present
	gotCfg, err = f.GetConfig()
	require.NoError(t, err)
	require.Equal(t, cfg, gotCfg)
	targets = runner.updater.opt.Targets
	require.Len(t, targets, 1)
	ti, ok := targets["nudge"]
	require.True(t, ok)
	require.EqualValues(t, NudgeMacOSTarget, ti)

	// override the custom check since we don't really have an executable
	ti.CustomCheckExec = func(path string) error {
		require.Contains(t, path, "/Nudge.app/Contents/MacOS/Nudge")
		return nil
	}
	runner.updater.opt.Targets["nudge"] = ti

	// trigger an update check
	updated, err := runner.UpdateAction()
	require.NoError(t, err)
	require.True(t, updated)

	// doesn't re-update after an update
	gotCfg, err = f.GetConfig()
	require.NoError(t, err)
	require.Equal(t, cfg, gotCfg)
	updated, err = runner.UpdateAction()
	require.NoError(t, err)
	require.False(t, updated)

	// runner hashes are updated
	b, ok := runner.localHashes["nudge"]
	require.True(t, ok)
	require.NotEmpty(t, b)

	// a config is created on the next run after install
	gotCfg, err = f.GetConfig()
	require.NoError(t, err)
	require.Equal(t, cfg, gotCfg)
	configBytes, err := os.ReadFile(configPath)
	require.NoError(t, err)
	var savedConfig fleet.NudgeConfig
	err = json.Unmarshal(configBytes, &savedConfig)
	require.NoError(t, err)
	require.Equal(t, cfg.NudgeConfig, &savedConfig)

	// config on disk changes if the config from the server changes
	cfg.NudgeConfig.OSVersionRequirements[0].RequiredMinimumOSVersion = "13.1.1"
	gotCfg, err = f.GetConfig()
	require.NoError(t, err)
	require.Equal(t, cfg, gotCfg)
	configBytes, err = os.ReadFile(configPath)
	require.NoError(t, err)
	savedConfig = fleet.NudgeConfig{}
	err = json.Unmarshal(configBytes, &savedConfig)
	require.NoError(t, err)
	require.Equal(t, cfg.NudgeConfig, &savedConfig)

	// config permissions are always validated and set to the right value
	err = os.Chmod(configPath, constant.DefaultFileMode)
	require.NoError(t, err)
	gotCfg, err = f.GetConfig()
	require.NoError(t, err)
	require.Equal(t, cfg, gotCfg)
	fileInfo, err := os.Stat(configPath)
	require.NoError(t, err)
	require.Equal(t, fileInfo.Mode(), nudgeConfigFileMode)

	configBytes, err = os.ReadFile(configPath)
	require.NoError(t, err)
	savedConfig = fleet.NudgeConfig{}
	err = json.Unmarshal(configBytes, &savedConfig)
	require.NoError(t, err)
	require.Equal(t, cfg.NudgeConfig, &savedConfig)

	// mock exec command to test handling of nudge launch errors
	wantCmd := filepath.Join(
		tmpDir,
		"bin",
		"nudge",
		NudgeMacOSTarget.Platform,
		NudgeMacOSTarget.Channel,
		NudgeMacOSTarget.ExtractedExecSubPath[0],
	)
	wantArgs := []string{fmt.Sprintf("file://%s", configPath)}
	runNudgeFnInvoked = false

	// nudge launches successfully
	time.Sleep(1 * time.Second)
	execCmd = mockExecCommand(t, "mock stdout", "", wantCmd, wantArgs...)
	gotCfg, err = f.GetConfig()
	require.NoError(t, err)
	require.Equal(t, "mock stdout", execOut)
	require.True(t, runNudgeFnInvoked)
	runNudgeFnInvoked = false
	execOut = ""

	// nudge fails to launch, stderr is captured and logged
	time.Sleep(1 * time.Second)
	execCmd = mockExecCommand(t, "", "mock stderr", wantCmd, wantArgs...)
	gotCfg, err = f.GetConfig()
	require.ErrorContains(t, err, "exit status 1: mock stderr")
	require.Equal(t, "", execOut)
	require.True(t, runNudgeFnInvoked)
	runNudgeFnInvoked = false

	// after launch error, nudge will not launch again
	time.Sleep(1 * time.Second)
	gotCfg, err = f.GetConfig()
	require.NoError(t, err)
	require.Equal(t, "", execOut)
	require.False(t, runNudgeFnInvoked)
	time.Sleep(1 * time.Second)
	gotCfg, err = f.GetConfig()
	require.NoError(t, err)
	require.Equal(t, "", execOut)
	require.False(t, runNudgeFnInvoked)
	time.Sleep(1 * time.Second)
	gotCfg, err = f.GetConfig()
	require.NoError(t, err)
	require.NoError(t, err)
	require.Equal(t, "", execOut)
	require.False(t, runNudgeFnInvoked)

	// nudge is removed from targets when the config is not present
	cfg.NudgeConfig = nil
	gotCfg, err = f.GetConfig()
	require.NoError(t, err)
	require.Equal(t, cfg, gotCfg)
	targets = runner.updater.opt.Targets
	require.Empty(t, targets)
	ti, ok = targets["nudge"]
	require.False(t, ok)
	require.Empty(t, ti)
}

// mockExecCommand returns a function that can be used to mock exec.Command using TestHelperProcess.
func mockExecCommand(t *testing.T, mockStdout string, mockStderr string, wantCommand string, wantArgs ...string) func(command string, args ...string) *exec.Cmd {
	return func(command string, args ...string) *exec.Cmd {
		cs := []string{"-test.run=TestHelperProcess", "--", command}
		cs = append(cs, args...)

		cmd := exec.Command(os.Args[0], cs...)
		cmd.Env = []string{
			"GO_WANT_HELPER_PROCESS=1",
			fmt.Sprintf("GO_WANT_HELPER_PROCESS_COMMAND=%s", wantCommand),
			fmt.Sprintf("GO_WANT_HELPER_PROCESS_ARGS=%s", strings.Join(wantArgs, " ")),
		}
		if mockStdout != "" {
			cmd.Env = append(cmd.Env, fmt.Sprintf("GO_WANT_HELPER_PROCESS_STDOUT=%s", mockStdout))
		}
		if mockStderr != "" {
			cmd.Env = append(cmd.Env, fmt.Sprintf("GO_WANT_HELPER_PROCESS_STDERR=%s", mockStderr))
		}

		return cmd
	}
}

// TestHelperProcess is a helper process used for tests that mock exec.Command
//
// Inspired by: https://npf.io/2015/06/testing-exec-command/
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	wantCmd := os.Getenv("GO_WANT_HELPER_PROCESS_COMMAND")
	if gotCmd := os.Args[3]; gotCmd != wantCmd {
		fmt.Fprint(os.Stderr, fmt.Sprintf("expected command %s but got %s", wantCmd, gotCmd))
		os.Exit(1)
		return
	}
	wantArgs := os.Getenv("GO_WANT_HELPER_PROCESS_ARGS")
	if gotArgs := os.Args[4]; gotArgs != string(wantArgs) {
		fmt.Fprint(os.Stderr, fmt.Sprintf("expected arg %s but got %s", string(wantArgs), gotArgs))
		os.Exit(1)
		return
	}
	fmt.Fprintf(os.Stdout, os.Getenv("GO_WANT_HELPER_PROCESS_STDOUT"))

	err := os.Getenv("GO_WANT_HELPER_PROCESS_STDERR")
	if err != "" {
		fmt.Fprintf(os.Stderr, err)
		os.Exit(1)
	}

	os.Exit(0)

	return
}
