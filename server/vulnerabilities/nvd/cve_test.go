package nvd

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/fleetdm/fleet/v4/pkg/nettest"
	"github.com/fleetdm/fleet/v4/server/fleet"
	"github.com/fleetdm/fleet/v4/server/mock"
	"github.com/fleetdm/fleet/v4/server/vulnerabilities/nvd/tools/cvefeed"
	"github.com/fleetdm/fleet/v4/server/vulnerabilities/nvd/tools/wfn"
	"github.com/go-kit/log"
	kitlog "github.com/go-kit/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// firefox93WindowsVulnerabilities was manually generated by visiting:
// https://nvd.nist.gov/vuln/search/results?form_type=Advanced&results_type=overview&isCpeNameSearch=true&seach_type=all&query=cpe:2.3:a:mozilla:firefox:93.0:*:*:*:*:*:*:*
type cve struct {
	ID                string
	resolvedInVersion string
}

var firefox93WindowsVulnerabilities = []cve{
	{ID: "CVE-2021-43540", resolvedInVersion: "95.0"},
	{ID: "CVE-2021-38503", resolvedInVersion: "94.0"},
	{ID: "CVE-2021-38504", resolvedInVersion: "94.0"},
	{ID: "CVE-2021-38506", resolvedInVersion: "94.0"},
	{ID: "CVE-2021-38507", resolvedInVersion: "94.0"},
	{ID: "CVE-2021-38508", resolvedInVersion: "94.0"},
	{ID: "CVE-2021-38509", resolvedInVersion: "94.0"},
	{ID: "CVE-2021-43534", resolvedInVersion: "94.0"},
	{ID: "CVE-2021-43532", resolvedInVersion: "94.0"},
	{ID: "CVE-2021-43531", resolvedInVersion: "94.0"},
	{ID: "CVE-2021-43533", resolvedInVersion: "94.0"},
	{ID: "CVE-2021-43538", resolvedInVersion: "95.0"},
	{ID: "CVE-2021-43542", resolvedInVersion: "95.0"},
	{ID: "CVE-2021-43543", resolvedInVersion: "95.0"},
	{ID: "CVE-2021-30547"},
	{ID: "CVE-2021-43546", resolvedInVersion: "95.0"},
	{ID: "CVE-2021-43537", resolvedInVersion: "95.0"},
	{ID: "CVE-2021-43541", resolvedInVersion: "95.0"},
	{ID: "CVE-2021-43536", resolvedInVersion: "95.0"},
	{ID: "CVE-2021-43545", resolvedInVersion: "95.0"},
	{ID: "CVE-2021-43539", resolvedInVersion: "95.0"},
	{ID: "CVE-2022-34480", resolvedInVersion: "102.0"},
	{ID: "CVE-2022-26387", resolvedInVersion: "98.0"},
	{ID: "CVE-2022-22759", resolvedInVersion: "97.0"},
	{ID: "CVE-2022-28281", resolvedInVersion: "99.0"},
	{ID: "CVE-2022-45415", resolvedInVersion: "107.0"},
	{ID: "CVE-2022-42930", resolvedInVersion: "106.0"},
	{ID: "CVE-2022-0511", resolvedInVersion: "97.0"},
	{ID: "CVE-2022-22763", resolvedInVersion: "96.0"},
	{ID: "CVE-2022-22737", resolvedInVersion: "96.0"},
	{ID: "CVE-2022-22751", resolvedInVersion: "96.0"},
	{ID: "CVE-2022-38478", resolvedInVersion: "104.0"},
	{ID: "CVE-2022-22761", resolvedInVersion: "97.0"},
	{ID: "CVE-2022-34482", resolvedInVersion: "102.0"},
	{ID: "CVE-2022-26486", resolvedInVersion: "97.0.2"},
	{ID: "CVE-2022-22739", resolvedInVersion: "96.0"},
	{ID: "CVE-2022-22755", resolvedInVersion: "97.0"},
	{ID: "CVE-2022-22757", resolvedInVersion: "97.0"},
	{ID: "CVE-2022-1097", resolvedInVersion: "99.0"},
	{ID: "CVE-2022-22754", resolvedInVersion: "97.0"},
	{ID: "CVE-2022-22748", resolvedInVersion: "96.0"},
	{ID: "CVE-2022-22736", resolvedInVersion: "96.0"},
	{ID: "CVE-2022-22745", resolvedInVersion: "96.0"},
	{ID: "CVE-2022-26385", resolvedInVersion: "98.0"},
	{ID: "CVE-2022-26383", resolvedInVersion: "98.0"},
	{ID: "CVE-2022-3266", resolvedInVersion: "105.0"},
	{ID: "CVE-2022-34468", resolvedInVersion: "102.0"},
	{ID: "CVE-2022-34481", resolvedInVersion: "102.0"},
	{ID: "CVE-2022-28289", resolvedInVersion: "99.0"},
	{ID: "CVE-2022-22741", resolvedInVersion: "96.0"},
	{ID: "CVE-2022-28284", resolvedInVersion: "99.0"},
	{ID: "CVE-2022-34484", resolvedInVersion: "102.0"},
	{ID: "CVE-2022-22752", resolvedInVersion: "96.0"},
	{ID: "CVE-2022-26485", resolvedInVersion: "97.0.2"},
	{ID: "CVE-2022-28286", resolvedInVersion: "99.0"},
	{ID: "CVE-2022-28283", resolvedInVersion: "99.0"},
	{ID: "CVE-2022-28285", resolvedInVersion: "99.0"},
	{ID: "CVE-2022-0843", resolvedInVersion: "97.0"},
	{ID: "CVE-2022-29909", resolvedInVersion: "100.0"},
	{ID: "CVE-2022-22749", resolvedInVersion: "96.0"},
	{ID: "CVE-2022-26384", resolvedInVersion: "98.0"},
	{ID: "CVE-2022-28282", resolvedInVersion: "99.0"},
	{ID: "CVE-2022-28287", resolvedInVersion: "99.0"},
	{ID: "CVE-2022-40956", resolvedInVersion: "105.0"},
	{ID: "CVE-2022-22740", resolvedInVersion: "96.0"},
	{ID: "CVE-2022-22743", resolvedInVersion: "96.0"},
	{ID: "CVE-2022-22764", resolvedInVersion: "97.0"},
	{ID: "CVE-2022-22738", resolvedInVersion: "96.0"},
	{ID: "CVE-2022-1529", resolvedInVersion: "100.0.2"},
	{ID: "CVE-2022-22760", resolvedInVersion: "97.0"},
	{ID: "CVE-2022-29916", resolvedInVersion: "100.0"},
	{ID: "CVE-2022-29917", resolvedInVersion: "100.0"},
	{ID: "CVE-2022-22747", resolvedInVersion: "96.0"},
	{ID: "CVE-2022-26382", resolvedInVersion: "98.0"},
	{ID: "CVE-2022-22742", resolvedInVersion: "96.0"},
	{ID: "CVE-2022-28288", resolvedInVersion: "99.0"},
	{ID: "CVE-2022-22756", resolvedInVersion: "97.0"},
	{ID: "CVE-2022-26381", resolvedInVersion: "98.0"},
	{ID: "CVE-2022-1802", resolvedInVersion: "100.0.2"},
	{ID: "CVE-2022-34483", resolvedInVersion: "102.0"},
	{ID: "CVE-2022-29915", resolvedInVersion: "100.0"},
}

