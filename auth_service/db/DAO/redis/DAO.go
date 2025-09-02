package redis

import (
	"auth_service/db/DAO"
	"auth_service/models"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ttlHours   = 24
	userLabel  = "tokens:user:"
	tokenLabel = "tokens:token:"
)

func CreateTokensDAO() (dao DAO.ITokensDAO, err error) {
	connection, ctx, err := connect()
	if err != nil {
		err = fmt.Errorf("CreateTokensDAO(): %v", err)
		return
	}

	return &TokensDAO{
		connection: connection,
		ctx:        ctx,
	}, nil
}

type TokensDAO struct {
	connection *redis.Client
	ctx        context.Context
}

func (td *TokensDAO) GetUser(token string) (user models.User, exist bool, err error) {
	logLabel := "GetUser():"
	localTokenLable := fmt.Sprintf("%s%s", tokenLabel, token)

	u, err := td.connection.HGetAll(td.ctx, localTokenLable).Result()

	if err != nil {
		err = fmt.Errorf("%s %v", logLabel, err)
		return
	} else if len(u) == 0 {
		return
	} else {
		exist = true
	}

	user.Login = u["login"]
	user.Group = u["group"]

	return
}

func (td *TokensDAO) SetUser(token string, user models.User) (err error) {
	logLabel := "SetUser():"
	localUserLable := fmt.Sprintf("%s%s", userLabel, user.Login)
	localTokenLable := fmt.Sprintf("%s%s", tokenLabel, token)

	t, err := td.connection.Get(td.ctx, localUserLable).Result()
	if err != redis.Nil {
		if err != nil {
			err = fmt.Errorf("%s %v", logLabel, err)
			return
		} else {
			err = td.connection.Del(td.ctx, fmt.Sprintf("%s%s", tokenLabel, t)).Err()
			if err != nil {
				err = fmt.Errorf("%s %v", logLabel, err)
				return
			}
			err = td.connection.Del(td.ctx, localUserLable).Err()
			if err != nil {
				err = fmt.Errorf("%s %v", logLabel, err)
				return
			}
		}
	}

	err = td.connection.HSet(td.ctx, localTokenLable, user).Err()
	if err != nil {
		err = fmt.Errorf("%s %v", logLabel, err)
		return
	}
	err = td.connection.Expire(td.ctx, localTokenLable, time.Hour*time.Duration(ttlHours)).Err()
	if err != nil {
		err = fmt.Errorf("%s %v", logLabel, err)
		return
	}

	err = td.connection.Set(td.ctx, localUserLable, token, time.Hour*time.Duration(ttlHours)).Err()
	if err != nil {
		err = fmt.Errorf("%s %v", logLabel, err)
		return
	}

	return
}
