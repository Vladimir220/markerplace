package crypto

import (
	"errors"
	"fmt"
	"main/db/DAO"
	"main/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const defaultSigningKey = "secretKEY123"
const defaultExpirationHours = 256

type ITokenManager interface {
	GenerateToken(user models.User) (token string, err error)
	ValidateToken(token string) (user models.User, err error)
}

func CreateTokenManager() ITokenManager {
	return &TokenManager{
		tokensDAO: DAO.CreateTokensDAO(),
	}
}

type TokenManager struct {
	tokensDAO DAO.ITokensDAO
}

func (tm *TokenManager) GenerateToken(user models.User) (token string, err error) {
	id, err := uuid.NewRandom()
	if err != nil {
		err = fmt.Errorf("TokenManager:GenerateToken(): %v", err)
		return
	}

	token = id.String()
	err = tm.tokensDAO.SetUser(token, user)
	if err != nil {
		err = fmt.Errorf("TokenManager:GenerateToken(): %v", err)
		return
	}

	expirationTime := time.Now().Add(defaultExpirationHours)

	claims := &jwt.RegisteredClaims{
		Subject:   token,
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}

	tokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err = tokenJWT.SignedString(defaultSigningKey)
	if err != nil {
		err = fmt.Errorf("TokenManager:GenerateToken(): %v", err)
		return
	}

	return
}

func (tm *TokenManager) ValidateToken(token string) (user models.User, err error) {

	claims := &jwt.RegisteredClaims{}
	tokenJWT, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		return defaultSigningKey, nil
	})

	if err != nil {
		err = fmt.Errorf("TokenManager:ValidateToken(): %v", err)
		return
	}

	if !tokenJWT.Valid {
		err = errors.New("TokenManager:ValidateToken(): invalid token")
		return
	}

	token = claims.Subject
	user, exist, err := tm.tokensDAO.GetUser(token)
	if !exist {
		err = errors.New("TokenManager:ValidateToken(): invalid token")
		return
	}
	if err != nil {
		err = fmt.Errorf("TokenManager:ValidateToken(): %v", err)
		return
	}

	return
}
