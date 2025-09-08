package authentication

import (
	"context"
	"fmt"
	"marketplace/crypto"
	"marketplace/db/DAO/postgres"
	"marketplace/env"
	"marketplace/network/auth/tools"

	"github.com/Vladimir220/markerplace/logger_lib"
)

type IAuthentication interface {
	Register(login, password string) (token string, err error)
	Login(login, password string) (token string, err error)
}

func CreateAuthentication(ctx context.Context, tokenManager crypto.ITokenManager, dao postgres.IMarketplaceDAO) IAuthentication {
	return &Authentication{
		tokenManager: tokenManager,
		dao:          dao,
		logger:       logger_lib.CreateLoggerGateway(ctx, "Authentication"),
		infoLogs:     env.GetLogsConfig().PrintAuthenticationInfo,
	}
}

type Authentication struct {
	tokenManager crypto.ITokenManager
	dao          postgres.IMarketplaceDAO
	logger       logger_lib.ILogger
	infoLogs     bool
}

func (auth *Authentication) Register(login, password string) (token string, err error) {
	logLabel := fmt.Sprintf("Register():[params:%s,%s]:", login, "***")
	ok := auth.checkLogin(login)
	if !ok {
		err = tools.ErrLoginFormat
		return
	}

	password, err = crypto.GetHashedPassword(password)
	if err != nil {
		auth.logger.WriteError(fmt.Sprintf("%s %v", logLabel, err))
		err = tools.ErrServer
		return
	}

	user, isAlreadyExist, err := auth.dao.Registr(login, password)
	if err != nil {
		err = tools.ErrServer
		return
	}
	if isAlreadyExist {
		err = tools.ErrLoginIsTaken
		return
	}

	token, isErr := auth.tokenManager.GenerateToken(user)
	if isErr {
		err = tools.ErrServer
		return
	}

	if auth.infoLogs {
		auth.logger.WriteInfo(fmt.Sprintf("%s %s", logLabel, "registered"))
	}
	return
}

func (auth *Authentication) Login(login, password string) (token string, err error) {
	logLabel := fmt.Sprintf("Login():[params:%s,%s]:", login, "***")

	user, realPassword, isFound, err := auth.dao.GetUser(login)
	if !isFound {
		err = tools.ErrLogin
		return
	}
	if err != nil {
		err = tools.ErrServer
		return
	}

	equal, err := crypto.ComparePassword(password, realPassword)
	if err != nil {
		auth.logger.WriteError(fmt.Sprintf("%s %v", logLabel, err))
		err = tools.ErrServer
		return
	}
	if !equal {
		err = tools.ErrLogin
		return
	}

	token, isErr := auth.tokenManager.GenerateToken(user)
	if isErr {
		err = tools.ErrServer
		return
	}

	if auth.infoLogs {
		auth.logger.WriteInfo(fmt.Sprintf("%s %s", logLabel, "login"))
	}

	return
}
func (auth *Authentication) checkLogin(login string) bool {
	if len(login) <= 0 || len(login) > 30 {
		return false
	} else {
		return true
	}
}
