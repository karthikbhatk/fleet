package service

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/fleetdm/fleet/v4/server/fleet"
	"github.com/stretchr/testify/require"
)

func TestUrlGeneration(t *testing.T) {
	t.Run("without prefix", func(t *testing.T) {
		bc, err := newBaseClient("https://test.com", true, "", "", fleet.CapabilityMap{})
		require.NoError(t, err)
		require.Equal(t, "https://test.com/test/path", bc.url("test/path", "").String())
		require.Equal(t, "https://test.com/test/path?raw=query", bc.url("test/path", "raw=query").String())
	})

	t.Run("with prefix", func(t *testing.T) {
		bc, err := newBaseClient("https://test.com", true, "", "prefix/", fleet.CapabilityMap{})
		require.NoError(t, err)
		require.Equal(t, "https://test.com/prefix/test/path", bc.url("test/path", "").String())
		require.Equal(t, "https://test.com/prefix/test/path?raw=query", bc.url("test/path", "raw=query").String())
	})
}

func TestParseResponseKnownErrors(t *testing.T) {
	cases := []struct {
		message string
		code    int
		out     error
	}{
		{"not found errors", http.StatusNotFound, notFoundErr{}},
		{"unauthenticated errors", http.StatusUnauthorized, ErrUnauthenticated},
		{"license errors", http.StatusPaymentRequired, ErrMissingLicense},
	}

	for _, c := range cases {
		t.Run(c.message, func(t *testing.T) {
			bc, err := newBaseClient("https://test.com", true, "", "", fleet.CapabilityMap{})
			require.NoError(t, err)
			response := &http.Response{
				StatusCode: c.code,
				Body:       io.NopCloser(bytes.NewBufferString(`{"test": "ok"}`)),
			}
			err = bc.parseResponse("GET", "", response, &struct{}{})
			require.ErrorIs(t, err, c.out)
		})
	}
}

func TestParseResponseOK(t *testing.T) {
	bc, err := newBaseClient("https://test.com", true, "", "", fleet.CapabilityMap{})
	require.NoError(t, err)
	response := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(`{"test": "ok"}`)),
	}

	var resDest struct{ Test string }
	err = bc.parseResponse("", "", response, &resDest)
	require.NoError(t, err)
	require.Equal(t, "ok", resDest.Test)
}

func TestParseResponseGeneralErrors(t *testing.T) {
	t.Run("general HTTP errors", func(t *testing.T) {
		bc, err := newBaseClient("https://test.com", true, "", "", fleet.CapabilityMap{})
		require.NoError(t, err)
		response := &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       io.NopCloser(bytes.NewBufferString(`{"test": "ok"}`)),
		}
		err = bc.parseResponse("GET", "", response, &struct{}{})
		require.Error(t, err)
	})

	t.Run("parse errors", func(t *testing.T) {
		bc, err := newBaseClient("https://test.com", true, "", "", fleet.CapabilityMap{})
		require.NoError(t, err)
		response := &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       io.NopCloser(bytes.NewBufferString(`invalid json`)),
		}
		err = bc.parseResponse("GET", "", response, &struct{}{})
		require.Error(t, err)
	})
}

func TestNewBaseClient(t *testing.T) {
	t.Run("invalid addresses are an error", func(t *testing.T) {
		_, err := newBaseClient("invalid", true, "", "", fleet.CapabilityMap{})
		require.Error(t, err)
	})

	t.Run("http is only valid in development", func(t *testing.T) {
		_, err := newBaseClient("http://test.com", true, "", "", fleet.CapabilityMap{})
		require.Error(t, err)

		_, err = newBaseClient("http://localhost:8080", true, "", "", fleet.CapabilityMap{})
		require.NoError(t, err)

		_, err = newBaseClient("http://127.0.0.1:8080", true, "", "", fleet.CapabilityMap{})
		require.NoError(t, err)
	})
}

func TestClientCapabilities(t *testing.T) {
	cases := []struct {
		name         string
		capabilities fleet.CapabilityMap
		expected     string
	}{
		{"no capabilities", fleet.CapabilityMap{}, ""},
		{"one capability", fleet.CapabilityMap{fleet.CapabilityTokenRotation: {}}, "token_rotation"},
		{
			"multiple capabilities",
			fleet.CapabilityMap{
				fleet.CapabilityTokenRotation:       {},
				fleet.Capability("test_capability"): {},
			},
			"token_rotation,test_capability"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			bc, err := newBaseClient("https://test.com", true, "", "", c.capabilities)
			require.NoError(t, err)

			var req http.Request
			bc.setClientCapabilitiesHeader(&req)
			require.Equal(t, c.expected, req.Header.Get("X-Fleet-Capabilities"))
		})
	}
}

func TestServerCapabilities(t *testing.T) {
	// initial response has a single capability
	response := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(`{}`)),
		Header:     http.Header{"X-Fleet-Capabilities": []string{"token_rotation"}},
	}
	bc, err := newBaseClient("https://test.com", true, "", "", fleet.CapabilityMap{})
	require.NoError(t, err)

	err = bc.parseResponse("", "", response, &struct{}{})
	require.NoError(t, err)
	require.Equal(t, fleet.CapabilityMap{fleet.CapabilityTokenRotation: {}}, bc.serverCapabilities)
	require.True(t, bc.HasServerCapability(fleet.CapabilityTokenRotation))

	// later on, the server is downgraded and no longer has the capability
	response = &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(`{}`)),
		Header:     http.Header{},
	}
	err = bc.parseResponse("", "", response, &struct{}{})
	require.NoError(t, err)
	require.Equal(t, map[fleet.Capability]struct{}{}, bc.serverCapabilities)
	require.False(t, bc.HasServerCapability(fleet.CapabilityTokenRotation))

	// after an upgrade, the server has many capabilities
	response = &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(`{}`)),
		Header:     http.Header{"X-Fleet-Capabilities": []string{"token_rotation,test_capability"}},
	}
	err = bc.parseResponse("", "", response, &struct{}{})
	require.NoError(t, err)
	require.Equal(t, map[fleet.Capability]struct{}{
		fleet.CapabilityTokenRotation:       {},
		fleet.Capability("test_capability"): {},
	}, bc.serverCapabilities)
	require.True(t, bc.HasServerCapability(fleet.CapabilityTokenRotation))
	require.True(t, bc.HasServerCapability(fleet.Capability("test_capability")))
}