type threadSafeDSMock struct {
	mu sync.Mutex
	*mock.Store
}

func (d *threadSafeDSMock) ListSoftwareCPEs(ctx context.Context) ([]fleet.SoftwareCPE, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.Store.ListSoftwareCPEs(ctx)
}

func (d *threadSafeDSMock) InsertSoftwareVulnerability(ctx context.Context, vuln fleet.SoftwareVulnerability, src fleet.VulnerabilitySource) (bool, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.Store.InsertSoftwareVulnerability(ctx, vuln, src)
}

func TestTranslateCPEToCVE(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	// NVD_TEST_VULNDB_DIR can be used to speed up development (sync vulnerability data only once).
	tempDir := os.Getenv("NVD_TEST_VULNDB_DIR")
	if tempDir == "" {
		// download the CVEs once for all sub-tests, and then disable syncing
		tempDir = t.TempDir()
		err := nettest.RunWithNetRetry(t, func() error {
			return DownloadCVEFeed(tempDir, "", false, log.NewNopLogger())
		})
		require.NoError(t, err)
	} else {
		require.DirExists(t, tempDir)
		t.Logf("Using %s as database path", tempDir)
	}

	cveTests := map[string]struct {
		cpe          string
		excludedCVEs []string
		includedCVEs []cve
		// continuesToUpdate indicates if the product/software
		// continues to register new CVE vulnerabilities.
		continuesToUpdate bool
	}{
		"cpe:2.3:a:1password:1password:3.9.9:*:*:*:*:macos:*:*": {
			includedCVEs: []cve{
				{ID: "CVE-2012-6369"},
			},
			continuesToUpdate: false,
		},
		"cpe:2.3:a:1password:1password:3.9.9:*:*:*:*:*:*:*": {
			includedCVEs: []cve{
				{ID: "CVE-2012-6369"},
			},
			continuesToUpdate: false,
		},
		"cpe:2.3:a:pypa:pip:9.0.3:*:*:*:*:python:*:*": {
			includedCVEs: []cve{
				{ID: "CVE-2019-20916", resolvedInVersion: "19.2"},
				{ID: "CVE-2021-3572", resolvedInVersion: "21.1"},
				{ID: "CVE-2023-5752", resolvedInVersion: "23.3"},
			},
			continuesToUpdate: false,
		},
		"cpe:2.3:a:mozilla:firefox:93.0:*:*:*:*:windows:*:*": {
			includedCVEs:      firefox93WindowsVulnerabilities,
			continuesToUpdate: true,
		},
		"cpe:2.3:a:mozilla:firefox:93.0.100:*:*:*:*:windows:*:*": {
			includedCVEs:      firefox93WindowsVulnerabilities,
			continuesToUpdate: true,
		},
		"cpe:2.3:a:apple:icloud:1.0:*:*:*:*:macos:*:*": {
			excludedCVEs: []string{
				"CVE-2017-13797",
				"CVE-2017-2383",
				"CVE-2017-2366",
				"CVE-2016-4613",
				"CVE-2016-4692",
				"CVE-2016-4743",
				"CVE-2016-7578",
				"CVE-2016-7583",
				"CVE-2016-7586",
				"CVE-2016-7587",
				"CVE-2016-7589",
				"CVE-2016-7592",
				"CVE-2016-7598",
				"CVE-2016-7599",
				"CVE-2016-7610",
				"CVE-2016-7611",
				"CVE-2016-7614",
				"CVE-2016-7632",
				"CVE-2016-7635",
				"CVE-2016-7639",
				"CVE-2016-7640",
				"CVE-2016-7641",
				"CVE-2016-7642",
				"CVE-2016-7645",
				"CVE-2016-7646",
				"CVE-2016-7648",
				"CVE-2016-7649",
				"CVE-2016-7652",
				"CVE-2016-7654",
				"CVE-2016-7656",
				"CVE-2017-2383",
			},
			continuesToUpdate: true,
		},
		"cpe:2.3:a:clickstudios:passwordstate:9.5.8.4:*:*:*:*:chrome:*:*": {
			includedCVEs: []cve{
				{ID: "CVE-2022-4610"},
				{ID: "CVE-2022-4611"},
				{ID: "CVE-2022-4613"},
				{ID: "CVE-2022-4612"},
			},
			continuesToUpdate: true,
		},
		"cpe:2.3:a:avira:password_manager:2.18.4.38471:*:*:*:*:firefox:*:*": {
			includedCVEs: []cve{
				{ID: "CVE-2022-28795"},
			},
			continuesToUpdate: true,
		},
		"cpe:2.3:a:zoom:zoom:5.0.4301.0407:*:*:*:*:chrome:*:*": {
			excludedCVEs:      []string{"CVE-2021-28133"}, // CVE-2021-28133 is a vulnerability in the Zoom application, not the extension.
			continuesToUpdate: true,
		},
		"cpe:2.3:a:bitwarden:bitwarden:1.55.0:*:*:*:*:firefox:*:*": {
			excludedCVEs:      []string{"CVE-2023-38840"}, // CVE-2023-38840 is a vulnerability in the Bitwarden application, not the extension.
			continuesToUpdate: true,
		},
		// No Version Start
		"cpe:2.3:a:python:setuptools:64:*:*:*:*:*:*:*": {
			includedCVEs: []cve{
				{ID: "CVE-2022-40897", resolvedInVersion: "65.5.1"},
			},
			continuesToUpdate: true,
		},
		"cpe:2.3:a:mozilla:firefox:93.0.100:*:*:*:*:*:*:*": {
			includedCVEs: []cve{
				// CVE matches multiple products
				{ID: "CVE-2022-40956", resolvedInVersion: "105.0"},
			},
			continuesToUpdate: true,
		},
		// Multiple product matches with different version ranges
		"cpe:2.3:o:apple:macos:14.1:*:*:*:*:*:*:*": {
			includedCVEs: []cve{
				{ID: "CVE-2023-42919", resolvedInVersion: "14.2"},
			},
			continuesToUpdate: true,
		},
		"cpe:2.3:a:microsoft:windows_subsystem_for_linux:0.63.10:*:*:*:*:visual_studio_code:*:*": {
			includedCVEs: []cve{
				{ID: "CVE-2021-43907", resolvedInVersion: "0.63.11"},
			},
			continuesToUpdate: false,
		},
		"cpe:2.3:a:github:pull_requests_and_issues:0.66.1:*:*:*:*:visual_studio_code:*:*": {
			includedCVEs: []cve{
				{ID: "CVE-2023-36867", resolvedInVersion: "0.66.2"},
			},
			continuesToUpdate: false,
		},
		"cpe:2.3:a:microsoft:python_extension:2020.9.1:*:*:*:*:visual_studio_code:*:*": {
			includedCVEs: []cve{
				{ID: "CVE-2020-17163", resolvedInVersion: "2020.9.2"},
			},
			continuesToUpdate: false,
		},
		"cpe:2.3:a:microsoft:jupyter:2023.10.10:*:*:*:*:visual_studio_code:*:*": {
			includedCVEs: []cve{
				{ID: "CVE-2023-36018", resolvedInVersion: "2023.10.1100000000"},
			},
			continuesToUpdate: false,
		},
		"cpe:2.3:a:microsoft:jupyter:2024.2.0:*:*:*:*:visual_studio_code:*:*": {
			includedCVEs:      []cve{},
			continuesToUpdate: false,
		},
		"cpe:2.3:a:microsoft:visual_studio_code_eslint_extension:2.0.0:*:*:*:*:visual_studio_code:*:*": {
			includedCVEs: []cve{
				{ID: "CVE-2020-1481", resolvedInVersion: "2.1.7"},
			},
			continuesToUpdate: false,
		},
		"cpe:2.3:a:microsoft:python_extension:2020.4.0:*:*:*:*:visual_studio_code:*:*": {
			includedCVEs: []cve{
				{ID: "CVE-2020-1171", resolvedInVersion: "2020.5.0"},
				{ID: "CVE-2020-1192", resolvedInVersion: "2020.5.0"},
				{ID: "CVE-2020-17163", resolvedInVersion: "2020.9.2"},
			},
			continuesToUpdate: false,
		},
		"cpe:2.3:a:adobe:animate:*:*:*:*:*:macos:*:*": {
			includedCVEs: []cve{
				{ID: "CVE-2023-44325"},
			},
			continuesToUpdate: true,
		},
		"cpe:2.3:a:apple:safari:17.0:*:*:*:*:macos:*:*": {
			includedCVEs: []cve{
				{ID: "CVE-2023-42852", resolvedInVersion: "17.1"},
				{ID: "CVE-2023-42950", resolvedInVersion: "17.2"},
				{ID: "CVE-2024-23273", resolvedInVersion: "17.4"},
			},
			excludedCVEs:      []string{"CVE-2023-28205"},
			continuesToUpdate: true,
		},
		"cpe:2.3:a:apple:safari:16.4.0:*:*:*:*:macos:*:*": {
			includedCVEs: []cve{
				{ID: "CVE-2023-28205", resolvedInVersion: "16.4.1"},
			},
			continuesToUpdate: true,
		},
	}

	cveOSTests := []struct {
		platform     string
		version      string
		osID         uint
		includedCVEs []string
	}{
		{
			platform: "darwin",
			version:  "14.1.2",
			osID:     1,
			includedCVEs: []string{
				"CVE-2023-45866",
				"CVE-2023-42886",
				"CVE-2023-42891",
				"CVE-2023-42906",
				"CVE-2023-42910",
				"CVE-2023-42924",
				"CVE-2023-42883",
				"CVE-2023-42894",
				"CVE-2023-42926",
				"CVE-2023-42932",
				"CVE-2023-42907",
				"CVE-2023-42922",
				"CVE-2023-42904",
				"CVE-2023-42901",
				"CVE-2023-42898",
				"CVE-2023-42903",
				"CVE-2023-42902",
				"CVE-2023-42909",
				"CVE-2023-42914",
				"CVE-2023-42874",
				"CVE-2023-42882",
				"CVE-2023-42912",
				"CVE-2023-42911",
				"CVE-2023-42890",
				"CVE-2023-42905",
				"CVE-2023-42919",
				"CVE-2023-42900",
				"CVE-2023-42899",
				"CVE-2023-42908",
				"CVE-2023-42884",
			},
		},
		{
			platform: "darwin",
			version:  "13.6.2",
			osID:     2,
			// This is a subset of vulnerabilities for macOS 13.6.2
			includedCVEs: []string{
				"CVE-2023-32361",
				"CVE-2023-35990",
				"CVE-2023-40541",
				"CVE-2023-40400",
				"CVE-2023-41980",
				"CVE-2023-38615",
				"CVE-2023-39233",
				"CVE-2023-40402",
				"CVE-2023-40450",
				"CVE-2023-42891",
				"CVE-2023-41079",
				"CVE-2023-42932",
				"CVE-2023-38586",
				"CVE-2023-41067",
				"CVE-2023-40407",
				"CVE-2023-42924",
				"CVE-2023-40395",
				"CVE-2023-38596",
				"CVE-2023-32396",
				"CVE-2023-29497",
			},
		},
	}

	t.Run("find_vulns_on_cpes", func(t *testing.T) {
		t.Parallel()

		ds := new(mock.Store)

		softwareIDToCPEs := make(map[uint]string)
		ds.ListSoftwareCPEsFunc = func(ctx context.Context) ([]fleet.SoftwareCPE, error) {
			var softwareCPEs []fleet.SoftwareCPE
			i := uint(0)
			for cpe := range cveTests {
				softwareCPEs = append(softwareCPEs, fleet.SoftwareCPE{CPE: cpe, SoftwareID: i})
				softwareIDToCPEs[i] = cpe
				i++
			}
			return softwareCPEs, nil
		}

		var osIDs []uint
		ds.ListOperatingSystemsForPlatformFunc = func(ctx context.Context, p string) ([]fleet.OperatingSystem, error) {
			var oss []fleet.OperatingSystem
			for _, os := range cveOSTests {
				oss = append(oss, fleet.OperatingSystem{
					ID:       os.osID,
					Platform: os.platform,
					Version:  os.version,
				})
				osIDs = append(osIDs, os.osID)
			}
			return oss, nil
		}

		cveLock := &sync.Mutex{}
		cvesFound := make(map[string][]cve)
		ds.InsertSoftwareVulnerabilityFunc = func(ctx context.Context, vuln fleet.SoftwareVulnerability, src fleet.VulnerabilitySource) (bool, error) {
			cveLock.Lock()
			defer cveLock.Unlock()

			cpe, ok := softwareIDToCPEs[vuln.SoftwareID]
			if !ok {
				return false, fmt.Errorf("software id -> cpe not found: %d", vuln.SoftwareID)
			}
			cve := cve{
				ID:                vuln.CVE,
				resolvedInVersion: *vuln.ResolvedInVersion,
			}
			cvesFound[cpe] = append(cvesFound[cpe], cve)
			return false, nil
		}

		osCVELock := &sync.Mutex{}
		osCVEsFound := make(map[uint][]string)
		ds.InsertOSVulnerabilityFunc = func(ctx context.Context, vuln fleet.OSVulnerability, src fleet.VulnerabilitySource) (bool, error) {
			osCVELock.Lock()
			defer osCVELock.Unlock()

			osCVEsFound[vuln.OSID] = append(osCVEsFound[vuln.OSID], vuln.CVE)

			return false, nil
		}

		ds.DeleteOutOfDateVulnerabilitiesFunc = func(ctx context.Context, source fleet.VulnerabilitySource, duration time.Duration) error {
			return nil
		}
		ds.DeleteOutOfDateOSVulnerabilitiesFunc = func(ctx context.Context, source fleet.VulnerabilitySource, duration time.Duration) error {
			return nil
		}

		_, err := TranslateCPEToCVE(ctx, ds, tempDir, kitlog.NewNopLogger(), false, 1*time.Hour)
		require.NoError(t, err)

		require.True(t, ds.DeleteOutOfDateVulnerabilitiesFuncInvoked)
		require.True(t, ds.DeleteOutOfDateOSVulnerabilitiesFuncInvoked)

		for cpe, tc := range cveTests {
			if tc.continuesToUpdate {
				// Given that new vulnerabilities can be found on these
				// packages/products, we check that at least the
				// known ones are found.
				for _, cve := range tc.includedCVEs {
					require.Contains(t, cvesFound[cpe], cve, fmt.Sprintf("%s does not contain CVE %#v", cpe, cve))
				}
			} else {
				// Check for exact match of CVEs found.
				require.ElementsMatch(t, cvesFound[cpe], tc.includedCVEs, cpe)
			}

			for _, cve := range tc.excludedCVEs {
				for _, cveFound := range cvesFound[cpe] {
					require.NotEqual(t, cve, cveFound.ID, fmt.Sprintf("%s should not contain %s", cpe, cve))
				}
			}
		}

		for _, tc := range cveOSTests {
			for _, cve := range tc.includedCVEs {
				require.Contains(t, osCVEsFound[tc.osID], cve)
			}
		}
	})

	t.Run("recent_vulns", func(t *testing.T) {
		t.Parallel()

		ds := new(mock.Store)
		safeDS := &threadSafeDSMock{Store: ds}

		softwareCPEs := []fleet.SoftwareCPE{
			{CPE: "cpe:2.3:a:google:chrome:*:*:*:*:*:*:*:*", ID: 1, SoftwareID: 1},
			{CPE: "cpe:2.3:a:mozilla:firefox:*:*:*:*:*:*:*:*", ID: 2, SoftwareID: 2},
			{CPE: "cpe:2.3:a:haxx:curl:*:*:*:*:*:*:*:*", ID: 3, SoftwareID: 3},
		}
		ds.DeleteOutOfDateVulnerabilitiesFunc = func(ctx context.Context, source fleet.VulnerabilitySource, duration time.Duration) error {
			return nil
		}
		ds.ListSoftwareCPEsFunc = func(ctx context.Context) ([]fleet.SoftwareCPE, error) {
			return softwareCPEs, nil
		}
		ds.InsertSoftwareVulnerabilityFunc = func(ctx context.Context, vuln fleet.SoftwareVulnerability, src fleet.VulnerabilitySource) (bool, error) {
			return true, nil
		}
		ds.ListOperatingSystemsForPlatformFunc = func(ctx context.Context, p string) ([]fleet.OperatingSystem, error) {
			return nil, nil
		}
		ds.DeleteOutOfDateOSVulnerabilitiesFunc = func(ctx context.Context, source fleet.VulnerabilitySource, duration time.Duration) error {
			return nil
		}

		recent, err := TranslateCPEToCVE(ctx, safeDS, tempDir, kitlog.NewNopLogger(), true, 1*time.Hour)
		require.NoError(t, err)

		byCPE := make(map[uint]int)
		for _, cpe := range recent {
			byCPE[cpe.Affected()]++
		}

		// even if it's somewhat far in the past, I've seen the exact numbers
		// change a bit between runs with different downloads, so allow for a bit
		// of wiggle room.
		assert.Greater(t, byCPE[softwareCPEs[0].SoftwareID], 150, "google chrome CVEs")
		assert.Greater(t, byCPE[softwareCPEs[1].SoftwareID], 280, "mozilla firefox CVEs")
		assert.Greater(t, byCPE[softwareCPEs[2].SoftwareID], 10, "curl CVEs")

		// call it again but now return false from this call, simulating CVE-CPE pairs
		// that already existed in the DB.
		ds.InsertSoftwareVulnerabilityFunc = func(ctx context.Context, vuln fleet.SoftwareVulnerability, src fleet.VulnerabilitySource) (bool, error) {
			return false, nil
		}
		recent, err = TranslateCPEToCVE(ctx, safeDS, tempDir, kitlog.NewNopLogger(), true, 1*time.Hour)
		require.NoError(t, err)

		// no recent vulnerability should be reported
		assert.Len(t, recent, 0)
	})
}

