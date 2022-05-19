package fleet

import (
	"fmt"
	"time"
)

type SoftwareCVE struct {
	CVE         string `json:"cve" db:"cve"`
	DetailsLink string `json:"details_link" db:"details_link"`
}

// Software is a named and versioned piece of software installed on a device.
type Software struct {
	ID uint `json:"id" db:"id"`
	// Name is the reported name.
	Name string `json:"name" db:"name"`
	// Version is reported version.
	Version string `json:"version" db:"version"`
	// BundleIdentifier is the CFBundleIdentifier label from the info properties
	BundleIdentifier string `json:"bundle_identifier,omitempty" db:"bundle_identifier"`
	// Source is the source of the data (osquery table name).
	Source string `json:"source" db:"source"`

	// Release is the version of the OS this software was released on
	// (e.g. "30.el7" for a CentOS package).
	Release string `json:"release,omitempty" db:"release"`
	// Vendor is the supplier of the software (e.g. "CentOS").
	Vendor string `json:"vendor,omitempty" db:"vendor"`
	// Arch is the architecture of the software (e.g. "x86_64").
	Arch string `json:"arch,omitempty" db:"arch"`

	// GenerateCPE is the CPE23 string that corresponds to the current software
	GenerateCPE string `json:"generated_cpe" db:"generated_cpe"`

	// GeneratedCPEID is the ID of the matched CPE
	GeneratedCPEID uint `json:"-" db:"generated_cpe_id"`

	// Vulnerabilities lists all the found CVEs for the CPE
	Vulnerabilities VulnerabilitiesSlice `json:"vulnerabilities"`
	// HostsCount indicates the number of hosts with that software, filled only
	// if explicitly requested.
	HostsCount int `json:"hosts_count,omitempty" db:"hosts_count"`
	// CountsUpdatedAt is the timestamp when the hosts count was last updated
	// for that software, filled only if hosts count is requested.
	CountsUpdatedAt time.Time `json:"-" db:"counts_updated_at"`
	// LastOpenedAt is the timestamp when that software was last opened on the
	// corresponding host. Only filled when the software list is requested for
	// a specific host (host_id is provided).
	LastOpenedAt *time.Time `json:"last_opened_at,omitempty" db:"last_opened_at"`
}

func (Software) AuthzType() string {
	return "software"
}

// AuthzSoftwareInventory is used for access controls on software inventory.
type AuthzSoftwareInventory struct {
	// TeamID is the ID of the team. A value of nil means global scope.
	TeamID *uint `json:"team_id"`
}

// AuthzType implements authz.AuthzTyper.
func (s *AuthzSoftwareInventory) AuthzType() string {
	return "software_inventory"
}

type VulnerabilitiesSlice []SoftwareCVE

// HostSoftware is the set of software installed on a specific host
type HostSoftware struct {
	// Software is the software information.
	Software []Software `json:"software,omitempty" csv:"-"`
}

type SoftwareIterator interface {
	Next() bool
	Value() (*Software, error)
	Err() error
	Close() error
}

type SoftwareListOptions struct {
	ListOptions

	TeamID         *uint `query:"team_id,optional"`
	VulnerableOnly bool  `query:"vulnerable,optional"`

	SkipLoadingCVEs bool

	// WithHostCounts indicates that the list of software should include the
	// counts of hosts per software, and include only those software that have
	// a count of hosts > 0.
	WithHostCounts bool
}

// SoftwareVulnerability identifies a vulnerability on a specific software (CPE).
type SoftwareVulnerability struct {
	ID    uint   `db:"software_id"`
	CPE   string `db:"cpe"`
	CPEID uint   `db:"cpe_id"`
	CVE   string `db:"cve"`
}

// String implements fmt.Stringer.
func (sv SoftwareVulnerability) String() string {
	return fmt.Sprintf("{%d,%s}", sv.CPEID, sv.CVE)
}

// SoftwareWithCPE holds a software piece alongside its CPE ID.
type SoftwareWithCPE struct {
	// Software holds the software data.
	Software
	// CPEID is the ID of the software CPE in the system.
	CPEID uint
}

type VulnerabilitySource int

const (
	NVD VulnerabilitySource = iota
	OVAL
)
