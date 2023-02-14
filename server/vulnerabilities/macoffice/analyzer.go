package macoffice

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"

	"github.com/fleetdm/fleet/v4/server/contexts/ctxerr"
	"github.com/fleetdm/fleet/v4/server/fleet"
	"github.com/fleetdm/fleet/v4/server/vulnerabilities/io"
)

func latestReleaseNotes(vulnPath string) (ReleaseNotes, error) {
	fs := io.NewFSClient(vulnPath)

	files, err := fs.MacOfficeReleaseNotes()
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, nil
	}

	sort.Slice(files, func(i, j int) bool { return files[j].Before(files[i]) })
	filePath := filepath.Join(vulnPath, files[0].String())

	payload, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	relNotes := ReleaseNotes{}
	err = json.Unmarshal(payload, &relNotes)
	if err != nil {
		return nil, err
	}

	// Ensure the release notes are sorted by release date, this is because the vuln. processing
	// algo. will stop when a release note older than the current software version is found.
	sort.Slice(relNotes, func(i, j int) bool { return relNotes[j].Date.Before(relNotes[i].Date) })

	return relNotes, nil
}

func collectVulnerabilities(
	software *fleet.Software,
	relNotes ReleaseNotes,
) []fleet.SoftwareVulnerability {
	var vulns []fleet.SoftwareVulnerability

	product, ok := OfficeProductFromBundleId(software.BundleIdentifier)

	// If we don't have an Office Product
	if !ok {
		return vulns
	}

	for _, relNote := range relNotes {
		// We only care about release notes with set versions and with security updates,
		// 'relNotes' should only contain valid release notes, but this check not expensive.
		if !relNote.Valid() {
			continue
		}

		if relNote.CmpVersion(software.Version) <= 0 {
			return vulns
		}

		for _, cve := range relNote.CollectVulnerabilities(product) {
			vulns = append(vulns, fleet.SoftwareVulnerability{
				SoftwareID: software.ID,
				CVE:        cve,
			})
		}
	}

	return vulns
}

func Analyze(
	ctx context.Context,
	ds fleet.Datastore,
	vulnPath string,
	collectVulns bool,
) ([]fleet.SoftwareVulnerability, error) {
	relNotes, err := latestReleaseNotes(vulnPath)
	if err != nil {
		return nil, err
	}

	if len(relNotes) == 0 {
		return nil, nil
	}

	iter, err := ds.ListSoftwareBySourceIter(ctx, []string{"apps"})
	if err != nil {
		return nil, err
	}
	defer iter.Close()

	var collectedVulns []fleet.SoftwareVulnerability
	for iter.Next() {
		software, err := iter.Value()
		if err != nil {
			return nil, ctxerr.Wrap(ctx, err, "getting software from iterator")
		}

		detected := collectVulnerabilities(software, relNotes)

		if len(detected) != 0 {
			// Figure out what to delete and what to insert
			inserted, err := ds.InsertSoftwareVulnerabilities(ctx, detected, fleet.MacOfficeReleaseNotesSource)
			if err != nil {
				return nil, err
			}

			if collectVulns && inserted > 0 {
				collectedVulns = append(collectedVulns, detected...)
			}
		}
	}

	return collectedVulns, nil
}
