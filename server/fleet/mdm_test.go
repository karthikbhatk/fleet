package fleet_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/fleetdm/fleet/v4/server/fleet"
	apple_mdm "github.com/fleetdm/fleet/v4/server/mdm/apple"
	"github.com/fleetdm/fleet/v4/server/mock"
	"github.com/go-kit/log"
	nanodep_client "github.com/micromdm/nanodep/client"
	"github.com/micromdm/nanodep/godep"
	"github.com/stretchr/testify/require"
)

type mockStorage struct {
	token string
	url   string
}

func (s mockStorage) RetrieveAuthTokens(ctx context.Context, name string) (*nanodep_client.OAuth1Tokens, error) {
	return &nanodep_client.OAuth1Tokens{AccessToken: s.token}, nil
}

func (s mockStorage) RetrieveConfig(context.Context, string) (*nanodep_client.Config, error) {
	return &nanodep_client.Config{BaseURL: s.url}, nil
}

func TestDEPClient(t *testing.T) {
	ctx := context.Background()

	rxToken := regexp.MustCompile(`oauth_token="(\w+)"`)
	const (
		validToken                 = "OK"
		invalidToken               = "FAIL"
		termsChangedToken          = "TERMS"
		termsChangedAfterAuthToken = "TERMS_AFTER"
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/session" {
			matches := rxToken.FindStringSubmatch(r.Header.Get("Authorization"))
			require.NotNil(t, matches)
			token := matches[1]

			switch token {
			case validToken:
				_, _ = w.Write([]byte(`{"auth_session_token": "ok"}`))
			case termsChangedAfterAuthToken:
				_, _ = w.Write([]byte(`{"auth_session_token": "fail"}`))
			case termsChangedToken:
				w.WriteHeader(http.StatusForbidden)
				_, _ = w.Write([]byte(`{"code": "T_C_NOT_SIGNED"}`))
			case invalidToken:
				w.WriteHeader(http.StatusForbidden)
				_, _ = w.Write([]byte(`{"code": "ACCESS_DENIED"}`))
			default:
				w.WriteHeader(http.StatusUnauthorized)
			}
			return
		}

		require.Equal(t, "/account", r.URL.Path)
		authSsn := r.Header.Get("X-Adm-Auth-Session")
		if authSsn == "fail" {
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(`{"code": "T_C_NOT_SIGNED"}`))
			return
		}

		// otherwise, return account information, details not important for this
		// test.
		_, _ = w.Write([]byte(`{"admin_id": "test"}`))
	}))
	defer srv.Close()

	logger := log.NewNopLogger()
	ds := new(mock.Store)

	appCfg := fleet.AppConfig{}
	ds.AppConfigFunc = func(ctx context.Context) (*fleet.AppConfig, error) {
		return &appCfg, nil
	}
	ds.SaveAppConfigFunc = func(ctx context.Context, config *fleet.AppConfig) error {
		appCfg = *config
		return nil
	}

	checkDSCalled := func(readInvoked, writeInvoked bool) {
		require.Equal(t, readInvoked, ds.AppConfigFuncInvoked)
		require.Equal(t, writeInvoked, ds.SaveAppConfigFuncInvoked)
		ds.AppConfigFuncInvoked = false
		ds.SaveAppConfigFuncInvoked = false
	}

	cases := []struct {
		token        string
		wantErr      bool
		readInvoked  bool
		writeInvoked bool
		termsFlag    bool
	}{
		// use a valid token, appconfig should not be updated (already unflagged)
		{token: validToken, wantErr: false, readInvoked: true, writeInvoked: false, termsFlag: false},

		// use an invalid token, appconfig should not even be read (not a terms error)
		{token: invalidToken, wantErr: true, readInvoked: false, writeInvoked: false, termsFlag: false},

		// terms changed during the auth request
		{token: termsChangedToken, wantErr: true, readInvoked: true, writeInvoked: true, termsFlag: true},

		// use of an invalid token does not update the flag
		{token: invalidToken, wantErr: true, readInvoked: false, writeInvoked: false, termsFlag: true},

		// use of a valid token resets the flag
		{token: validToken, wantErr: false, readInvoked: true, writeInvoked: true, termsFlag: false},

		// use of a valid token again does not update the appConfig
		{token: validToken, wantErr: false, readInvoked: true, writeInvoked: false, termsFlag: false},

		// terms changed during the actual account request, after auth
		{token: termsChangedAfterAuthToken, wantErr: true, readInvoked: true, writeInvoked: true, termsFlag: true},

		// again terms changed after auth, doesn't update appConfig
		{token: termsChangedAfterAuthToken, wantErr: true, readInvoked: true, writeInvoked: false, termsFlag: true},

		// terms changed during auth, doesn't update appConfig
		{token: termsChangedToken, wantErr: true, readInvoked: true, writeInvoked: false, termsFlag: true},

		// valid token, resets the flag
		{token: validToken, wantErr: false, readInvoked: true, writeInvoked: true, termsFlag: false},
	}

	// order of calls is important, and test must not be parallelized as it would
	// be racy. For that reason, subtests are not used (it would make it possible
	// to run one subtest in isolation, which could fail).
	for i, c := range cases {
		t.Logf("case %d", i)

		store := mockStorage{token: c.token, url: srv.URL}
		dep := fleet.NewDEPClient(store, ds, logger)
		res, err := dep.AccountDetail(ctx, apple_mdm.DEPName)

		if c.wantErr {
			var httpErr *godep.HTTPError
			require.Error(t, err)
			if errors.As(err, &httpErr) {
				require.Equal(t, http.StatusForbidden, httpErr.StatusCode)
			} else {
				var authErr *nanodep_client.AuthError
				require.ErrorAs(t, err, &authErr)
				require.Equal(t, http.StatusForbidden, authErr.StatusCode)
			}
			if c.token == termsChangedToken || c.token == termsChangedAfterAuthToken {
				require.True(t, godep.IsTermsNotSigned(err))
			} else {
				require.False(t, godep.IsTermsNotSigned(err))
			}
		} else {
			require.NoError(t, err)
			require.Equal(t, "test", res.AdminID)
		}
		checkDSCalled(c.readInvoked, c.writeInvoked)
		require.Equal(t, c.termsFlag, appCfg.MDM.AppleBMTermsExpired)
	}
}
