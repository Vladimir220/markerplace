package auth

import (
	"context"
	"main/tools/crypto"
)

type IAuthorization interface {
	Authorize(ctx context.Context, token string) (updatedCtx context.Context, success bool)
}

func CreateAuthorization(tokenManager crypto.ITokenManager) IAuthorization {
	return &Authorization{
		tokenManager: tokenManager,
	}
}

type Authorization struct {
	tokenManager crypto.ITokenManager
}

func (auth *Authorization) Authorize(ctx context.Context, token string) (updatedCtx context.Context, success bool) {
	user, err := auth.tokenManager.ValidateToken(token)
	if err != nil {
		return
	}

	updatedCtx = context.WithValue(ctx, "user", user)
	success = true
	return
}
