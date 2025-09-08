package proxies

import (
	"auth_service/db/DAO/postgres"
	"auth_service/env"
	"auth_service/models"
	"context"
	"fmt"

	"github.com/Vladimir220/markerplace/logger_lib"
)

func CreateDAOWithLog(ctx context.Context, dao postgres.IAuthMarketplaceDAO) postgres.IAuthMarketplaceDAO {
	return &DAOWithLog{
		original: dao,
		logger:   logger_lib.CreateLoggerGateway(ctx, "IAuthMarketplaceDAO"),
		infoLogs: env.GetLogsConfig().PrintMarketplaceDAOInfo,
	}
}

type DAOWithLog struct {
	original postgres.IAuthMarketplaceDAO
	logger   logger_lib.ILogger
	infoLogs bool
}

func (mdwl *DAOWithLog) GetUser(login string) (user models.User, password string, isFound bool, err error) {
	logLabel := fmt.Sprintf("GetUser():[params:%s]:", login)
	if mdwl.infoLogs {
		mdwl.logger.WriteInfo(fmt.Sprintf("%s %s", logLabel, "Received"))
	}
	user, password, isFound, err = mdwl.original.GetUser(login)
	if err != nil {
		mdwl.logger.WriteError(fmt.Sprintf("%s %v", logLabel, err))
	}
	return
}

func (mdwl *DAOWithLog) Registr(login, password string) (user models.User, isAlreadyExist bool, err error) {
	logLabel := fmt.Sprintf("Registr():[params:%s,%s]:", login, "***")
	if mdwl.infoLogs {
		mdwl.logger.WriteInfo(fmt.Sprintf("%s %s", logLabel, "Received"))
	}
	user, isAlreadyExist, err = mdwl.original.Registr(login, password)
	if err != nil {
		mdwl.logger.WriteError(fmt.Sprintf("%s %v", logLabel, err))
	}
	return
}

func (mdwl *DAOWithLog) Close() {
	if mdwl.infoLogs {
		mdwl.logger.WriteInfo("Close()")
	}
	mdwl.original.Close()
}
