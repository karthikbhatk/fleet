package vpp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/fleetdm/fleet/v4/pkg/fleethttp"
)

// Asset is a product in the store.
//
// https://developer.apple.com/documentation/devicemanagement/asset
type Asset struct {
	// AdamID is the unique identifier for a product in the store.
	AdamID string `json:"adamId"`
	// PricingParam is the quality of a product in the store.
	// Possible Values are `STDQ` and `PLUS`
	PricingParam string `json:"pricingParam"`
}

// ErrorResponse represents the response that contains the error that occurs.
//
// https://developer.apple.com/documentation/devicemanagement/errorresponse
type ErrorResponse struct {
	ErrorInfo    ResponseErrorInfo `json:"errorInfo"`
	ErrorMessage string            `json:"errorMessage"`
	ErrorNumber  int32             `json:"errorNumber"`
}

// Error implements the Erorrer interface
func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("Apple VPP endpoint returned error: %s (error number: %d)", e.ErrorMessage, e.ErrorNumber)
}

// ResponseErrorInfo represents the request-specific information regarding the
// failure.
//
// https://developer.apple.com/documentation/devicemanagement/responseerrorinfo
type ResponseErrorInfo struct {
	Assets        []Asset  `json:"assets"`
	ClientUserIds []string `json:"clientUserIds"`
	SerialNumbers []string `json:"serialNumbers"`
}

// client is a package-level client (similar to http.DefaultClient) so it can
// be reused instead of crated as needed, as the internal Transport typically
// has internal state (cached connections, etc) and it's safe for concurrent
// use.
var client = fleethttp.NewClient(fleethttp.WithTimeout(10 * time.Second))

// GetConfig fetches the VPP config from Apple's VPP API. This doubles as a
// verification that the user-provided VPP token is valid.
//
// https://developer.apple.com/documentation/devicemanagement/client_config-a40
func GetConfig(token string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, getBaseURL()+"/client/config", nil)
	if err != nil {
		return "", fmt.Errorf("creating request to Apple VPP endpoint: %w", err)
	}

	var respJSON struct {
		LocationName string `json:"locationName"`
	}

	if err := do(req, token, &respJSON); err != nil {
		return "", fmt.Errorf("making request to Apple VPP endpoint: %w", err)
	}

	return respJSON.LocationName, nil
}

// AssociateAssetsRequest is the request for asset management.
type AssociateAssetsRequest struct {
	// Assets are the assets to assign.
	Assets []Asset `json:"assets"`
	// SerialNumbers is the set of identifiers for devices to assign the
	// assets to.
	SerialNumbers []string `json:"serialNumbers"`
}

// AssociateAssets associates assets to serial numbers according the the
// request parameters provided.
//
// https://developer.apple.com/documentation/devicemanagement/associate_assets
func AssociateAssets(token string, params *AssociateAssetsRequest) error {
	var reqBody bytes.Buffer
	if err := json.NewEncoder(&reqBody).Encode(params); err != nil {
		return fmt.Errorf("encoding params as JSON: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, getBaseURL()+"/assets/associate", &reqBody)
	if err != nil {
		return fmt.Errorf("creating request to Apple VPP endpoint: %w", err)
	}

	if err := do[any](req, token, nil); err != nil {
		return fmt.Errorf("making request to Apple VPP endpoint: %w", err)
	}
	return nil
}

// AssetFilter represents the filters for querying assets.
type AssetFilter struct {
	// PageIndex is the requested page index.
	PageIndex int32 `json:"pageIndex"`

	// ProductType is the filter for the asset product type.
	// Possible Values: App, Book
	ProductType string `json:"productType"`

	// PricingParam is the filter for the asset product quality.
	// Possible Values: STDQ, PLUS
	PricingParam string `json:"pricingParam"`

	// Revocable is the filter for asset revocability.
	Revocable *bool `json:"revocable"`

	// DeviceAssignable is the filter for asset device assignability.
	DeviceAssignable *bool `json:"deviceAssignable"`

	// MaxAvailableCount is the filter for the maximum inclusive assets available count.
	MaxAvailableCount int32 `json:"maxAvailableCount"`

	// MinAvailableCount is the filter for the minimum inclusive assets available count.
	MinAvailableCount int32 `json:"minAvailableCount"`

	// MaxAssignedCount is the filter for the maximum inclusive assets assigned count.
	MaxAssignedCount int32 `json:"maxAssignedCount"`

	// MinAssignedCount is the filter for the minimum inclusive assets assigned count.
	MinAssignedCount int32 `json:"minAssignedCount"`

	// AdamID is the filter for the asset product unique identifier.
	AdamID string `json:"adamId"`
}

// GetAssets fetches the assets from Apple's VPP API with optional filters.
func GetAssets(token string, filter *AssetFilter) ([]Asset, error) {
	baseURL := getBaseURL() + "/assets"
	reqURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("parsing base URL: %w", err)
	}

	if filter != nil {
		query := url.Values{}
		addFilter(query, "adamId", filter.AdamID)
		addFilter(query, "pricingParam", filter.PricingParam)
		addFilter(query, "productType", filter.ProductType)
		addFilter(query, "revocable", filter.Revocable)
		addFilter(query, "deviceAssignable", filter.DeviceAssignable)
		addFilter(query, "maxAvailableCount", filter.MaxAvailableCount)
		addFilter(query, "minAvailableCount", filter.MinAvailableCount)
		addFilter(query, "maxAssignedCount", filter.MaxAssignedCount)
		addFilter(query, "minAssignedCount", filter.MinAssignedCount)
		addFilter(query, "pageIndex", filter.PageIndex)
		reqURL.RawQuery = query.Encode()
	}

	req, err := http.NewRequest(http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("creating request to Apple VPP endpoint: %w", err)
	}

	var bodyResp struct {
		Assets []Asset `json:"assets"`
	}

	if err = do(req, token, &bodyResp); err != nil {
		return nil, fmt.Errorf("retrieving assets: %w", err)
	}

	return bodyResp.Assets, nil
}

func do[T any](req *http.Request, token string, dest *T) error {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("making request to Apple VPP endpoint: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body from Apple VPP endpoint: %w", err)
	}

	// For some reason, Apple returns 200 OK even if you pass an invalid token in the Auth header.
	// We will need to parse the response and check to see if it contains an error.
	var errResp ErrorResponse
	if err := json.Unmarshal(body, &errResp); err == nil && (errResp.ErrorMessage != "" || errResp.ErrorNumber != 0) {
		return &errResp
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("calling Apple VPP endpoint failed with status %d", resp.StatusCode)
	}

	if dest != nil {
		if err := json.Unmarshal(body, dest); err != nil {
			return fmt.Errorf("decoding response data from Apple VPP endpoint: %w", err)
		}
	}

	return nil
}

func getBaseURL() string {
	devURL := os.Getenv("FLEET_DEV_VPP_URL")
	if devURL != "" {
		return devURL
	}
	return "https://vpp.itunes.apple.com/mdm/v2"
}

// addFilter adds a filter to the query values if it is not the zero value.
func addFilter(query url.Values, key string, value any) {
	switch v := value.(type) {
	case string:
		if v != "" {
			query.Add(key, v)
		}
	case *bool:
		if v != nil {
			query.Add(key, strconv.FormatBool(*v))
		}
	case int32:
		if v != 0 {
			query.Add(key, fmt.Sprintf("%d", v))
		}
	}
}
