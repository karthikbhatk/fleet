package msrc_input

import (
	"strings"
	"time"
)

// XML elements related to the 'vuln' namespace used to describe vulnerabilities and their remediations.

type VulnerabilityXML struct {
	CVE          string                        `xml:"CVE"`
	Score        float64                       `xml:"CVSSScoreSets>ScoreSet>BaseScore"`
	Revisions    []RevisionHistoryXML          `xml:"RevisionHistory>Revision"`
	Remediations []VulnerabilityRemediationXML `xml:"Remediations>Remediation"`
}

type RevisionHistoryXML struct {
	Date        string `xml:"Date"`
	Description string `xml:"Description"`
}

type VulnerabilityRemediationXML struct {
	Type            string   `xml:"Type,attr"`
	FixedBuild      string   `xml:"FixedBuild"`
	RestartRequired string   `xml:"RestartRequired"`
	ProductIDs      []string `xml:"ProductID"`
	Description     string   `xml:"Description"`
	Supercedence    string   `xml:"Supercedence"`
}

// IncludesVendorFix returns true if the vulnerability has a vendor fix targeting the product
// identified by pID.
func (r *VulnerabilityXML) IncludesVendorFix(pID string) bool {
	for _, rem := range r.Remediations {
		if rem.IsVendorFix() {
			for _, vfPID := range rem.ProductIDs {
				if vfPID == pID {
					return true
				}
			}
		}
	}

	return false
}

// PublishedDate returns the date the vuln was published (if any)
func (v *VulnerabilityXML) PublishedDate() *time.Time {
	var dPublished *time.Time

	for _, rev := range v.Revisions {
		if strings.Index(rev.Description, "Information published") != -1 {
			dPublished, err := time.Parse("2006-01-02T15:04:05", rev.Date)
			if err != nil {
				return nil
			}
			return &dPublished
		}
	}
	return dPublished
}

func (rem *VulnerabilityRemediationXML) IsVendorFix() bool {
	return rem.Type == "Vendor Fix"
}
