package redis

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func connect() (connection *redis.Client, ctx context.Context, err error) {
	logLable := "connect()"
	host, password, dbNum, err := getEnvLoginData()
	if err != nil {
		err = fmt.Errorf("%s%v", logLable, err)
		return
	}

	connection = redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       dbNum,
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

func getEnvLoginData() (host, password string, dbNum int, err error) {
	logLabel := "getEnvLoginData():"

	host = os.Getenv("REDIS_HOST")
	password = os.Getenv("REDIS_PASSWORD")
	dbStr := os.Getenv("REDIS_DB")

	if host == "" || dbStr == "" {
		err = fmt.Errorf("%s one of the following variables is not specified in env: REDIS_HOST, REDIS_PASSWORD, REDIS_DB", logLabel)
		return
	}

	dbNum, err = strconv.Atoi(dbStr)
	if err != nil {
		err = fmt.Errorf("%s %v", logLabel, err)
		return
	}

	return
}
