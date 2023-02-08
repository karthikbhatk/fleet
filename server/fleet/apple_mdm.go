package fleet

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/micromdm/nanodep/godep"
	"github.com/micromdm/nanomdm/mdm"
	"go.mozilla.org/pkcs7"
	"howett.net/plist"
)

// MDMAppleEnrollmentType is the type for Apple MDM enrollments.
type MDMAppleEnrollmentType string

const (
	// MDMAppleEnrollmentTypeAutomatic is the value for automatic enrollments.
	MDMAppleEnrollmentTypeAutomatic MDMAppleEnrollmentType = "automatic"
	// MDMAppleEnrollmentTypeManual is the value for manual enrollments.
	MDMAppleEnrollmentTypeManual MDMAppleEnrollmentType = "manual"
)

// MDMAppleEnrollmentProfilePayload contains the data necessary to create
// an enrollment profile in Fleet.
type MDMAppleEnrollmentProfilePayload struct {
	// Type is the type of the enrollment.
	Type MDMAppleEnrollmentType `json:"type"`
	// DEPProfile is the JSON object with the following Apple-defined fields:
	// https://developer.apple.com/documentation/devicemanagement/profile
	//
	// DEPProfile is nil when Type is MDMAppleEnrollmentTypeManual.
	DEPProfile *json.RawMessage `json:"dep_profile"`
	// Token should be auto-generated.
	Token string `json:"-"`
}

// MDMAppleEnrollmentProfile represents an Apple MDM enrollment profile in Fleet.
// Such enrollment profiles are used to enroll Apple devices to Fleet.
type MDMAppleEnrollmentProfile struct {
	// ID is the unique identifier of the enrollment in Fleet.
	ID uint `json:"id" db:"id"`
	// Token is a random identifier for an enrollment. Currently as the authentication
	// token to protect access to the enrollment.
	Token string `json:"token" db:"token"`
	// Type is the type of the enrollment.
	Type MDMAppleEnrollmentType `json:"type" db:"type"`
	// DEPProfile is the JSON object with the following Apple-defined fields:
	// https://developer.apple.com/documentation/devicemanagement/profile
	//
	// DEPProfile is nil when Type is MDMAppleEnrollmentTypeManual.
	DEPProfile *json.RawMessage `json:"dep_profile" db:"dep_profile"`
	// EnrollmentURL is the URL where an enrollement is served.
	EnrollmentURL string `json:"enrollment_url" db:"-"`

	UpdateCreateTimestamps
}

// AuthzType implements authz.AuthzTyper.
func (m MDMAppleEnrollmentProfile) AuthzType() string {
	return "mdm_apple_enrollment_profile"
}

// MDMAppleDEPKeyPair contains the DEP public key certificate and private key pair. Both are PEM encoded.
type MDMAppleDEPKeyPair struct {
	PublicKey  []byte `json:"public_key"`
	PrivateKey []byte `json:"private_key"`
}

// MDMAppleCommandResult holds the result of a command execution provided by the target device.
type MDMAppleCommandResult struct {
	// ID is the enrollment ID. This should be the same as the device ID.
	ID string `json:"id" db:"id"`
	// CommandUUID is the unique identifier of the command.
	CommandUUID string `json:"command_uuid" db:"command_uuid"`
	// Status is the command status. One of Acknowledged, Error, or NotNow.
	Status string `json:"status" db:"status"`
	// Result is the original command result XML plist. If the status is Error, it will include the
	// ErrorChain key with more information.
	Result []byte `json:"result" db:"result"`
}

// AuthzType implements authz.AuthzTyper.
func (m MDMAppleCommandResult) AuthzType() string {
	return "mdm_apple_command_result"
}

// MDMAppleInstaller holds installer packages for Apple devices.
type MDMAppleInstaller struct {
	// ID is the unique identifier of the installer in Fleet.
	ID uint `json:"id" db:"id"`
	// Name is the name of the installer (usually the package file name).
	Name string `json:"name" db:"name"`
	// Size is the size of the installer package.
	Size int64 `json:"size" db:"size"`
	// Manifest is the manifest of the installer. Generated from the installer
	// contents and ready to use in `InstallEnterpriseApplication` commands.
	Manifest string `json:"manifest" db:"manifest"`
	// Installer is the actual installer contents.
	Installer []byte `json:"-" db:"installer"`
	// URLToken is a random token used for authentication to protect access to installers.
	// Applications deployede via InstallEnterpriseApplication must be publicly accessible,
	// this hard to guess token provides some protection.
	URLToken string `json:"url_token" db:"url_token"`
	// URL is the full URL where the installer is served.
	URL string `json:"url"`
}

// AuthzType implements authz.AuthzTyper.
func (m MDMAppleInstaller) AuthzType() string {
	return "mdm_apple_installer"
}

// MDMAppleDevice represents an MDM enrolled Apple device.
type MDMAppleDevice struct {
	// ID is the device hardware UUID.
	ID string `json:"id" db:"id"`
	// SerialNumber is the serial number of the Apple device.
	SerialNumber string `json:"serial_number" db:"serial_number"`
	// Enabled indicates whether the device is currently enrolled.
	// It's set to false when a device unenrolls from Fleet.
	Enabled bool `json:"enabled" db:"enabled"`
}

// AuthzType implements authz.AuthzTyper.
func (m MDMAppleDevice) AuthzType() string {
	return "mdm_apple_device"
}

// MDMAppleDEPDevice represents an Apple device in Apple Business Manager (ABM).
type MDMAppleDEPDevice struct {
	godep.Device
}

// AuthzType implements authz.AuthzTyper.
func (m MDMAppleDEPDevice) AuthzType() string {
	return "mdm_apple_dep_device"
}

