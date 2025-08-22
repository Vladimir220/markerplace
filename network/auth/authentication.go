package auth

import (
	"errors"
	"fmt"
	"main/crypto"
	"main/db/DAO/postgres"
)

type IAuthentication interface {
	Register(login, password string) (token string, err error)
	Login(login, password string) (token string, err error)
}

func CreateAuthentication(tokenManager crypto.ITokenManager) IAuthentication {
	return &Authentication{
		tokenManager: tokenManager,
		dao:          postgres.CreateMarketplaceDAO(),
	}
}

type Authentication struct {
	tokenManager crypto.ITokenManager
	dao          postgres.IMarketplaceDAO
}

func (auth *Authentication) Register(login, password string) (token string, err error) {
	ok := auth.checkLogin(login)
	if !ok {
		err = errors.New("Authentication:Register: incorrect login")
		return
	}

	password, err = crypto.GetHashedPassword(password)
	if err != nil {
		err = fmt.Errorf("Authentication:Registr: %v", err)
		return
	}

	user, _, err := auth.dao.Registr(login, password)
	if err != nil {
		err = fmt.Errorf("Authentication:Registr: %v", err)
		return
	}

	token, err = auth.tokenManager.GenerateToken(user)

	return
}

func (auth *Authentication) Login(login, password string) (token string, err error) {
	user, realPassword, isFound, err := auth.dao.GetUser(login)
	if !isFound {
		err = fmt.Errorf("Authentication:Login: user '%s' not found", login)
		return
	}
	if err != nil {
		err = fmt.Errorf("Authentication:Login: %v", err)
		return
	}

	equal := crypto.ComparePassword(password, realPassword)
	if !equal {
		err = errors.New("Authentication:Login: incorrect password")
		return
	}

	token, err = auth.tokenManager.GenerateToken(user)

	return
}
func (auth *Authentication) checkLogin(login string) bool {
	if len(login) <= 0 || len(login) > 30 {
		return false
	} else {
		return true
	}
}
