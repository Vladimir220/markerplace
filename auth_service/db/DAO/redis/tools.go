package redis

import (
	"auth_service/env"
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func connect() (connection *redis.Client, ctx context.Context, err error) {
	logLable := "connect()"
	data, err := env.GetRedisEnvData()
	if err != nil {
		err = fmt.Errorf("%s%v", logLable, err)
		return
	}

	connection = redis.NewClient(&redis.Options{
		Addr:     data.Host,
		Password: data.Password,
		DB:       data.DbNum,
	})

	_, err = connection.Ping(context.Background()).Result()
	if err != nil {
		err = fmt.Errorf("%s %v", logLable, err)
		return
	}

	ctx = context.Background()

	connection.FlushDB(ctx)

	return
}
