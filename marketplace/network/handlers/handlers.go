package handlers

import (
	"main/crypto"
	"main/db/DAO/postgres"
	"main/log"
)

type Handlers struct {
	logger       log.ILogger
	dao          postgres.IMarketplaceDAO
	tokenManager crypto.ITokenManager
	infoLogs     bool
}

func CreateHandlers(tokenManager crypto.ITokenManager, dao postgres.IMarketplaceDAO, infoLogs bool) Handlers {
	return Handlers{
		logger:       log.CreateLogger("Handlers"),
		dao:          dao,
		tokenManager: tokenManager,
		infoLogs:     infoLogs,
	}
}
