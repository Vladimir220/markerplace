package crypto

import (
	"auth_service/db/DAO"
	"auth_service/models"
	"context"
	"fmt"
	"time"

	"github.com/Vladimir220/markerplace/logger_lib"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const defaultSigningKey = "secretKEY123"
const defaultExpirationHours = 256

type ITokenManager interface {
	GenerateToken(user models.User) (token string, isErr bool)
	ValidateToken(token string) (user models.User, isValid, isErr bool)
}

func CreateTokenManager(ctx context.Context, tokensDAO DAO.ITokensDAO) ITokenManager {
	return &TokenManager{
		tokensDAO: tokensDAO,
		logger:    logger_lib.CreateLoggerGateway(ctx, "TokenManager"),
	}
}

type TokenManager struct {
	tokensDAO DAO.ITokensDAO
	logger    logger_lib.ILogger
}

func (tm *TokenManager) GenerateToken(user models.User) (token string, isErr bool) {
	logLabel := fmt.Sprintf("GenerateToken():[params:%v]:", user)

	id, err := uuid.NewRandom()
	if err != nil {
		tm.logger.WriteError(fmt.Sprintf("%s %v", logLabel, err))
		isErr = true
		return
	}

	token = id.String()
	err = tm.tokensDAO.SetUser(token, user)
	if err != nil {
		isErr = true
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
		tm.logger.WriteError(fmt.Sprintf("%s %v", logLabel, err))
		isErr = true
		return
	}

	return
}

func (tm *TokenManager) ValidateToken(token string) (user models.User, isValid, isErr bool) {
	logLabel := fmt.Sprintf("ValidateToken():[params:%s]:", token)

	claims := &jwt.RegisteredClaims{}
	tokenJWT, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		return []byte(defaultSigningKey), nil
	})
	if err != nil {
		tm.logger.WriteError(fmt.Sprintf("%s %v", logLabel, err))
		isErr = true
		return
	}

	if !tokenJWT.Valid {
		return
	} else {
		isValid = true
	}

	token = claims.Subject
	user, exist, err := tm.tokensDAO.GetUser(token)
	if !exist || err != nil {
		isValid = false
		return
	}

	return
}
