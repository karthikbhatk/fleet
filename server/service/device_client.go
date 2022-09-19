package service

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/fleetdm/fleet/v4/server/fleet"
)

// Device client is used consume the `device/...` endpoints and meant to be used by Fleet Desktop
type DeviceClient struct {
	*baseClient
	token string
}

func (dc *DeviceClient) request(verb string, path string, query string, responseDest interface{}) error {
	var bodyBytes []byte
	request, err := http.NewRequest(
		verb,
		dc.url(path, query).String(),
		bytes.NewBuffer(bodyBytes),
	)
	if err != nil {
		return err
	}

	dc.setClientCapabilitiesHeader(request)
	response, err := dc.http.Do(request)
	if err != nil {
		return fmt.Errorf("%s %s: %w", verb, path, err)
	}
	defer response.Body.Close()

	return dc.parseResponse(verb, path, response, responseDest)
}

// NewDeviceClient instantiates a new client to perform requests against device
// endpoints
func NewDeviceClient(addr, token string, insecureSkipVerify bool, rootCA string, capabilities []fleet.Capability) (*DeviceClient, error) {
	baseClient, err := newBaseClient(addr, insecureSkipVerify, rootCA, "", capabilities)
	if err != nil {
		return nil, err
	}

	return &DeviceClient{
		baseClient: baseClient,
		token:      token,
	}, nil
}

// Get fetches payload used by Fleet Desktop.
func (dc *DeviceClient) GetDesktopPayload() (*FleetDesktopResponse, error) {
	verb, path := "GET", "/api/latest/fleet/device/"+dc.token+"/desktop"

	var r FleetDesktopResponse
	err := dc.request(verb, path, "", &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
