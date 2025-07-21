package auth

import (
	"context"
	"main/tools/crypto"
)

type IAuthorization interface {
	Authorize(ctx context.Context, token string) (updatedCtx context.Context, success bool)
}

func CreateAuthorization() IAuthorization {
	return &Authorization{
		tokenManager: crypto.CreateTokenManager(),
	}
}

type Authorization struct {
	tokenManager crypto.ITokenManager
}

func (auth *Authorization) Authorize(ctx context.Context, token string) (updatedCtx context.Context, success bool) {
	updatedCtx = ctx

	user, err := auth.tokenManager.ValidateToken(token)
	if err != nil {
		return
	}

	ctx = context.WithValue(ctx, "user", user)
	success = true
	return
}
