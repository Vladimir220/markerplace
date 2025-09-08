package authentication

import (
	"context"
	"fmt"
	"marketplace/crypto"
	"marketplace/db/DAO/postgres"
	"marketplace/network/auth/tools"

	"github.com/Vladimir220/markerplace/logger_lib"
)

func CreateAuthenticationGateway(ctx context.Context, reserveTokenManager crypto.ITokenManager, reserveDao postgres.IMarketplaceDAO) (auth IAuthentication) {
	logger := logger_lib.CreateLoggerGateway(ctx, "AuthenticationGateway")

	grpcAuth, err := CreateGrpcAuthentication(ctx)
	var remoteUnavailable bool
	if err != nil {
		logger.WriteWarning(fmt.Sprintf("%s: %v", "CreateAuthenticationGateway(): remoteAuth unavailable", err))
		remoteUnavailable = true
	}

	reserveAuth := CreateAuthentication(ctx, reserveTokenManager, reserveDao)

	auth = &AuthenticationGateway{
		remoteUnavailable: remoteUnavailable,
		remoteAuth:        grpcAuth,
		reserveAuth:       reserveAuth,
		ctx:               ctx,
		logger:            logger,
	}

	return
}

// Proxy for IAuthentication
type AuthenticationGateway struct {
	remoteUnavailable bool
	reserveAuth       IAuthentication
	remoteAuth        IAuthentication
	ctx               context.Context
	logger            logger_lib.ILogger
}

func (auth AuthenticationGateway) Register(login, password string) (token string, err error) {
	if !auth.remoteUnavailable {
		token, err = auth.remoteAuth.Register(login, password)
		if err == tools.ErrLoginFormat || err == tools.ErrPasswordFormat || err == tools.ErrLoginIsTaken || err == nil {
			return
		}
	}

	auth.logger.WriteWarning(fmt.Sprintf("%s: %v", "Register(): remoteAuth unavailable", err))
	token, err = auth.reserveAuth.Register(login, password)

	return
}

func (auth AuthenticationGateway) Login(login, password string) (token string, err error) {
	if !auth.remoteUnavailable {
		token, err = auth.remoteAuth.Login(login, password)
		if err == tools.ErrLogin || err == tools.ErrLoginFormat || err == nil {
			return
		}
	}

	auth.logger.WriteWarning(fmt.Sprintf("%s: %v", "Login(): remoteAuth unavailable", err))
	token, err = auth.reserveAuth.Login(login, password)
	return
}
