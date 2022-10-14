package fleet

import (
	"encoding/json"
	"time"
)

type StatisticsPayload struct {
	AnonymousIdentifier            string                             `json:"anonymousIdentifier"`
	FleetVersion                   string                             `json:"fleetVersion"`
	LicenseTier                    string                             `json:"licenseTier"`
	Organization                   string                             `json:"organization"`
	NumHostsEnrolled               int                                `json:"numHostsEnrolled"`
	NumUsers                       int                                `json:"numUsers"`
	NumTeams                       int                                `json:"numTeams"`
	NumPolicies                    int                                `json:"numPolicies"`
	NumLabels                      int                                `json:"numLabels"`
	SoftwareInventoryEnabled       bool                               `json:"softwareInventoryEnabled"`
	VulnDetectionEnabled           bool                               `json:"vulnDetectionEnabled"`
	SystemUsersEnabled             bool                               `json:"systemUsersEnabled"`
	HostsStatusWebHookEnabled      bool                               `json:"hostsStatusWebHookEnabled"`
	NumWeeklyActiveUsers           int                                `json:"numWeeklyActiveUsers"`
	HostsEnrolledByOperatingSystem map[string][]HostsCountByOSVersion `json:"hostsEnrolledByOperatingSystem"`
	// HostsEnrolledByOrbitVersion is a count of hosts enrolled to Fleet grouped by orbit version
	HostsEnrolledByOrbitVersion []HostsCountByOrbitVersion `json:"hostsEnrolledByOrbitVersion"`
	// HostsEnrolledByOsqueryVersion is a count of hosts enrolled to Fleet grouped by osquery version
	HostsEnrolledByOsqueryVersion []HostsCountByOsqueryVersion `json:"hostsEnrolledByOsqueryVersion"`
	StoredErrors                  json.RawMessage              `json:"storedErrors"`
	// NumHostsNotResponding is a count of hosts that connect to Fleet successfully but fail to submit results for distributed queries.
	NumHostsNotResponding int `json:"numHostsNotResponding"`
}

type HostsCountByOrbitVersion struct {
	OrbitVersion string `json:"orbitVersion" db:"orbit_version"`
	NumHosts     int    `json:"numHosts" db:"num_hosts"`
}
type HostsCountByOsqueryVersion struct {
	OsqueryVersion string `json:"osqueryVersion" db:"osquery_version"`
	NumHosts       int    `json:"numHosts" db:"num_hosts"`
}

type HostsCountByOSVersion struct {
	Version     string `json:"version"`
	NumEnrolled int    `json:"numEnrolled"`
}

const (
	StatisticsFrequency = time.Hour * 24 * 7
)
