package fleet

import "time"

// VPPApp represents a VPP (Volume Purchase Program) application,
// this is used by Apple MDM to manage applications via Apple
// Bussines Manager.
type VPPApp struct {
	// AdamID is a unique identifier assigned to each app in
	// the App Store, this value is managed by Apple.
	AdamID string `db:"adam_id" json:"app_store_id"`
	// BundleIdentifier is the unique bundle identifier of the
	// Application.
	BundleIdentifier string `db:"bundle_identifier" json:"bundle_identifier"`
	// IconURL is the URL of this App icon
	IconURL string `db:"icon_url" json:"icon_url"`
	// Name is the user-facing name of this app.
	Name string `db:"name" json:"name"`
	// LatestVersion is the latest version of this app.
	LatestVersion string `db:"latest_version" json:"latest_version"`
	// Added indicates whether or not this app has been added to Fleet.
	Added   bool  `json:"added"`
	TeamID  *uint `db:"-" json:"-"`
	TitleID uint  `db:"title_id" json:"-"`

	CreatedAt time.Time `db:"created_at" json:"-"`
	UpdatedAt time.Time `db:"updated_at" json:"-"`
}

// AuthzType implements authz.AuthzTyper.
func (v *VPPApp) AuthzType() string {
	return "installable_entity"
}
