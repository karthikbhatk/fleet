package service

import (
	"fmt"

	"github.com/WatchBeam/clock"
	"github.com/fleetdm/fleet/v4/server/authz"
	"github.com/fleetdm/fleet/v4/server/config"
	"github.com/fleetdm/fleet/v4/server/fleet"
	kitlog "github.com/go-kit/kit/log"
)

type Service struct {
	fleet.Service

	ds      fleet.Datastore
	logger  kitlog.Logger
	config  config.FleetConfig
	clock   clock.Clock
	authz   *authz.Authorizer
	license *fleet.LicenseInfo
}

func NewService(
	svc fleet.Service,
	ds fleet.Datastore,
	logger kitlog.Logger,
	config config.FleetConfig,
	mailService fleet.MailService,
	c clock.Clock,
	license *fleet.LicenseInfo,
) (*Service, error) {

	authorizer, err := authz.NewAuthorizer()
	if err != nil {
		return nil, fmt.Errorf("new authorizer: %w", err)
	}

	return &Service{
		Service: svc,
		ds:      ds,
		logger:  logger,
		config:  config,
		clock:   c,
		authz:   authorizer,
		license: license,
	}, nil
}

// TODO(mna): copied from server/service/transport_error.go for now, should
// eventually have common implementations of HTTP-related errors. #4406
type badRequestError struct {
	message string
}

func (e *badRequestError) Error() string {
	return e.message
}

func (e *badRequestError) BadRequestError() []map[string]string {
	return nil
}
