package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/fleetdm/fleet/v4/orbit/pkg/packaging"
	"github.com/fleetdm/fleet/v4/orbit/pkg/update"
	"github.com/fleetdm/fleet/v4/pkg/nettest"
	"github.com/stretchr/testify/require"
	"pault.ag/go/debian/deb"
)

func TestPackage(t *testing.T) {
	nettest.Run(t)

	updateOpt := update.DefaultOptions
	updateOpt.RootDirectory = t.TempDir()
	updatesData, err := packaging.InitializeUpdates(updateOpt)
	require.NoError(t, err)

	// --type is required
	runAppCheckErr(t, []string{"package", "deb"}, "Required flag \"type\" not set")

	// if you provide -fleet-url & --enroll-secret are required together
	runAppCheckErr(t, []string{"package", "--type=deb", "--fleet-url=https://localhost:8080"}, "--enroll-secret and --fleet-url must be provided together")
	runAppCheckErr(t, []string{"package", "--type=deb", "--enroll-secret=foobar"}, "--enroll-secret and --fleet-url must be provided together")

	// --insecure and --fleet-certificate are mutually exclusive
	runAppCheckErr(t, []string{"package", "--type=deb", "--insecure", "--fleet-certificate=test123"}, "--insecure and --fleet-certificate may not be provided together")

	// Test invalid PEM file provided in --fleet-certificate.
	certDir := t.TempDir()
	fleetCertificate := filepath.Join(certDir, "fleet.pem")
	err = os.WriteFile(fleetCertificate, []byte("undefined"), os.FileMode(0o644))
	require.NoError(t, err)
	runAppCheckErr(t, []string{"package", "--type=deb", fmt.Sprintf("--fleet-certificate=%s", fleetCertificate)}, fmt.Sprintf("failed to read fleet server certificate %q: invalid PEM file", fleetCertificate))

	if runtime.GOOS != "linux" {
		runAppCheckErr(t, []string{"package", "--type=msi", "--native-tooling"}, "native tooling is only available in Linux")
	}

	t.Run("deb", func(t *testing.T) {
		shorterEnrollSecret := "aa"
		longerEnrollSecret := "aaaaaa"
		runAppForTest(t, []string{"package", "--type=deb", "--insecure", "--disable-open-folder", "--enroll-secret=" + longerEnrollSecret, "--fleet-url=https://localhost:8080"})
		info, err := os.Stat(fmt.Sprintf("fleet-osquery_%s_amd64.deb", updatesData.OrbitVersion))
		require.NoError(t, err)
		require.Greater(t, info.Size(), int64(0))

		fd, err := os.Open(fmt.Sprintf("fleet-osquery_%s_amd64.deb", updatesData.OrbitVersion))
		if err != nil {
			panic(err)
		}

		debFile, err := deb.Load(fd, fmt.Sprintf("fleet-osquery_%s_amd64.deb", updatesData.OrbitVersion))
		if err != nil {
			panic(err)
		}
		for {
			hdr, err := debFile.Data.Next()
			if err == io.EOF {
				break
			}
			require.NoError(t, err)

			if hdr.Name != "./etc/default/orbit" {
				continue
			}

			data, err := io.ReadAll(debFile.Data)
			require.NoError(t, err)

			s := string(data)
			require.Contains(t, s, fmt.Sprintf("ORBIT_ENROLL_SECRET=%s\n", longerEnrollSecret))
		}

		fd.Close()

		runAppForTest(t, []string{"package", "--type=deb", "--insecure", "--disable-open-folder", "--enroll-secret=" + shorterEnrollSecret, "--fleet-url=https://localhost:8080"})
		info, err = os.Stat(fmt.Sprintf("fleet-osquery_%s_amd64.deb", updatesData.OrbitVersion))
		require.NoError(t, err)
		require.Greater(t, info.Size(), int64(0))

		fd, err = os.Open(fmt.Sprintf("fleet-osquery_%s_amd64.deb", updatesData.OrbitVersion))
		if err != nil {
			panic(err)
		}
		defer fd.Close()

		debFile, err = deb.Load(fd, fmt.Sprintf("fleet-osquery_%s_amd64.deb", updatesData.OrbitVersion))
		if err != nil {
			panic(err)
		}
		for {
			hdr, err := debFile.Data.Next()
			if err == io.EOF {
				break
			}
			require.NoError(t, err)

			if hdr.Name != "./etc/default/orbit" {
				continue
			}

			data, err := io.ReadAll(debFile.Data)
			require.NoError(t, err)

			s := string(data)
			require.Contains(t, s, fmt.Sprintf("ORBIT_ENROLL_SECRET=%s\n", shorterEnrollSecret))
		}
	})

	t.Run("--use-sytem-configuration can't be used on installers that aren't pkg", func(t *testing.T) {
		for _, p := range []string{"deb", "msi", "rpm", ""} {
			runAppCheckErr(
				t,
				[]string{"package", fmt.Sprintf("--type=%s", p), "--use-system-configuration"},
				"--use-system-configuration is only available for pkg installers",
			)
		}
	})

	// fleet-osquery.msi
	// runAppForTest(t, []string{"package", "--type=msi", "--insecure"}) TODO: this is currently failing on Github runners due to permission issues
	// info, err = os.Stat("orbit-osquery_0.0.3.msi")
	// require.NoError(t, err)
	// require.Greater(t, info.Size(), int64(0))

	// runAppForTest(t, []string{"package", "--type=pkg", "--insecure"}) TODO: had a hard time getting xar installed on Ubuntu
}