func TestSyncsCVEFromURL(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.RequestURI, ".meta") {
			fmt.Fprint(w, "lastModifiedDate:2021-08-04T11:10:30-04:00\r\n")
			fmt.Fprint(w, "size:20967174\r\n")
			fmt.Fprint(w, "zipSize:1453429\r\n")
			fmt.Fprint(w, "gzSize:1453293\r\n")
			fmt.Fprint(w, "sha256:10D7338A1E2D8DB344C381793110B67FCA7D729ADA21624EF089EBA78CCE7B53\r\n")
		}
	}))
	defer ts.Close()

	tempDir := t.TempDir()
	cveFeedPrefixURL := ts.URL + "/feeds/json/cve/1.1/"
	err := DownloadCVEFeed(tempDir, cveFeedPrefixURL, false, log.NewNopLogger())
	require.Error(t, err)
	require.Contains(t,
		err.Error(),
		fmt.Sprintf("1 synchronisation error:\n\tunexpected size for \"%s/feeds/json/cve/1.1/nvdcve-1.1-2002.json.gz\" (200 OK): want 1453293, have 0", ts.URL),
	)
}

// This test is using real data from the 2022 NVD feed
func TestGetMatchingVersionEndExcluding(t *testing.T) {
	ctx := context.Background()
	testDict := loadDict(t, "../testdata/nvdcve-1.1-2022.json.gz")

	tests := []struct {
		name    string
		cve     string
		meta    *wfn.Attributes
		want    string
		wantErr bool
	}{
		{
			name: "happy path with version with no Version Start",
			cve:  "CVE-2022-40897",
			meta: &wfn.Attributes{
				Vendor:  "python",
				Product: "setuptools",
				Version: "64",
			},
			want:    "65.5.1",
			wantErr: false,
		},
		{
			name: "CVE matches multiple products",
			cve:  "CVE-2022-40956",
			meta: &wfn.Attributes{
				Vendor:  "mozilla",
				Product: "firefox",
				Version: "93.0.100",
			},
			want:    "105.0",
			wantErr: false,
		},
		{
			name: "Nodes has nested Children",
			cve:  "CVE-2022-40961",
			meta: &wfn.Attributes{
				Vendor:  "mozilla",
				Product: "firefox",
				Version: "93.0.100",
			},
			want:    "105.0",
			wantErr: false,
		},
		{
			name: "Multiple product matches with different version ranges",
			cve:  "CVE-2022-26697",
			meta: &wfn.Attributes{
				Vendor:  "apple",
				Product: "macos",
				Version: "12.0",
			},
			want:    "12.4",
			wantErr: false,
		},
		{
			name: "No version end excluding",
			cve:  "CVE-2022-26834",
			meta: &wfn.Attributes{
				Vendor:  "cybozu",
				Product: "remote_service_manager",
				Version: "3.1.2",
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "Version exceeds version end excluding",
			cve:  "CVE-2022-4610",
			meta: &wfn.Attributes{
				Vendor:   "clickstudios",
				Product:  "passwordstate",
				Version:  "9.5.8.4",
				TargetSW: "chrome",
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "Can compare 4th version part",
			cve:  "CVE-2022-45889",
			meta: &wfn.Attributes{
				Vendor:  "planetestream",
				Product: "planet_estream",
				Version: "6.72.10.06",
			},
			want:    "6.72.10.07",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getMatchingVersionEndExcluding(ctx, tt.cve, tt.meta, testDict, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("getMatchingVersionEndExcluding() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPreprocessVersion(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"2.3.0.2", "2.3.0-2"},
		{"2.3.0+2", "2.3.0+2"},
		{"v5.3.0.2", "v5.3.0-2"},
		{"5.3.0-2", "5.3.0-2"},
		{"2.3.0.2.5", "2.3.0-2.5"},
		{"2.3.0", "2.3.0"},
		{"2.3", "2.3"},
		{"v2.3.0", "v2.3.0"},
		{"notAVersion", "notAVersion"},
		{"2.0.0+svn315-7fakesync1ubuntu0.22.04.1", "2.0.0+svn315-7fakesync1ubuntu0.22.04.1"},
		{"1.21.1ubuntu2", "1.21.1-ubuntu2"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			output := preprocessVersion(tc.input)
			if output != tc.expected {
				t.Fatalf("input: %s, expected: %s, got: %s", tc.input, tc.expected, output)
			}
		})
	}
}

func TestGetMacOSCPEs(t *testing.T) {
	ctx := context.Background()
	ds := new(mock.Store)
	os := fleet.OperatingSystem{
		ID:            1,
		Name:          "macOS",
		Version:       "11.6.2",
		Arch:          "x86_64",
		KernelVersion: "20.6.0",
		Platform:      "darwin",
	}

	ds.ListOperatingSystemsForPlatformFunc = func(ctx context.Context, p string) ([]fleet.OperatingSystem, error) {
		return []fleet.OperatingSystem{os}, nil
	}

	CVEs, err := GetMacOSCPEs(ctx, ds)
	require.NoError(t, err)
	require.Len(t, CVEs, 2)

	expected := map[osCPEWithNVDMeta]struct{}{
		{
			OperatingSystem: os,
			meta: &wfn.Attributes{
				Part:      "o",
				Vendor:    "apple",
				Product:   "mac_os_x",
				Version:   CVEs[0].Version,
				Update:    wfn.Any,
				Edition:   wfn.Any,
				SWEdition: wfn.Any,
				TargetSW:  wfn.Any,
				TargetHW:  wfn.Any,
				Other:     wfn.Any,
				Language:  wfn.Any,
			},
		}: {},
		{
			OperatingSystem: os,
			meta: &wfn.Attributes{
				Part:      "o",
				Vendor:    "apple",
				Product:   "macos",
				Version:   CVEs[0].Version,
				Update:    wfn.Any,
				Edition:   wfn.Any,
				SWEdition: wfn.Any,
				TargetSW:  wfn.Any,
				TargetHW:  wfn.Any,
				Other:     wfn.Any,
				Language:  wfn.Any,
			},
		}: {},
	}

	for _, cve := range CVEs {
		require.Contains(t, expected, cve)
	}
}

// loadDict loads a cvefeed.Dictionary from a JSON NVD feed file.
func loadDict(t *testing.T, path string) cvefeed.Dictionary {
	dict, err := cvefeed.LoadJSONDictionary(path)
	if err != nil {
		t.Fatal(err)
	}
	return dict
}

func TestExpandCPEAliases(t *testing.T) {
	firefox := &wfn.Attributes{
		Vendor:  "mozilla",
		Product: "firefox",
		Version: "93.0.100",
	}
	chromePlugin := &wfn.Attributes{
		Vendor:   "google",
		Product:  "plugin foobar",
		Version:  "93.0.100",
		TargetSW: "chrome",
	}

	vsCodeExtension := &wfn.Attributes{
		Vendor:   "Microsoft",
		Product:  "foo.extension",
		Version:  "2024.2.1",
		TargetSW: "visual_studio_code",
	}
	vsCodeExtensionAlias := *vsCodeExtension
	vsCodeExtensionAlias.TargetSW = "visual_studio"

	pythonCodeExtension := &wfn.Attributes{
		Vendor:   "microsoft",
		Product:  "python_extension",
		Version:  "2024.2.1",
		TargetSW: "visual_studio_code",
	}
	pythonCodeExtensionAlias1 := *pythonCodeExtension
	pythonCodeExtensionAlias1.TargetSW = "visual_studio"
	pythonCodeExtensionAlias2 := *pythonCodeExtension
	pythonCodeExtensionAlias2.Product = "visual_studio_code"
	pythonCodeExtensionAlias2.TargetSW = "python"

	for _, tc := range []struct {
		name            string
		cpeItem         *wfn.Attributes
		expectedAliases []*wfn.Attributes
	}{
		{
			name:            "no expansion without target_sw",
			cpeItem:         firefox,
			expectedAliases: []*wfn.Attributes{firefox},
		},
		{
			name:            "no expansion with target_sw",
			cpeItem:         chromePlugin,
			expectedAliases: []*wfn.Attributes{chromePlugin},
		},
		{
			name:            "visual studio code extension",
			cpeItem:         vsCodeExtension,
			expectedAliases: []*wfn.Attributes{vsCodeExtension, &vsCodeExtensionAlias},
		},
		{
			name:            "python visual studio code extension",
			cpeItem:         pythonCodeExtension,
			expectedAliases: []*wfn.Attributes{pythonCodeExtension, &pythonCodeExtensionAlias1, &pythonCodeExtensionAlias2},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			aliases := expandCPEAliases(tc.cpeItem)
			require.Equal(t, tc.expectedAliases, aliases)
		})
	}
}
