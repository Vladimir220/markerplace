package handlers

import (
	"context"
	"marketplace/crypto"
	"marketplace/db/DAO/postgres"
	"marketplace/env"

	"github.com/Vladimir220/markerplace/logger_lib"
)

type Handlers struct {
	logger       logger_lib.ILogger
	dao          postgres.IMarketplaceDAO
	tokenManager crypto.ITokenManager
	infoLogs     bool
	ctx          context.Context
}

func CreateHandlers(ctx context.Context, tokenManager crypto.ITokenManager, dao postgres.IMarketplaceDAO) Handlers {
	return Handlers{
		logger:       logger_lib.CreateLoggerGateway(ctx, "Handlers"),
		dao:          dao,
		tokenManager: tokenManager,
		infoLogs:     env.GetLogsConfig().PrintHandlersInfo,
		ctx:          ctx,
	}
}
