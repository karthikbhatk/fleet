package authzcheck

import (
	"context"
	"testing"

	"github.com/fleetdm/fleet/server/contexts/authz"
	"github.com/fleetdm/fleet/server/fleet"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthzCheck(t *testing.T) {
	t.Parallel()

	checker := NewMiddleware()

	check := func(ctx context.Context, req interface{}) (interface{}, error) {
		authCtx, ok := authz.FromContext(ctx)
		require.True(t, ok)
		authCtx.Checked = true
		return struct{}{}, nil
	}
	check = checker.AuthzCheck()(check)

	_, err := check(context.Background(), struct{}{})
	assert.NoError(t, err)
}

func TestAuthzCheckAuthFailed(t *testing.T) {
	t.Parallel()

	checker := NewMiddleware()

	check := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, fleet.NewAuthFailedError("failed")
	}
	check = checker.AuthzCheck()(check)

	_, err := check(context.Background(), struct{}{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed")
}

func TestAuthzCheckAuthRequired(t *testing.T) {
	t.Parallel()

	checker := NewMiddleware()

	check := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, fleet.NewAuthRequiredError("required")
	}
	check = checker.AuthzCheck()(check)

	_, err := check(context.Background(), struct{}{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "required")
}

func TestAuthzCheckMissing(t *testing.T) {
	t.Parallel()

	checker := NewMiddleware()

	nocheck := func(context.Context, interface{}) (interface{}, error) { return struct{}{}, nil }
	nocheck = checker.AuthzCheck()(nocheck)

	_, err := nocheck(context.Background(), struct{}{})
	assert.Error(t, err)
}
