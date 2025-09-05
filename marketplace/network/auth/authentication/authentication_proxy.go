package authentication

import (
	"context"
	"fmt"
	"marketplace/crypto"
	"marketplace/network/auth/tools"

	"github.com/Vladimir220/markerplace/logger_lib"
)

func CreateAuthenticationProxy(ctx context.Context, localTokenManager crypto.ITokenManager, infoLogs bool) (auth IAuthentication, err error) {
	logger := logger_lib.CreateLoggerAdapter(ctx, "AuthenticationProxy")

	grpcAuth, err := CreateGrpcAuthentication(ctx)
	var remoteUnavailable bool
	if err != nil {
		logger.WriteWarning(fmt.Sprintf("%s: %v", "CreateAuthenticationProxy(): remoteAuth unavailable", err))
		remoteUnavailable = true
	}

	localAuth, err := CreateAuthentication(ctx, localTokenManager, infoLogs)
	if err != nil {
		err = fmt.Errorf("CreateAuthenticationProxy(): %v", err)
		return
	}

	auth = &AuthenticationProxy{
		remoteUnavailable: remoteUnavailable,
		remoteAuth:        grpcAuth,
		localAuth:         localAuth,
		ctx:               ctx,
		logger:            logger,
	}

	return
}

type AuthenticationProxy struct {
	remoteUnavailable bool
	localAuth         IAuthentication
	remoteAuth        IAuthentication
	ctx               context.Context
	logger            logger_lib.ILogger
}

func (auth AuthenticationProxy) Register(login, password string) (token string, err error) {
	if !auth.remoteUnavailable {
		token, err = auth.remoteAuth.Register(login, password)
		if err == tools.ErrLoginFormat || err == tools.ErrPasswordFormat || err == tools.ErrLoginIsTaken || err == nil {
			return
		}
	}

	auth.logger.WriteWarning(fmt.Sprintf("%s: %v", "Register(): remoteAuth unavailable", err))
	token, err = auth.localAuth.Register(login, password)

	return
}

func (auth AuthenticationProxy) Login(login, password string) (token string, err error) {
	if !auth.remoteUnavailable {
		token, err = auth.remoteAuth.Login(login, password)
		if err == tools.ErrLogin || err == tools.ErrLoginFormat || err == nil {
			return
		}
	}

	auth.logger.WriteWarning(fmt.Sprintf("%s: %v", "Login(): remoteAuth unavailable", err))
	token, err = auth.localAuth.Login(login, password)
	return
}
