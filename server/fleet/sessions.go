package fleet

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/fleetdm/fleet/v4/server/ptr"
)

// Auth contains methods to fetch information from a valid SSO auth response
type Auth interface {
	// UserID returns the Subject Name Identifier associated with the request,
	// this can be an email address, an entity identifier, or any other valid
	// Name Identifier as described in the spec:
	// http://docs.oasis-open.org/security/saml/v2.0/saml-core-2.0-os.pdf
	//
	// Fleet requires users to configure this value to be the email of the Subject
	UserID() string
	// UserDisplayName finds a display name in the SSO response Attributes, there
	// isn't a defined spec for this, so the return value is in a best-effort
	// basis
	UserDisplayName() string
	// RequestID returns the request id associated with this SSO session
	RequestID() string
	// AssertionAttributes returns the attributes of the SAML response.
	AssertionAttributes() []SAMLAttribute
}

// SAMLAttribute holds the name and values of a custom attribute.
type SAMLAttribute struct {
	Name   string
	Values []SAMLAttributeValue
}

// SAMLAttributeValue holds the type and value of a custom attribute.
type SAMLAttributeValue struct {
	// Type is the type of attribute value.
	Type string
	// Value is the actual value of the attribute.
	Value string
}

type SSOSession struct {
	Token       string
	RedirectURL string
}

// SessionSSOSettings SSO information used prior to authentication.
type SessionSSOSettings struct {
	// IDPName is a human readable name for the IDP
	IDPName string `json:"idp_name"`
	// IDPImageURL https link to a logo image for the IDP.
	IDPImageURL string `json:"idp_image_url"`
	// SSOEnabled true if single sign on is enabled.
	SSOEnabled bool `json:"sso_enabled"`
}

// Session is the model object which represents what an active session is
type Session struct {
	CreateTimestamp
	ID         uint
	AccessedAt time.Time `db:"accessed_at"`
	UserID     uint      `json:"user_id" db:"user_id"`
	Key        string
	APIOnly    *bool `json:"-" db:"api_only"`
}

func (s Session) AuthzType() string {
	return "session"
}

// SSORolesInfo holds the configuration parsed from SAML custom attributes.
//
// `Global` and `Teams` are never both set (by design, users must be either global
// or member of teams).
type SSORolesInfo struct {
	// Global holds the role for the Global domain.
	Global *string
	// Teams holds the roles for teams.
	Teams []TeamRole
}

// TeamRole holds a user's role on a Team.
type TeamRole struct {
	// ID is the unique identifier of the team.
	ID uint
	// Role is the role of the user in the team.
	Role string
}

func (s SSORolesInfo) verify() error {
	if s.Global != nil && len(s.Teams) > 0 {
		return errors.New("cannot set both global and team roles")
	}
	// Check for duplicate entries for the same team.
	// This is just in case some IdP allows duplicating attributes.
	teamMap := make(map[uint]struct{})
	for _, teamRole := range s.Teams {
		if _, ok := teamMap[teamRole.ID]; ok {
			return fmt.Errorf("duplicate team entry: %d", teamRole.ID)
		}
		teamMap[teamRole.ID] = struct{}{}
	}
	return nil
}

func (s SSORolesInfo) empty() bool {
	return s.Global == nil && len(s.Teams) == 0
}

const (
	globalUserRoleSSOAttrName     = "FLEET_JIT_USER_ROLE_GLOBAL"
	teamUserRoleSSOAttrNamePrefix = "FLEET_JIT_USER_ROLE_TEAM_"
)

// RolesFromSSOAttributes loads Global and Team roles from SAML custom attributes.
//   - Custom attribute `FLEET_JIT_USER_ROLE_GLOBAL` is used for setting global role.
//   - Custom attributes of the form `FLEET_JIT_USER_ROLE_TEAM_<TEAM_ID>` are used
//     for setting role for a team with ID <TEAM_ID>.
//
// For both attributes currently supported values are `admin`, `maintainer` and `observer`
func RolesFromSSOAttributes(attributes []SAMLAttribute) (SSORolesInfo, error) {
	ssoRoleInfo := SSORolesInfo{}
	for _, attribute := range attributes {
		switch {
		case attribute.Name == globalUserRoleSSOAttrName:
			role, err := parseRole(attribute.Values)
			if err != nil {
				return SSORolesInfo{}, fmt.Errorf("parse global role: %w", err)
			}
			ssoRoleInfo.Global = ptr.String(role)
		case strings.HasPrefix(attribute.Name, teamUserRoleSSOAttrNamePrefix):
			teamIDSuffix := strings.TrimPrefix(attribute.Name, teamUserRoleSSOAttrNamePrefix)
			teamID, err := strconv.ParseUint(teamIDSuffix, 10, 64)
			if err != nil {
				return SSORolesInfo{}, fmt.Errorf("parse team ID: %w", err)
			}
			teamRole, err := parseRole(attribute.Values)
			if err != nil {
				return SSORolesInfo{}, fmt.Errorf("parse team role: %w", err)
			}
			ssoRoleInfo.Teams = append(ssoRoleInfo.Teams, TeamRole{
				ID:   uint(teamID),
				Role: teamRole,
			})
		default:
			continue
		}
	}
	if err := ssoRoleInfo.verify(); err != nil {
		return SSORolesInfo{}, err
	}
	if ssoRoleInfo.empty() {
		// When the configuration is not set, the default is to
		// make the user a global observer.
		return SSORolesInfo{Global: ptr.String(RoleObserver)}, nil
	}
	return ssoRoleInfo, nil
}

func parseRole(values []SAMLAttributeValue) (string, error) {
	if len(values) == 0 {
		return "", errors.New("empty role")
	}
	// Using last value by default.
	value := values[len(values)-1].Value
	if value != RoleAdmin && value != RoleMaintainer && value != RoleObserver {
		return "", fmt.Errorf("invalid role: %s", value)
	}
	return value, nil
}
