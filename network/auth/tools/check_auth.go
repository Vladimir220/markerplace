package tools

import (
	"context"
	"main/models"
)

func CheckAuth(ctx context.Context) (user models.User, ok bool) {
	user, ok = ctx.Value("user").(models.User)
	return
}
