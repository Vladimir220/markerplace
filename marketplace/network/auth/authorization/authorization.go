package authorization

import (
	"context"
	"fmt"
	"marketplace/crypto"
	"marketplace/env"

	"github.com/Vladimir220/markerplace/logger_lib"
)

type IAuthorization interface {
	Authorize(ctx context.Context, token string) (updatedCtx context.Context, success bool)
}

func CreateAuthorization(ctx context.Context, tokenManager crypto.ITokenManager) IAuthorization {
	return &Authorization{
		tokenManager: tokenManager,
		logger:       logger_lib.CreateLoggerGateway(ctx, "Authorization"),
		infoLogs:     env.GetLogsConfig().PrintAuthorizationInfo,
	}
}

type Authorization struct {
	tokenManager crypto.ITokenManager
	logger       logger_lib.ILogger
	infoLogs     bool
}

func (auth *Authorization) Authorize(ctx context.Context, token string) (updatedCtx context.Context, success bool) {
	logLabel := "Authorize():"

	user, isValid, isErr := auth.tokenManager.ValidateToken(token)
	if !isValid || isErr {
		return
	}

	updatedCtx = context.WithValue(ctx, "user", user)
	success = true

	if auth.infoLogs {
		auth.logger.WriteInfo(fmt.Sprintf("%s user \"%s\" authorize", logLabel, user.Login))
	}
	return
}
