package msrc_parsed

import (
	"time"
)

type VulnGraph struct {
	// The name portion of the full product name. We will have one graph per 'product' name (e.g. Windows 10)
	ProductName string
	// When was the graph last updated.
	LastUpdated time.Time
	// All products contained in this graph (Product ID => Product full name).
	// We can have many different 'products' under a single name, for example, for 'Windows 10':
	// - Windows 10 Version 1809 for 32-bit Systems
	// - Windows 10 Version 1909 for x64-based Systems
	Products map[string]string
	// All vulnerabilities contained in this graph, by CVE
	Vulnerabities map[string]VulnNode
	// All vendor fixes for remediating the vulnerabilities contained in this graph, by KBID
	VendorFixes map[string]VendorFixNode
}

func NewVulnGraph(pName string) *VulnGraph {
	return &VulnGraph{
		ProductName:   pName,
		LastUpdated:   time.Now().UTC(),
		Products:      make(map[string]string),
		Vulnerabities: make(map[string]VulnNode),
		VendorFixes:   make(map[string]VendorFixNode),
	}
}

type NodeRef struct {
	RefName  string
	RefValue string
}

type VulnNode struct {
	Published *time.Time
	// What products are susceptible to this vuln.
	ProductsIDs []string
	// References to what Vendor fixes remediate this vuln.
	RemediatedBy []NodeRef
}

func NewVendorFixNodeRef(val string) NodeRef {
	return NodeRef{
		RefName:  "vendor_fixes",
		RefValue: val,
	}
}

type VendorFixNode struct {
	FixedBuild        string
	TargetProductsIDs []string
	// Reference to what vendor fix this particular vendor fix 'replaces'.
	Supersedes NodeRef
}
