package env

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type PostgresEnvData struct {
	User, Password, DbName, Host string
}

type RedisEnvData struct {
	Host, Password string
	DbNum          int
}

type LogsConfig struct {
	PrintTokenDAOInfo       bool
	PrintMarketplaceDAOInfo bool
	PrintAuthenticationInfo bool
	PrintHandlersInfo       bool
}

func GetLogsConfig() (data LogsConfig) {
	printTokenDAOInfoStr := os.Getenv("PRINT_TOKEN_DAO_INFO")
	if strings.ToUpper(printTokenDAOInfoStr) == "TRUE" {
		data.PrintTokenDAOInfo = true
	}

	printMarketplaceDAOInfoStr := os.Getenv("PRINT_MARKETPLACE_DAO_INFO")
	if strings.ToUpper(printMarketplaceDAOInfoStr) == "TRUE" {
		data.PrintMarketplaceDAOInfo = true
	}

	printAuthenticationInfoStr := os.Getenv("PRINT_AUTHENTICATION_INFO")
	if strings.ToUpper(printAuthenticationInfoStr) == "TRUE" {
		data.PrintAuthenticationInfo = true
	}

	printHandlersInfoStr := os.Getenv("PRINT_HANDLERS_INFO")
	if strings.ToUpper(printHandlersInfoStr) == "TRUE" {
		data.PrintHandlersInfo = true
	}

	return
}

func GetPostgresEnvData() (data PostgresEnvData, err error) {
	data.User = os.Getenv("DB_USER")
	data.DbName = os.Getenv("DB_NAME")
	data.Password = os.Getenv("DB_PASSWORD")
	data.Host = os.Getenv("DB_HOST")
	if data.User == "" || data.DbName == "" || data.Host == "" {
		err = errors.New("GetEnvLoginData(): one of the following variables is not specified in .env: DB_USER, DB_NAME, DB_HOST")
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
