package DAO

import (
	"main/models"
	"sync"
)

type ITokensDAO interface {
	GetUser(token string) (user models.User, exist bool, err error)
	SetUser(token string, user models.User) (err error)
}

func CreateTokensDAO() ITokensDAO {
	return &TokensDAO{
		t:  make(map[string]models.User),
		tt: make(map[models.User]string),
		mu: sync.Mutex{},
	}
}

type TokensDAO struct {
	t  map[string]models.User
	tt map[models.User]string
	mu sync.Mutex
}

func (td *TokensDAO) GetUser(token string) (user models.User, exist bool, err error) {
	td.mu.Lock()
	defer td.mu.Unlock()
	user, exist = td.t[token]
	return
}

func (td *TokensDAO) SetUser(token string, user models.User) (err error) {
	td.mu.Lock()
	defer td.mu.Unlock()
	if t, exst := td.tt[user]; exst {
		delete(td.t, t)
		delete(td.tt, user)
	}
	td.tt[user] = token
	td.t[token] = user
	return
}
