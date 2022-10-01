package main

import (
	"context"
	"strconv"

	"github.com/fleetdm/fleet/v4/orbit/pkg/build"
	orbit_table "github.com/fleetdm/fleet/v4/orbit/pkg/table"
	"github.com/fleetdm/fleet/v4/server/service"
	"github.com/osquery/osquery-go/plugin/table"
)

// orbitInfoExtension implements an extension table that provides info about Orbit.
type orbitInfoExtension struct {
	orbitClient     *service.OrbitClient
	deviceAuthToken string
}

var _ orbit_table.Extension = orbitInfoExtension{}

// Name partially implements orbit_table.Extension.
func (o orbitInfoExtension) Name() string {
	return "orbit_info"
}

// Columns partially implements orbit_table.Extension.
func (o orbitInfoExtension) Columns() []table.ColumnDefinition {
	return []table.ColumnDefinition{
		table.TextColumn("version"),
		table.TextColumn("device_auth_token"),
		table.TextColumn("enrolled"),
		table.TextColumn("last_request_error"),
	}
}

// GenerateFunc partially implements orbit_table.Extension.
func (o orbitInfoExtension) GenerateFunc(_ context.Context, _ table.QueryContext) ([]map[string]string, error) {
	v := build.Version
	if v == "" {
		v = "unknown"
	}
	lastRequestError := ""
	if err := o.orbitClient.LastRequestError(); err != nil {
		lastRequestError = err.Error()
	}

	return []map[string]string{
		{
			"version":            v,
			"device_auth_token":  o.deviceAuthToken,
			"enrolled":           strconv.FormatBool(o.orbitClient.Enrolled()),
			"last_request_error": lastRequestError,
		},
	}, nil
}
