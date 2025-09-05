package env

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type KafkaEnvData struct {
	InfoTopicName    string
	ErrorTopicName   string
	WarningTopicName string
	BrokerHosts      []string
}

type PostgresEnvData struct {
	User, Password, DbName, Host string
}

type RedisEnvData struct {
	Host, Password string
	DbNum          int
}

func GetKafkaEnvData() (data KafkaEnvData, err error) {
	brokerHostsStr := os.Getenv("KAFKA_BROKER_HOSTS")
	data.InfoTopicName = os.Getenv("INFO_TOPIC_NAME")
	data.ErrorTopicName = os.Getenv("ERROR_TOPIC_NAME")
	data.WarningTopicName = os.Getenv("WARNING_TOPIC_NAME")

	if brokerHostsStr == "" || data.InfoTopicName == "" || data.ErrorTopicName == "" || data.WarningTopicName == "" {
		err = errors.New("GetKafkaEnvData(): one of the following variables is not specified in .env: KAFKA_BROKER_HOSTS, INFO_TOPIC_NAME, ERROR_TOPIC_NAME, WARNING_TOPIC_NAME")
		return
	}

	data.BrokerHosts = strings.Split(strings.TrimSpace(brokerHostsStr), ",")

	return
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
