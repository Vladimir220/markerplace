package tools

import (
	"context"
	"marketplace/models"
)

func CheckAuth(ctx context.Context) (user models.User, ok bool) {
	user, ok = ctx.Value("user").(models.User)
	return
}
