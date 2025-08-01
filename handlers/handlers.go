package handlers

import (
	"main/db/DAO/postgres"
	"main/tools/crypto"
	"main/tools/log"
)

type Handlers struct {
	logger       log.ILogger
	dao          postgres.IMarcketplaceDAO
	tokenManager crypto.ITokenManager
}

func CreateHandlers(tokenManager crypto.ITokenManager) Handlers {
	return Handlers{
		logger:       log.CreateLogger("Handlers"),
		dao:          postgres.CreateMarcketplaceDAO(),
		tokenManager: tokenManager,
	}
}

func (h Handlers) Close() {
	h.dao.Close()
}
