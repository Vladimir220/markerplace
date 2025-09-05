package crypto

import (
	"context"
	"fmt"
	"marketplace/db/DAO"
	"marketplace/log"
	"marketplace/models"
)

func CreateTokenManagerProxy(ctx context.Context, tokensDAO DAO.ITokensDAO, infoLogs bool) (tokenManager ITokenManager) {
	logger := log.CreateLoggerAdapter(ctx, "TokenManagerProxy")

	grpcTokenManager, err := CreateGrpcTokenManager(ctx)
	var remoteUnavailable bool
	if err != nil {
		logger.WriteWarning(fmt.Sprintf("%s: %v", "CreateTokenManagerProxy(): remoteAuth unavailable", err))
		remoteUnavailable = true
	}

	localTokenManager := CreateTokenManager(ctx, tokensDAO, infoLogs)

	tokenManager = &TokenManagerProxy{
		remoteUnavailable:  remoteUnavailable,
		remoteTokenManager: grpcTokenManager,
		localTokenManager:  localTokenManager,
		ctx:                ctx,
		logger:             logger,
	}

	return
}

type TokenManagerProxy struct {
	remoteUnavailable  bool
	localTokenManager  ITokenManager
	remoteTokenManager ITokenManager
	ctx                context.Context
	logger             log.ILogger
}

func (tm TokenManagerProxy) GenerateToken(user models.User) (token string, isErr bool) {
	var err error
	if !tm.remoteUnavailable {
		token, isErr = tm.remoteTokenManager.GenerateToken(user)
		if !isErr {
			return
		}
	}

	tm.logger.WriteWarning(fmt.Sprintf("%s: %v", "GenerateToken(): remote TokenManager unavailable", err))
	token, isErr = tm.localTokenManager.GenerateToken(user)

	return
}

func (tm TokenManagerProxy) ValidateToken(token string) (user models.User, isValid, isErr bool) {
	var err error
	if !tm.remoteUnavailable {
		user, isValid, isErr = tm.remoteTokenManager.ValidateToken(token)
		if !isErr {
			return
		}
	}

	tm.logger.WriteWarning(fmt.Sprintf("%s: %v", "ValidateToken(): remote TokenManager unavailable", err))
	user, isValid, isErr = tm.localTokenManager.ValidateToken(token)

	return
}
