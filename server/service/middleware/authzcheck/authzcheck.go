// Package authzcheck implements a middleware that ensures that an authorization
// check was performed. This does not ensure that the correct authorization
// check was performed, but offers a backstop in case that a developer misses a
// check.
package authzcheck

import (
	"context"
	"reflect"
	"runtime"

	"github.com/fleetdm/fleet/server/authz"
	authz_ctx "github.com/fleetdm/fleet/server/contexts/authz"
	"github.com/go-kit/kit/endpoint"
)

// Middleware is the authzcheck middleware type.
type Middleware struct{}

// NewMiddleware returns a new authzcheck middleware.
func NewMiddleware() *Middleware {
	return &Middleware{}
}

func (m *Middleware) AuthzCheck() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			authzctx := &authz_ctx.AuthorizationContext{}
			ctx = authz_ctx.NewContext(ctx, authzctx)

			response, error := next(ctx, req)

			// If authorization was not checked, return a response that will
			// marshal to a generic error and log that the check was missed.
			if !authzctx.Checked {
				funcName := runtime.FuncForPC(reflect.ValueOf(next).Pointer()).Name()
				return nil, authz.ForbiddenWithInternal("missed authz check: " + funcName)
			}

			return response, error
		}
	}
}
