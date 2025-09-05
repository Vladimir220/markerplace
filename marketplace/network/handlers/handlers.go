package handlers

import (
	"context"
	"marketplace/crypto"
	"marketplace/db/DAO/postgres"
	"marketplace/log"
)

type Handlers struct {
	logger       log.ILogger
	dao          postgres.IMarketplaceDAO
	tokenManager crypto.ITokenManager
	infoLogs     bool
	ctx          context.Context
}

func CreateHandlers(ctx context.Context, tokenManager crypto.ITokenManager, dao postgres.IMarketplaceDAO, infoLogs bool) Handlers {
	return Handlers{
		logger:       log.CreateLoggerAdapter(ctx, "Handlers"),
		dao:          dao,
		tokenManager: tokenManager,
		infoLogs:     infoLogs,
		ctx:          ctx,
	}
}
