package oval

import (
	"compress/bzip2"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/fleetdm/fleet/v4/server"
	"github.com/fleetdm/fleet/v4/server/datastore/mysql"
	"github.com/fleetdm/fleet/v4/server/fleet"
	"github.com/stretchr/testify/require"
)

type softwareFixture struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Release string `json:"release"`
	Arch    string `json:"arch"`
}

func extract(src, dst string, t require.TestingT) {
	srcF, err := os.Open(src)
	require.NoError(t, err)
	defer srcF.Close()

	dstF, err := os.Create(dst)
	require.NoError(t, err)
	defer dstF.Close()

	r := bzip2.NewReader(srcF)
	// ignoring "G110: Potential DoS vulnerability via decompression bomb", as this is test code.
	_, err = io.Copy(dstF, r) //nolint:gosec
	require.NoError(t, err)
}

func loadSoftware(
	ds *mysql.Datastore,
	p Platform,
	s fleet.OSVersion,
	vulnPath string,
	t require.TestingT,
) *fleet.Host {
	osqueryHostID, err := server.GenerateRandomText(10)
	require.NoError(t, err)

	ctx := context.Background()

	h, err := ds.NewHost(context.Background(), &fleet.Host{
		Hostname:        string(p),
		NodeKey:         string(p),
		UUID:            string(p),
		DetailUpdatedAt: time.Now(),
		LabelUpdatedAt:  time.Now(),
		PolicyUpdatedAt: time.Now(),
		SeenTime:        time.Now(),
		OsqueryHostID:   osqueryHostID,
		Platform:        s.Platform,
		OSVersion:       s.Name,
	})
	require.NoError(t, err)

	var fixtures []softwareFixture
	contents, err := ioutil.ReadFile(filepath.Join(vulnPath, fmt.Sprintf("%s-software.json", p)))
	require.NoError(t, err)

	err = json.Unmarshal(contents, &fixtures)
	require.NoError(t, err)

	var software []fleet.Software
	for _, fi := range fixtures {
		software = append(software, fleet.Software{
			Name:    fi.Name,
			Version: fi.Version,
			Release: fi.Release,
			Arch:    fi.Arch,
		})
	}
	err = ds.UpdateHostSoftware(ctx, h.ID, software)
	require.NoError(t, err)

	err = ds.LoadHostSoftware(ctx, h, false)
	require.NoError(t, err)

	for _, s := range h.Software {
		err = ds.AddCPEForSoftware(ctx, s, fmt.Sprintf("%s-%s", s.Name, s.Version))
		require.NoError(t, err)
	}

	return h
}

func extractFixtures(
	p Platform,
	ovalFixtureDir string,
	softwareFixtureDir string,
	vulnPath string,
	t require.TestingT,
) {
	ovalFixPath := filepath.Join("..", "testdata", ovalFixtureDir)
	srcDefPath := filepath.Join(ovalFixPath, fmt.Sprintf("%s-oval_def.json.bz2", p))
	dstDefPath := filepath.Join(vulnPath, p.ToFilename(time.Now(), "json"))
	extract(srcDefPath, dstDefPath, t)

	softwareFixPath := filepath.Join("..", "testdata", softwareFixtureDir)
	srcSoftPath := filepath.Join(softwareFixPath, fmt.Sprintf("%s-software.json.bz2", p))
	dstSoftPath := filepath.Join(vulnPath, fmt.Sprintf("%s-software.json", p))
	extract(srcSoftPath, dstSoftPath, t)

	srcCvesPath := filepath.Join(softwareFixPath, fmt.Sprintf("%s-software_cves.csv.bz2", p))
	dstCvesPath := filepath.Join(vulnPath, fmt.Sprintf("%s-software_cves.csv", p))
	extract(srcCvesPath, dstCvesPath, t)
}

func withTestFixutre(
	version fleet.OSVersion,
	ovalFixtureDir string,
	softwareFixtureDir string,
	vulnPath string,
	ds *mysql.Datastore,
	afterLoad func(h *fleet.Host),
	t require.TestingT,
) {
	ctx := context.Background()
	p := NewPlatform(version.Platform, version.Name)

	extractFixtures(p, ovalFixtureDir, softwareFixtureDir, vulnPath, t)

	h := loadSoftware(ds, p, version, vulnPath, t)
	err := ds.UpdateOSVersions(ctx)
	require.NoError(t, err)
	afterLoad(h)
	err = ds.DeleteHost(ctx, h.ID)
	require.NoError(t, err)
}

