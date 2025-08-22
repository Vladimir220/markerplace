package auth

import (
	"context"
	"fmt"
	"main/crypto"
	"main/log"
)

type IAuthorization interface {
	Authorize(ctx context.Context, token string) (updatedCtx context.Context, success bool)
}

func CreateAuthorization(tokenManager crypto.ITokenManager) IAuthorization {
	return &Authorization{
		tokenManager: tokenManager,
		logger:       log.CreateLogger("Authorization"),
	}
}

type Authorization struct {
	tokenManager crypto.ITokenManager
	logger       log.ILogger
}

func (auth *Authorization) Authorize(ctx context.Context, token string) (updatedCtx context.Context, success bool) {
	logLabel := "Authorize():"

	user, success, err := auth.tokenManager.ValidateToken(token)
	if !success {
		if err != nil {
			auth.logger.WriteError(fmt.Sprintf("%s %v", logLabel, err))
			return
		}
		return
	}
	updatedCtx = context.WithValue(ctx, "user", user)
	success = true
	//auth.logger.WriteInfo(fmt.Sprintf("Authorization:Authorize(): user \"%s\" authorize", user.Login))
	return
}
