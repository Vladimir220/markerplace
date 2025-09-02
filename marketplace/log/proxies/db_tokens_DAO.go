package proxies

import (
	"fmt"
	"main/db/DAO"
	"main/log"
	"main/models"
)

func CreateTokensDAOWithLog(dao DAO.ITokensDAO, infoLogs bool) DAO.ITokensDAO {
	return &TokensDAOWithLog{
		original: dao,
		logger:   log.CreateLogger("ITokensDAO"),
		infoLogs: infoLogs,
	}
}

type TokensDAOWithLog struct {
	original DAO.ITokensDAO
	logger   log.ILogger
	infoLogs bool
}

func (tdwl *TokensDAOWithLog) GetUser(token string) (user models.User, exist bool, err error) {
	logLabel := fmt.Sprintf("GetUser():[params:%s]:", token)
	if tdwl.infoLogs {
		tdwl.logger.WriteInfo(fmt.Sprintf("%s %s", logLabel, "get"))
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
		tdwl.logger.WriteInfo(fmt.Sprintf("%s %s", logLabel, "get"))
	}
	err = tdwl.original.SetUser(token, user)
	if err != nil {
		tdwl.logger.WriteError(fmt.Sprintf("%s %v", logLabel, err))
	}
	return
}
