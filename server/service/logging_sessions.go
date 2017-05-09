package service

import (
	"context"
	"time"

	"github.com/kolide/kolide/server/kolide"
)

func (mw loggingMiddleware) Login(ctx context.Context, username, password string) (user *kolide.User, token string, err error) {

	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "Login",
			"user", username,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	user, token, err = mw.Service.Login(ctx, username, password)
	return
}

func (mw loggingMiddleware) Logout(ctx context.Context) (err error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "Logout",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.Service.Logout(ctx)
	return
}

func (mw loggingMiddleware) InitiateSSO(ctx context.Context, relayURL string) (idpURL string, err error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "InitiateSSO",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	idpURL, err = mw.Service.InitiateSSO(ctx, relayURL)
	return
}

func (mw loggingMiddleware) CallbackSSO(ctx context.Context, auth kolide.Auth) (sess *kolide.SSOSession, err error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "CallbackSSO",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	sess, err = mw.Service.CallbackSSO(ctx, auth)
	return
}