// These following types are copied from nanomdm.

// EnrolledAPIResult is a per-enrollment API result.
type EnrolledAPIResult struct {
	PushError    string `json:"push_error,omitempty"`
	PushResult   string `json:"push_result,omitempty"`
	CommandError string `json:"command_error,omitempty"`
}

// EnrolledAPIResults is a map of enrollments to a per-enrollment API result.
type EnrolledAPIResults map[string]*EnrolledAPIResult

// CommandEnqueueResult is the result of a command execution on enrolled Apple devices.
type CommandEnqueueResult struct {
	// Status is the status of the command.
	Status EnrolledAPIResults `json:"status,omitempty"`
	// NoPush indicates whether the command was issued with no_push.
	// If this is true, then Fleet won't send a push notification to devices.
	NoPush bool `json:"no_push,omitempty"`
	// PushError indicates the error when trying to send push notification
	// to target devices.
	PushError string `json:"push_error,omitempty"`
	// CommandError holds the error when enqueueing the command.
	CommandError string `json:"command_error,omitempty"`
	// CommandUUID is the unique identifier for the command.
	CommandUUID string `json:"command_uuid,omitempty"`
	// RequestType is the name of the command.
	RequestType string `json:"request_type,omitempty"`
}

// MDMAppleCommand represents an Apple MDM command.
type MDMAppleCommand struct {
	*mdm.Command
}

// AuthzType implements authz.AuthzTyper.
func (m MDMAppleCommand) AuthzType() string {
	return "mdm_apple_command"
}

// MDMAppleHostDetails represents the device identifiers used to ingest an MDM device as a Fleet
// host pending enrollment.
// See also https://developer.apple.com/documentation/devicemanagement/authenticaterequest.
type MDMAppleHostDetails struct {
	SerialNumber string
	UDID         string
	Model        string
}

type MDMAppleCommandTimeoutError struct{}

func (e MDMAppleCommandTimeoutError) Error() string {
	return "Timeout waiting for MDM device to acknowledge command"
}

func (e MDMAppleCommandTimeoutError) StatusCode() int {
	return http.StatusGatewayTimeout
}

// Mobileconfig is the byte slice corresponding to an XML property list (i.e. plist) representation
// of an Apple MDM configuration profile in Fleet.
//
// Configuration profiles are used to configure Apple devices. See also
// https://developer.apple.com/documentation/devicemanagement/configuring_multiple_devices_using_profiles.
type Mobileconfig []byte

// ParseConfigProfile attempts to parse the Mobileconfig byte slice as a Fleet MDMAppleConfigProfile.
//
// The byte slice must be XML or PKCS7 parseable. Fleet also requires that it contains both
// a PayloadIdentifier and a PayloadDisplayName and that it has PayloadType set to "Configuration".
//
// Adapted from https://github.com/micromdm/micromdm/blob/main/platform/profile/profile.go
func (mc *Mobileconfig) ParseConfigProfile() (*MDMAppleConfigProfile, error) {
	mcBytes := *mc
	if !bytes.HasPrefix(mcBytes, []byte("<?xml")) {
		p7, err := pkcs7.Parse(mcBytes)
		if err != nil {
			return nil, fmt.Errorf("mobileconfig is not XML nor PKCS7 parseable: %w", err)
		}
		err = p7.Verify()
		if err != nil {
			return nil, err
		}
		mcBytes = Mobileconfig(p7.Content)
	}
	var parsed struct {
		PayloadIdentifier  string
		PayloadDisplayName string
		PayloadType        string
	}
	_, err := plist.Unmarshal(mcBytes, &parsed)
	if err != nil {
		return nil, err
	}
	if parsed.PayloadType != "Configuration" {
		return nil, fmt.Errorf("invalid PayloadType: %s", parsed.PayloadType)
	}
	if parsed.PayloadIdentifier == "" {
		return nil, errors.New("empty PayloadIdentifier in profile")
	}
	if parsed.PayloadDisplayName == "" {
		return nil, errors.New("empty PayloadDisplayName in profile")
	}

	return &MDMAppleConfigProfile{
		Identifier:   parsed.PayloadIdentifier,
		Name:         parsed.PayloadDisplayName,
		Mobileconfig: mc,
	}, nil
}

// MDMAppleConfigProfile represents an Apple MDM configuration profile in Fleet.
// Configuration profiles are used to configure Apple devices .
// See also https://developer.apple.com/documentation/devicemanagement/configuring_multiple_devices_using_profiles.
type MDMAppleConfigProfile struct {
	// ProfileID is the unique id of the configuration profile in Fleet
	ProfileID uint `db:"profile_id"`
	// TeamID is the id of the team with which the configuration is associated. A team id of zero
	// represents a configuration profile that is not associated with any team.
	TeamID uint `db:"team_id"`
	// Identifier corresponds to the payload identifier of the associated mobileconfig payload.
	// Fleet requires that Identifier must be unique in combination with the Name and TeamID.
	Identifier string `db:"identifier"`
	// Name corresponds to the payload display name of the associated mobileconfig payload.
	// Fleet requires that Name must be unique in combination with the Identifier and TeamID.
	Name string `db:"name"`
	// Mobileconfig is the byte slice corresponding to the XML property list (i.e. plist)
	// representation of the configuration profile. It must be XML or PKCS7 parseable.
	Mobileconfig *Mobileconfig `db:"mobileconfig"`
	CreatedAt    time.Time     `db:"created_at"`
	UpdatedAt    time.Time     `db:"updated_at"`
}

func (cp *MDMAppleConfigProfile) Validate() error {
	// TODO(sarah): Additional validations for PayloadContent (e.g., screening out FileVault payloads)
	// should be handled here

	return nil
}
