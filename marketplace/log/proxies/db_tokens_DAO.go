package proxies

import (
	"context"
	"fmt"
	"marketplace/db/DAO"
	"marketplace/env"
	"marketplace/models"

	"github.com/Vladimir220/markerplace/logger_lib"
)

func CreateTokensDAOWithLog(ctx context.Context, dao DAO.ITokensDAO) DAO.ITokensDAO {
	return &TokensDAOWithLog{
		original: dao,
		logger:   logger_lib.CreateLoggerGateway(ctx, "ITokensDAO"),
		infoLogs: env.GetLogsConfig().PrintTokenDAOInfo,
	}
}

type TokensDAOWithLog struct {
	original DAO.ITokensDAO
	logger   logger_lib.ILogger
	infoLogs bool
}

func (tdwl *TokensDAOWithLog) GetUser(token string) (user models.User, exist bool, err error) {
	logLabel := fmt.Sprintf("GetUser():[params:%s]:", token)
	if tdwl.infoLogs {
		tdwl.logger.WriteInfo(fmt.Sprintf("%s %s", logLabel, "Received"))
	}
	user, exist, err = tdwl.original.GetUser(token)
	if err != nil {
		tdwl.logger.WriteError(fmt.Sprintf("%s %v", logLabel, err))
	}
	return
}

func (tdwl *TokensDAOWithLog) SetUser(token string, user models.User) (err error) {
	logLabel := fmt.Sprintf("SetUser():[params:%s,%s]:", token, user)
	if tdwl.infoLogs {
		tdwl.logger.WriteInfo(fmt.Sprintf("%s %s", logLabel, "Received"))
	}
	err = tdwl.original.SetUser(token, user)
	if err != nil {
		tdwl.logger.WriteError(fmt.Sprintf("%s %v", logLabel, err))
	}
	return
}

func (tdwl *TokensDAOWithLog) Close() {
	if tdwl.infoLogs {
		tdwl.logger.WriteInfo("Close()")
	}
	tdwl.original.Close()
}