func assertVulns(
	ds *mysql.Datastore,
	vulnPath string,
	h *fleet.Host,
	p Platform,
	t require.TestingT,
) {
	ctx := context.Background()

	fPath := filepath.Join(vulnPath, fmt.Sprintf("%s-software_cves.csv", p))
	f, err := os.Open(fPath)
	require.NoError(t, err)
	defer f.Close()

	r := csv.NewReader(f)
	var expected []string
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}

		if len(row) < 1 {
			continue
		}

		if len(row) > 1 && row[1] == "#ignore:" || strings.Index(row[0], "ignore") != -1 {
			continue
		}

		if !strings.HasPrefix(strings.ToLower(row[0]), "cve") {
			continue
		}

		expected = append(expected, row[0])
	}
	require.NotEmpty(t, expected)

	storedVulns, err := ds.ListSoftwareVulnerabilities(ctx, []uint{h.ID})
	require.NoError(t, err)

	uniq := make(map[string]bool)
	for _, v := range storedVulns[h.ID] {
		uniq[v.CVE] = true
	}
	actual := make([]string, 0, len(uniq))
	for k := range uniq {
		actual = append(actual, k)
	}

	require.ElementsMatch(t, actual, expected)
}

func BenchmarkTestOvalAnalyzer(b *testing.B) {
	b.Run("Ubuntu", func(b *testing.B) {
		ds := mysql.CreateMySQLDS(b)
		defer mysql.TruncateTables(b, ds)

		vulnPath := b.TempDir()

		systems := []fleet.OSVersion{
			{Platform: "ubuntu", Name: "Ubuntu 16.4.0"},
			{Platform: "ubuntu", Name: "Ubuntu 18.4.0"},
			{Platform: "ubuntu", Name: "Ubuntu 20.4.0"},
			{Platform: "ubuntu", Name: "Ubuntu 21.4.0"},
			{Platform: "ubuntu", Name: "Ubuntu 21.10.0"},
			{Platform: "ubuntu", Name: "Ubuntu 22.4.0"},
		}

		ovalFixtureDir := "ubuntu"
		softwareFixtureDir := filepath.Join("ubuntu", "software")

		for _, v := range systems {
			b.Run(fmt.Sprintf("for %s %s", v.Platform, v.Name), func(b *testing.B) {
				withTestFixutre(v, ovalFixtureDir, softwareFixtureDir, vulnPath, ds, func(h *fleet.Host) {
					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						_, err := Analyze(context.Background(), ds, v, vulnPath, true)
						require.NoError(b, err)
					}
				}, b)
			})
		}
	})

	b.Run("RHEL", func(b *testing.B) {
		ds := mysql.CreateMySQLDS(b)
		defer mysql.TruncateTables(b, ds)

		vulnPath := b.TempDir()

		systems := []struct {
			softwareFixtureDir string
			ovalFixtureDir     string
			version            fleet.OSVersion
		}{
			{
				ovalFixtureDir:     "rhel",
				softwareFixtureDir: filepath.Join("rhel", "software", "0709"),
				version:            fleet.OSVersion{Platform: "rhel", Name: "Red Hat Enterprise Linux Server 7.9.0"},
			},
			{
				ovalFixtureDir:     "rhel",
				softwareFixtureDir: filepath.Join("rhel", "software", "0802"),
				version:            fleet.OSVersion{Platform: "rhel", Name: "Red Hat Enterprise Linux Server 8.2.0"},
			},
			{
				ovalFixtureDir:     "rhel",
				softwareFixtureDir: filepath.Join("rhel", "software", "0804"),
				version:            fleet.OSVersion{Platform: "rhel", Name: "Red Hat Enterprise Linux 8.4.0"},
			},
			{
				ovalFixtureDir:     "rhel",
				softwareFixtureDir: filepath.Join("rhel", "software", "0806"),
				version:            fleet.OSVersion{Platform: "rhel", Name: "Red Hat Enterprise Linux 8.6.0"},
			},
			{
				ovalFixtureDir:     "rhel",
				softwareFixtureDir: filepath.Join("rhel", "software", "0900"),
				version:            fleet.OSVersion{Platform: "rhel", Name: "Red Hat Enterprise Linux 9.0.0"},
			},
		}

		for _, v := range systems {
			b.Run(fmt.Sprintf("for %s %s", v.version.Platform, v.version.Name), func(b *testing.B) {
				withTestFixutre(v.version, v.ovalFixtureDir, v.softwareFixtureDir, vulnPath, ds, func(h *fleet.Host) {
					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						_, err := Analyze(context.Background(), ds, v.version, vulnPath, true)
						require.NoError(b, err)
					}
				}, b)
			})
		}
	})
}

