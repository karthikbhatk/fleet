package maintainedapps

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/fleetdm/fleet/v4/pkg/fleethttp"
	"github.com/fleetdm/fleet/v4/server/contexts/ctxerr"
	"github.com/fleetdm/fleet/v4/server/fleet"
	kitlog "github.com/go-kit/log"
	"golang.org/x/sync/errgroup"
)

//go:embed apps.json
var appsJSON []byte

type maintainedApp struct {
	Identifier       string `json:"identifier"`
	BundleIdentifier string `json:"bundle_identifier"`
	InstallerFormat  string `json:"installer_format"`
}

const baseBrewAPIURL = "https://formulae.brew.sh/api/"

func Refresh(ctx context.Context, ds fleet.Datastore, logger kitlog.Logger) error {
	var apps []maintainedApp
	if err := json.Unmarshal(appsJSON, &apps); err != nil {
		return ctxerr.Wrap(ctx, err, "unmarshal embedded apps.json")
	}

	// allow mocking of the brew API for tests
	baseURL := baseBrewAPIURL
	if v := os.Getenv("FLEET_DEV_BREW_API_URL"); v != "" {
		baseURL = v
	}

	i := ingester{
		baseURL: baseURL,
		ds:      ds,
		logger:  logger,
	}
	return i.ingest(ctx, apps)
}

type ingester struct {
	baseURL string
	ds      fleet.Datastore
	logger  kitlog.Logger
}

func (i ingester) ingest(ctx context.Context, apps []maintainedApp) error {
	var g errgroup.Group

	client := fleethttp.NewClient(fleethttp.WithTimeout(10 * time.Second))

	// run at most 3 concurrent requests
	g.SetLimit(3)
	for _, app := range apps {
		app := app // capture loop variable, not required in Go 1.23+
		g.Go(func() error {
			return i.ingestOne(ctx, app, client)
		})
	}
	return ctxerr.Wrap(ctx, g.Wait(), "ingest apps")
}

func (i ingester) ingestOne(ctx context.Context, app maintainedApp, client *http.Client) error {
	apiURL := fmt.Sprintf("%scask/%s.json", i.baseURL, app.Identifier)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return ctxerr.Wrap(ctx, err, "create http request")
	}

	res, err := client.Do(req)
	if err != nil {
		return ctxerr.Wrap(ctx, err, "execute http request")
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return ctxerr.Wrap(ctx, err, "read http response body")
	}

	switch res.StatusCode {
	case http.StatusOK:
		// success, go on
	case http.StatusNotFound:
		// TODO: delete the existing entry? do nothing and succeed? doing the latter for now.
		return nil
	default:
		if len(body) > 512 {
			body = body[:512]
		}
		return ctxerr.Errorf(ctx, "brew API returned status %d: %s", res.StatusCode, string(body))
	}

	var cask brewCask
	if err := json.Unmarshal(body, &cask); err != nil {
		return ctxerr.Wrapf(ctx, err, "unmarshal brew cask for %s", app.Identifier)
	}
	panic("unimplemented")
}

type brewCask struct {
	Token     string          `json:"token"`
	FullToken string          `json:"full_token"`
	Tap       string          `json:"tap"`
	Name      []string        `json:"name"`
	Desc      string          `json:"desc"`
	URL       string          `json:"url"`
	Version   string          `json:"version"`
	SHA256    string          `json:"sha256"`
	Artifacts []*brewArtifact `json:"artifacts"`
}

// brew artifacts are objects that have one and only one of their fields set.
type brewArtifact struct {
	App []string `json:"app"`
	// TODO: Pkg is a bit like Binary, it is an array with a string and an object as first two elements.
	// The object has a choices field with an array of objects. See Microsoft Edge.
	Pkg       []string         `json:"pkg"`
	Uninstall []*brewUninstall `json:"uninstall"`
	Zap       []*brewZap       `json:"zap"`
	// TODO: Binary is a complex artifact - it can be provided multiple times
	// (not in an array, as items in Artifacts) and its value is an array where
	// the first element is a string - the binary artifact - and the second
	// element is an object with a "target" key). See the "docker" and "firefox" casks. Not
	// handling for now.
}

// unlike brewArtifact, a single brewUninstall can have many fields set.
// All fields can have one or multiple strings (string or []string).
type brewUninstall struct {
	LaunchCtl []string `json:"launchctl"`
	Quit      []string `json:"quit"`
	PkgUtil   []string `json:"pkgutil"`
	Script    []string `json:"script"`
	// format: [0]=signal, [1]=process name
	Signal []string `json:"signal"`
	Delete []string `json:"delete"`
	RmDir  []string `json:"rmdir"`
}

// same as brewUninstall, can be []string or string (see Microsoft Teams).
type brewZap struct {
	Trash []string `json:"trash"`
	RmDir []string `json:"rmdir"`
}
