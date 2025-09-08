package env

import (
	"errors"
	"os"
	"strings"
)

type PostgresEnvData struct {
	User, Password, DbName, Host string
}

type LogsConfig struct {
	PrintMarketplaceDAOInfo bool
}

func GetLogsConfig() (data LogsConfig) {
	printMarketplaceDAOInfoStr := os.Getenv("PRINT_MARKETPLACE_DAO_INFO")
	if strings.ToUpper(printMarketplaceDAOInfoStr) == "TRUE" {
		data.PrintMarketplaceDAOInfo = true
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