func TestOvalAnalyzer(t *testing.T) {
	t.Run("analyzing RHEL software", func(t *testing.T) {
		ds := mysql.CreateMySQLDS(t)
		defer mysql.TruncateTables(t, ds)

		vulnPath := t.TempDir()

		ctx := context.Background()

		systems := []struct {
			softwareFixtureDir string
			ovalFixtureDir     string
			version            fleet.OSVersion
		}{
			{
				ovalFixtureDir:     "rhel",
				softwareFixtureDir: filepath.Join("rhel", "software", "0709"),
				version:            fleet.OSVersion{Platform: "rhel", Name: "Red Hat Enterprise Linux Server 7.9.0"},
			},
			{
				ovalFixtureDir:     "rhel",
				softwareFixtureDir: filepath.Join("rhel", "software", "0802"),
				version:            fleet.OSVersion{Platform: "rhel", Name: "Red Hat Enterprise Linux Server 8.2.0"},
			},
			{
				ovalFixtureDir:     "rhel",
				softwareFixtureDir: filepath.Join("rhel", "software", "0804"),
				version:            fleet.OSVersion{Platform: "rhel", Name: "Red Hat Enterprise Linux 8.4.0"},
			},
			{
				ovalFixtureDir:     "rhel",
				softwareFixtureDir: filepath.Join("rhel", "software", "0806"),
				version:            fleet.OSVersion{Platform: "rhel", Name: "Red Hat Enterprise Linux 8.6.0"},
			},
			{
				ovalFixtureDir:     "rhel",
				softwareFixtureDir: filepath.Join("rhel", "software", "0900"),
				version:            fleet.OSVersion{Platform: "rhel", Name: "Red Hat Enterprise Linux 9.0.0"},
			},
		}

		for _, s := range systems {
			withTestFixutre(s.version, s.ovalFixtureDir, s.softwareFixtureDir, vulnPath, ds, func(h *fleet.Host) {
				_, err := Analyze(ctx, ds, s.version, vulnPath, true)
				require.NoError(t, err)
				p := NewPlatform(s.version.Platform, s.version.Name)
				assertVulns(ds, vulnPath, h, p, t)
			}, t)
		}
	})

	// For generating the vulnerability lists I used VMs and ran oscap (since it seems like oscap
	// does not work with Docker) and extracted all installed software vulnerabilities, then I had
	// the VMs join my local dev env, and extracted the installed software from the database.
	t.Run("analyzing Ubuntu software", func(t *testing.T) {
		ds := mysql.CreateMySQLDS(t)
		defer mysql.TruncateTables(t, ds)

		vulnPath := t.TempDir()

		ctx := context.Background()

		systems := []fleet.OSVersion{
			{Platform: "ubuntu", Name: "Ubuntu 16.4.0"},
			{Platform: "ubuntu", Name: "Ubuntu 18.4.0"},
			{Platform: "ubuntu", Name: "Ubuntu 20.4.0"},
			{Platform: "ubuntu", Name: "Ubuntu 21.4.0"},
			{Platform: "ubuntu", Name: "Ubuntu 21.10.0"},
			{Platform: "ubuntu", Name: "Ubuntu 22.4.0"},
		}

		ovalFixtureDir := "ubuntu"
		softwareFixtureDir := filepath.Join("ubuntu", "software")
		for _, v := range systems {
			withTestFixutre(v, ovalFixtureDir, softwareFixtureDir, vulnPath, ds, func(h *fleet.Host) {
				_, err := Analyze(ctx, ds, v, vulnPath, true)
				require.NoError(t, err)

				p := NewPlatform(v.Platform, v.Name)
				assertVulns(ds, vulnPath, h, p, t)
			}, t)
		}
	})

	t.Run("#vulnsDelta", func(t *testing.T) {
		t.Run("no existing vulnerabilities", func(t *testing.T) {
			var found []fleet.SoftwareVulnerability
			var existing []fleet.SoftwareVulnerability

			toInsert, toDelete := vulnsDelta(found, existing)
			require.Empty(t, toInsert)
			require.Empty(t, toDelete)
		})

		t.Run("existing match found", func(t *testing.T) {
			found := []fleet.SoftwareVulnerability{
				{CPEID: 1, CVE: "cve_1"},
				{CPEID: 1, CVE: "cve_2"},
				{CPEID: 2, CVE: "cve_3"},
				{CPEID: 2, CVE: "cve_4"},
			}

			existing := []fleet.SoftwareVulnerability{
				{CPEID: 1, CVE: "cve_1"},
				{CPEID: 1, CVE: "cve_2"},
				{CPEID: 2, CVE: "cve_3"},
				{CPEID: 2, CVE: "cve_4"},
			}

			toInsert, toDelete := vulnsDelta(found, existing)
			require.Empty(t, toInsert)
			require.Empty(t, toDelete)
		})

		t.Run("existing differ from found", func(t *testing.T) {
			found := []fleet.SoftwareVulnerability{
				{CPEID: 1, CVE: "cve_1"},
				{CPEID: 1, CVE: "cve_2"},
				{CPEID: 3, CVE: "cve_5"},
				{CPEID: 3, CVE: "cve_6"},
			}

			existing := []fleet.SoftwareVulnerability{
				{CPEID: 1, CVE: "cve_1"},
				{CPEID: 1, CVE: "cve_2"},
				{CPEID: 2, CVE: "cve_3"},
				{CPEID: 2, CVE: "cve_4"},
			}

			expectedToInsert := []fleet.SoftwareVulnerability{
				{CPEID: 3, CVE: "cve_5"},
				{CPEID: 3, CVE: "cve_6"},
			}

			expectedToDelete := []fleet.SoftwareVulnerability{
				{CPEID: 2, CVE: "cve_3"},
				{CPEID: 2, CVE: "cve_4"},
			}

			toInsert, toDelete := vulnsDelta(found, existing)
			require.Equal(t, expectedToInsert, toInsert)
			require.ElementsMatch(t, expectedToDelete, toDelete)
		})

		t.Run("nothing found but vulns exist", func(t *testing.T) {
			var found []fleet.SoftwareVulnerability

			existing := []fleet.SoftwareVulnerability{
				{CPEID: 1, CVE: "cve_1"},
				{CPEID: 1, CVE: "cve_2"},
				{CPEID: 2, CVE: "cve_3"},
				{CPEID: 2, CVE: "cve_4"},
			}

			toInsert, toDelete := vulnsDelta(found, existing)
			require.Empty(t, toInsert)
			require.ElementsMatch(t, existing, toDelete)
		})
	})

	t.Run("#load", func(t *testing.T) {
		t.Run("invalid vuln path", func(t *testing.T) {
			platform := NewPlatform("ubuntu", "Ubuntu 20.4.0")
			_, err := loadDef(platform, "")
			require.Error(t, err, "invalid vulnerabity path")
		})
	})

	t.Run("#latestOvalDefFor", func(t *testing.T) {
		t.Run("definition matching platform for date exists", func(t *testing.T) {
			path := t.TempDir()

			today := time.Now()
			platform := NewPlatform("ubuntu", "Ubuntu 20.4.0")
			def := filepath.Join(path, platform.ToFilename(today, "json"))

			f1, err := os.Create(def)
			require.NoError(t, err)
			f1.Close()

			result, err := latestOvalDefFor(platform, path, today)
			require.NoError(t, err)
			require.Equal(t, def, result)
		})

		t.Run("definition matching platform exists but not for date", func(t *testing.T) {
			path := t.TempDir()

			today := time.Now()
			yesterday := today.Add(-24 * time.Hour)

			platform := NewPlatform("ubuntu", "Ubuntu 20.4.0")
			def := filepath.Join(path, platform.ToFilename(yesterday, "json"))

			f1, err := os.Create(def)
			require.NoError(t, err)
			f1.Close()

			result, err := latestOvalDefFor(platform, path, today)
			require.NoError(t, err)
			require.Equal(t, def, result)
		})

		t.Run("definition does not exists for platform", func(t *testing.T) {
			path := t.TempDir()

			today := time.Now()

			platform1 := NewPlatform("ubuntu", "Ubuntu 20.4.0")
			def1 := filepath.Join(path, platform1.ToFilename(today, "json"))
			f1, err := os.Create(def1)
			require.NoError(t, err)
			f1.Close()

			platform2 := NewPlatform("ubuntu", "Ubuntu 18.4.0")

			_, err = latestOvalDefFor(platform2, path, today)
			require.Error(t, err, "file not found for platform")
		})
	})
}
