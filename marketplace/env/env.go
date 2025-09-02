package env

import (
	"errors"
	"os"
	"strings"
)

type KafkaEnvData struct {
	loggerData struct {
		InfoTopicName    string
		ErrorTopicName   string
		WarningTopicName string
	}

	writerData struct {
		NewAnnouncementTopicName string
	}

	BrokerHosts []string
}

type PostgresEnvData struct {
	User, Password, DbName, Host string
}

func GetKafkaEnvData() (data KafkaEnvData, err error) {
	brokerHostsStr := os.Getenv("KAFKA_LOGGER_HOST")
	data.loggerData.InfoTopicName = os.Getenv("INFO_TOPIC_NAME")
	data.loggerData.ErrorTopicName = os.Getenv("ERROR_TOPIC_NAME")
	data.loggerData.WarningTopicName = os.Getenv("WARNING_TOPIC_NAME")
	data.writerData.NewAnnouncementTopicName = os.Getenv("NEW_ANNOUNCEMENT_TOPIC_NAME")

	if brokerHostsStr == "" || data.loggerData.InfoTopicName == "" || data.loggerData.ErrorTopicName == "" || data.loggerData.WarningTopicName == "" || data.writerData.NewAnnouncementTopicName == "" {
		err = errors.New("GetKafkaEnvData(): one of the following variables is not specified in .env: KAFKA_BROKER_HOSTS, INFO_TOPIC_NAME, ERROR_TOPIC_NAME, WARNING_TOPIC_NAME, NEW_ANNOUNCEMENT_TOPIC_NAME")
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
	serviceName = os.Getenv("DB_PASSWORD")
	if host == "" || serviceName == "" {
		err = errors.New("GetServiceData(): one of the following variables is not specified in .env: GRPC_HOST, DB_PASSWORD")
		return
	}
	return
}
