package msrc

import (
	"testing"

	"github.com/fleetdm/fleet/v4/server/fleet"
	msrc_io "github.com/fleetdm/fleet/v4/server/vulnerabilities/msrc/io"
	"github.com/stretchr/testify/require"
)

type testData struct {
	remoteList          map[msrc_io.SecurityBulletinName]string
	remoteListError     error
	remoteDownloaded    []string
	remoteDownloadError error
	localList           []msrc_io.SecurityBulletinName
	localListError      error
	localDeleted        []msrc_io.SecurityBulletinName
	localDeleteError    error
}

type ghMock struct{ testData *testData }

func (gh ghMock) Bulletins() (map[msrc_io.SecurityBulletinName]string, error) {
	return gh.testData.remoteList, gh.testData.remoteDownloadError
}

func (gh ghMock) Download(b msrc_io.SecurityBulletinName, url string) error {
	gh.testData.remoteDownloaded = append(gh.testData.remoteDownloaded, url)
	return gh.testData.remoteDownloadError
}

type fsMock struct{ testData *testData }

func (fs fsMock) Bulletins() ([]msrc_io.SecurityBulletinName, error) {
	return fs.testData.localList, fs.testData.localListError
}

func (fs fsMock) Delete(d msrc_io.SecurityBulletinName) error {
	fs.testData.localDeleted = append(fs.testData.localDeleted, d)
	return fs.testData.localDeleteError
}

func TestSync(t *testing.T) {
	t.Run("#sync", func(t *testing.T) {
		os := []fleet.OperatingSystem{
			{
				Name:          "Microsoft Windows 11 Enterprise",
				Version:       "21H2",
				Arch:          "64-bit",
				KernelVersion: "10.0.22000.795",
			},
			{
				Name:          "Microsoft Windows 10 Pro",
				Version:       "10.0.19044",
				Arch:          "64-bit",
				KernelVersion: "10.0.19044",
			},
		}

		testData := testData{
			remoteList: map[msrc_io.SecurityBulletinName]string{
				msrc_io.NewSecurityBulletinName("Windows_10-2022_10_10.json"): "http://somebulletin.com",
			},
			localList: []msrc_io.SecurityBulletinName{"Windows_10-2022_09_10.json"},
		}

		err := sync(os, fsMock{testData: &testData}, ghMock{testData: &testData})
		require.NoError(t, err)
		require.ElementsMatch(t, testData.remoteDownloaded, []string{"http://somebulletin.com"})
		require.ElementsMatch(t, testData.localDeleted, []msrc_io.SecurityBulletinName{"Windows_10-2022_09_10.json"})
	})

	t.Run("#bulletinsDelta", func(t *testing.T) {
		t.Run("win OS provided", func(t *testing.T) {
			os := []fleet.OperatingSystem{
				{
					Name:          "Microsoft Windows 11 Enterprise",
					Version:       "21H2",
					Arch:          "64-bit",
					KernelVersion: "10.0.22000.795",
				},
				{
					Name:          "Microsoft Windows 10 Pro",
					Version:       "10.0.19044",
					Arch:          "64-bit",
					KernelVersion: "10.0.19044",
				},
			}
			t.Run("without remote bulletins", func(t *testing.T) {
				var remote []msrc_io.SecurityBulletinName
				local := []msrc_io.SecurityBulletinName{
					"Windows_10-2022_10_10.json",
				}
				toDownload, toDelete := bulletinsDelta(os, local, remote)
				require.Empty(t, toDownload)
				require.Empty(t, toDelete)
			})

			t.Run("with remote bulletins", func(t *testing.T) {
				remote := []msrc_io.SecurityBulletinName{
					"Windows_10-2022_10_10.json",
					"Windows_11-2022_10_10.json",
					"Windows_Server_2016-2022_10_10.json",
					"Windows_8.1-2022_10_10.json",
				}
				t.Run("no local bulletins", func(t *testing.T) {
					var local []msrc_io.SecurityBulletinName
					toDownload, toDelete := bulletinsDelta(os, local, remote)

					require.ElementsMatch(t, toDownload, []msrc_io.SecurityBulletinName{
						"Windows_10-2022_10_10.json",
						"Windows_11-2022_10_10.json",
					})
					require.Empty(t, toDelete)
				})

				t.Run("missing some local bulletin", func(t *testing.T) {
					local := []msrc_io.SecurityBulletinName{
						"Windows_10-2022_10_10.json",
					}
					toDownload, toDelete := bulletinsDelta(os, local, remote)

					require.ElementsMatch(t, toDownload, []msrc_io.SecurityBulletinName{
						"Windows_11-2022_10_10.json",
					})
					require.Empty(t, toDelete)
				})

				t.Run("out of date local bulletin", func(t *testing.T) {
					local := []msrc_io.SecurityBulletinName{
						"Windows_10-2022_09_10.json",
						"Windows_11-2022_10_10.json",
					}

					toDownload, toDelete := bulletinsDelta(os, local, remote)

					require.ElementsMatch(t, toDownload, []msrc_io.SecurityBulletinName{
						"Windows_10-2022_10_10.json",
					})
					require.ElementsMatch(t, toDelete, []msrc_io.SecurityBulletinName{
						"Windows_10-2022_09_10.json",
					})
				})

				t.Run("up to date local bulletins", func(t *testing.T) {
					local := []msrc_io.SecurityBulletinName{
						"Windows_10-2022_10_10.json",
						"Windows_11-2022_10_10.json",
					}

					toDownload, toDelete := bulletinsDelta(os, local, remote)

					require.Empty(t, toDownload)
					require.Empty(t, toDelete)
				})
			})
		})

		t.Run("no Win OS provided", func(t *testing.T) {
			os := []fleet.OperatingSystem{
				{
					Name:          "CentOS",
					Version:       "8.0.0",
					Platform:      "rhel",
					KernelVersion: "5.10.76-linuxkit",
				},
			}
			local := []msrc_io.SecurityBulletinName{"Windows_11-2022_10_10.json"}
			remote := []msrc_io.SecurityBulletinName{"Windows_10-2022_10_10.json"}

			t.Run("nothing to download, nothing to delete", func(t *testing.T) {
				toDownload, toDelete := bulletinsDelta(os, local, remote)
				require.Empty(t, toDownload)
				require.Empty(t, toDelete)
			})
		})

		t.Run("no OS provided", func(t *testing.T) {
			var os []fleet.OperatingSystem
			t.Run("no local bulletins", func(t *testing.T) {
				var local []msrc_io.SecurityBulletinName

				t.Run("returns all remote", func(t *testing.T) {
					remote := []msrc_io.SecurityBulletinName{
						"Windows_10-2022_10_10.json",
						"Windows_11-2022_10_10.json",
					}

					toDownload, toDelete := bulletinsDelta(os, local, remote)
					require.ElementsMatch(t, toDownload, remote)
					require.Empty(t, toDelete)
				})
			})
		})
	})
}
