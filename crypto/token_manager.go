package crypto

import (
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
	ValidateToken(token string) (user models.User, success bool, err error)
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
	logLabel := "TokenManager:GenerateToken():"

	id, err := uuid.NewRandom()
	if err != nil {
		err = fmt.Errorf("%s %v", logLabel, err)
		return
	}

	token = id.String()
	err = tm.tokensDAO.SetUser(token, user)

	if err != nil {
		err = fmt.Errorf("%s %v", logLabel, err)
		return
	}

	expirationTime := time.Now().Add(time.Hour * defaultExpirationHours)

	claims := &jwt.RegisteredClaims{
		Subject:   token,
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}

	tokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err = tokenJWT.SignedString([]byte(defaultSigningKey))
	if err != nil {
		err = fmt.Errorf("%s %v", logLabel, err)
		return
	}

	return
}

func (tm *TokenManager) ValidateToken(token string) (user models.User, success bool, err error) {
	logLabel := "TokenManager:ValidateToken():"

	claims := &jwt.RegisteredClaims{}
	tokenJWT, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		return []byte(defaultSigningKey), nil
	})
	if err != nil {
		err = fmt.Errorf("%s %v", logLabel, err)
		return
	}

	if !tokenJWT.Valid {
		err = fmt.Errorf("%s invalid token", logLabel)
		return
	}

	token = claims.Subject
	user, exist, err := tm.tokensDAO.GetUser(token)
	if !exist {
		if err != nil {
			err = fmt.Errorf("%s %v", logLabel, err)
			return
		}
		return
	}

	success = true
	return
}
