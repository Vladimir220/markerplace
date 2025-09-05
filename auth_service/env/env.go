package env

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type PostgresEnvData struct {
	User, Password, DbName, Host string
}

type RedisEnvData struct {
	Host, Password string
	DbNum          int
}

func GetPostgresEnvData() (data PostgresEnvData, err error) {
	data.User = os.Getenv("DB_USER")
	data.Password = os.Getenv("DB_PASSWORD")
	data.DbName = os.Getenv("DB_NAME")
	data.Host = os.Getenv("DB_HOST")
	if data.User == "" || data.Password == "" || data.DbName == "" || data.Host == "" {
		err = errors.New("GetEnvLoginData(): one of the following variables is not specified in .env: DB_USER, DB_PASSWORD, DB_NAME, DB_HOST")
		return
	}
	return
}

func GetServiceData() (host, serviceName string, err error) {
	host = os.Getenv("GRPC_HOST")
	serviceName = os.Getenv("SERVICE_NAME")
	if host == "" || serviceName == "" {
		err = errors.New("GetServiceData(): one of the following variables is not specified in .env: GRPC_HOST, DB_PASSWORD")
		return
	}
	return
}

func GetRedisEnvData() (data RedisEnvData, err error) {
	logLabel := "GetEnvLoginData():"

	data.Host = os.Getenv("REDIS_HOST")
	data.Password = os.Getenv("REDIS_PASSWORD")
	dbStr := os.Getenv("REDIS_DB")

	if data.Host == "" || dbStr == "" {
		err = fmt.Errorf("%s one of the following variables is not specified in env: REDIS_HOST, REDIS_PASSWORD, REDIS_DB", logLabel)
		return
	}

	data.DbNum, err = strconv.Atoi(dbStr)
	if err != nil {
		err = fmt.Errorf("%s %v", logLabel, err)
		return
	}

	return
}
