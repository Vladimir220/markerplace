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
		mu: sync.Mutex{},
	}
}

type TokensDAO struct {
	t  map[string]models.User
	mu sync.Mutex
}

func (td *TokensDAO) GetUser(token string) (user models.User, exist bool, err error) {
	td.mu.Lock()
	defer td.mu.Unlock()
	user, exist = td.t[token]
	delete(td.t, token)
	return
}

func (td *TokensDAO) SetUser(token string, user models.User) (err error) {
	td.mu.Lock()
	defer td.mu.Unlock()
	td.t[token] = user
	return
}
