package crypto

import (
	"context"
	"fmt"
	"marketplace/db/DAO"
	"marketplace/models"

	"github.com/Vladimir220/markerplace/logger_lib"
)

func CreateTokenManagerGateway(ctx context.Context, reserveTokensDAO DAO.ITokensDAO) (tokenManager ITokenManager) {
	logger := logger_lib.CreateLoggerGateway(ctx, "TokenManagerGateway")

	grpcTokenManager, err := CreateGrpcTokenManager(ctx)
	var remoteUnavailable bool
	if err != nil {
		logger.WriteWarning(fmt.Sprintf("%s: %v", "CreateTokenManagerGateway(): remoteAuth unavailable", err))
		remoteUnavailable = true
	}

	reserveTokenManager := CreateTokenManager(ctx, reserveTokensDAO)

	tokenManager = &TokenManagerGateway{
		remoteUnavailable:   remoteUnavailable,
		remoteTokenManager:  grpcTokenManager,
		reserveTokenManager: reserveTokenManager,
		ctx:                 ctx,
		logger:              logger,
	}

	return
}

// Proxy for ITokenManager
type TokenManagerGateway struct {
	remoteUnavailable   bool
	reserveTokenManager ITokenManager
	remoteTokenManager  ITokenManager
	ctx                 context.Context
	logger              logger_lib.ILogger
}

func (tm TokenManagerGateway) GenerateToken(user models.User) (token string, isErr bool) {
	var err error
	if !tm.remoteUnavailable {
		token, isErr = tm.remoteTokenManager.GenerateToken(user)
		if !isErr {
			return
		}
	}

	tm.logger.WriteWarning(fmt.Sprintf("%s: %v", "GenerateToken(): remote TokenManager unavailable", err))
	token, isErr = tm.reserveTokenManager.GenerateToken(user)

	return
}

func (tm TokenManagerGateway) ValidateToken(token string) (user models.User, isValid, isErr bool) {
	var err error
	if !tm.remoteUnavailable {
		user, isValid, isErr = tm.remoteTokenManager.ValidateToken(token)
		if !isErr {
			return
		}
	}

	tm.logger.WriteWarning(fmt.Sprintf("%s: %v", "ValidateToken(): remote TokenManager unavailable", err))
	user, isValid, isErr = tm.reserveTokenManager.ValidateToken(token)

	return
}
