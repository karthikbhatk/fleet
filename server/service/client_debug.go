package service

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/fleetdm/fleet/v4/server/fleet"
)

// DebugPprof calls the /debug/pprof/ endpoints.
func (c *Client) DebugPprof(name string) ([]byte, error) {
	endpoint := "/debug/pprof/" + name
	response, err := c.AuthenticatedDo("GET", endpoint, "", nil)
	if err != nil {
		return nil, fmt.Errorf("GET %s: %w", endpoint, err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"get pprof received status %d",
			response.StatusCode,
		)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("read pprof response body: %w", err)
	}

	return body, nil
}

func (c *Client) DebugMigrations() (*fleet.MigrationStatus, error) {
	var migrationStatus fleet.MigrationStatus
	err := c.authenticatedRequest(nil, "GET", "/debug/migrations", &migrationStatus)
	if err != nil {
		return nil, err
	}
	return &migrationStatus, nil
}

// DebugErrors calls the /debug/errors endpoint and on success writes its
// (potentially large) response body to w.
func (c *Client) DebugErrors(w io.Writer) error {
	endpoint := "/debug/errors"
	response, err := c.AuthenticatedDo("GET", endpoint, "", nil)
	if err != nil {
		return fmt.Errorf("GET %s: %w", endpoint, err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("get errors received status %d", response.StatusCode)
	}

	if _, err := io.Copy(w, response.Body); err != nil {
		return fmt.Errorf("read errors response body: %w", err)
	}
	return nil
}
