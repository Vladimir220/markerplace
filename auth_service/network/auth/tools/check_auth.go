package tools

import (
	"auth_service/models"
	"context"
)

func CheckAuth(ctx context.Context) (user models.User, ok bool) {
	user, ok = ctx.Value("user").(models.User)
	return
}
