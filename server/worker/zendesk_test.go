package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fleetdm/fleet/v4/server/fleet"
	"github.com/fleetdm/fleet/v4/server/mock"
	"github.com/fleetdm/fleet/v4/server/service/externalsvc"
	kitlog "github.com/go-kit/kit/log"
	"github.com/stretchr/testify/require"
)

func TestZendeskRun(t *testing.T) {
	ds := new(mock.Store)
	ds.HostsByCVEFunc = func(ctx context.Context, cve string) ([]*fleet.HostShort, error) {
		return []*fleet.HostShort{
			{
				ID:       1,
				Hostname: "test",
			},
		}, nil
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(501)
			return
		}
		if !strings.Contains(r.URL.Path, "/api/v2/tickets") {
			w.WriteHeader(502)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)
		require.Contains(t, string(body), `"subject":"Vulnerability CVE-1234-5678 detected on 1 host(s)"`)

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	client, err := externalsvc.NewZendeskTestClient(&externalsvc.ZendeskOptions{URL: fmt.Sprintf("%s/api/v2/tickets", srv.URL)})
	require.NoError(t, err)

	zendesk := &Zendesk{
		FleetURL:      "https://fleetdm.com",
		Datastore:     ds,
		Log:           kitlog.NewNopLogger(),
		ZendeskClient: client,
	}

	err = zendesk.Run(context.Background(), json.RawMessage(`{"cve":"CVE-1234-5678"}`))
	require.NoError(t, err)
}

func TestZendeskQueueJobs(t *testing.T) {
	ds := new(mock.Store)
	ctx := context.Background()
	logger := kitlog.NewNopLogger()

	t.Run("success", func(t *testing.T) {
		ds.NewJobFunc = func(ctx context.Context, job *fleet.Job) (*fleet.Job, error) {
			return job, nil
		}
		err := QueueZendeskJobs(ctx, ds, logger, map[string][]string{"CVE-1234-5678": nil})
		require.NoError(t, err)
		require.True(t, ds.NewJobFuncInvoked)
	})

	t.Run("failure", func(t *testing.T) {
		ds.NewJobFunc = func(ctx context.Context, job *fleet.Job) (*fleet.Job, error) {
			return nil, io.EOF
		}
		err := QueueZendeskJobs(ctx, ds, logger, map[string][]string{"CVE-1234-5678": nil})
		require.Error(t, err)
		require.ErrorIs(t, err, io.EOF)
		require.True(t, ds.NewJobFuncInvoked)
	})
}
