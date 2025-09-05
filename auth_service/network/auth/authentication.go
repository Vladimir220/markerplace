package auth

import (
	"auth_service/crypto"
	"auth_service/db/DAO/postgres"
	"context"
	"errors"
	"fmt"

	"github.com/Vladimir220/markerplace/logger_lib"
)

var ErrLogin = errors.New("Incorrect password or login")
var ErrLoginFormat = errors.New("Incorrect login format")
var ErrPasswordFormat = errors.New("Incorrect password format")
var ErrLoginIsTaken = errors.New("Login is taken")
var ErrServer = errors.New("Server error")

type IAuthentication interface {
	Register(login, password string) (token string, err error)
	Login(login, password string) (token string, err error)
}

func CreateAuthentication(ctx context.Context, tokenManager crypto.ITokenManager, infoLogs bool) (IAuthentication, error) {
	dao, err := postgres.CreateMarketplaceDAO()
	if err != nil {
		return nil, fmt.Errorf("CreateAuthentication():%v", err)
	}
	return &Authentication{
		tokenManager: tokenManager,
		dao:          dao,
		logger:       logger_lib.CreateLoggerAdapter(ctx, "Authentication"),
		infoLogs:     infoLogs,
	}, nil
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
		err = ErrLoginFormat
		return
	}

	password, err = crypto.GetHashedPassword(password)
	if err != nil {
		auth.logger.WriteError(fmt.Sprintf("%s %v", logLabel, err))
		err = ErrServer
		return
	}

	user, isAlreadyExist, err := auth.dao.Registr(login, password)
	if err != nil {
		err = ErrServer
		return
	}
	if isAlreadyExist {
		err = ErrLoginIsTaken
		return
	}

	token, isErr := auth.tokenManager.GenerateToken(user)
	if isErr {
		err = ErrServer
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
		err = ErrLogin
		return
	}
	if err != nil {
		err = ErrServer
		return
	}

	equal, err := crypto.ComparePassword(password, realPassword)
	if err != nil {
		auth.logger.WriteError(fmt.Sprintf("%s %v", logLabel, err))
		err = ErrServer
		return
	}
	if !equal {
		err = ErrLogin
		return
	}

	token, isErr := auth.tokenManager.GenerateToken(user)
	if isErr {
		err = ErrServer
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
