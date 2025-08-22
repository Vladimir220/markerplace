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
}

func CreateHandlers(tokenManager crypto.ITokenManager, dao postgres.IMarketplaceDAO) Handlers {
	return Handlers{
		logger:       log.CreateLogger("Handlers"),
		dao:          dao,
		tokenManager: tokenManager,
	}
}
