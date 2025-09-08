package env

import (
	"errors"
	"os"
	"strings"
)

type KafkaEnvData struct {
	WriterData struct {
		NewAnnouncementTopicName    string
		UpdateAnnouncementTopicName string
		DeleteAnnouncementTopicName string
	}

	BrokerHosts []string
}

type PostgresEnvData struct {
	User, Password, DbName, Host string
}

type ServicesData struct {
	ReaderGRPCHost string
	AuthGRPCHost   string
}

type LogsConfig struct {
	PrintTokenDAOInfo       bool
	PrintMarketplaceDAOInfo bool
	PrintAuthorizationInfo  bool
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

	printAuthorizationInfoStr := os.Getenv("PRINT_AUTHORIZATION_INFO")
	if strings.ToUpper(printAuthorizationInfoStr) == "TRUE" {
		data.PrintAuthorizationInfo = true
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

func GetKafkaEnvData() (data KafkaEnvData, err error) {
	brokerHostsStr := os.Getenv("KAFKA_BROKER_HOSTS")
	data.WriterData.NewAnnouncementTopicName = os.Getenv("NEW_ANNOUNCEMENT_TOPIC_NAME")
	data.WriterData.UpdateAnnouncementTopicName = os.Getenv("UPDATE_ANNOUNCEMENT_TOPIC_NAME")
	data.WriterData.DeleteAnnouncementTopicName = os.Getenv("DELETE_ANNOUNCEMENT_TOPIC_NAME")

	if brokerHostsStr == "" || data.WriterData.NewAnnouncementTopicName == "" {
		err = errors.New("GetKafkaEnvData(): one of the following variables is not specified in .env: KAFKA_BROKER_HOSTS, NEW_ANNOUNCEMENT_TOPIC_NAME, UPDATE_ANNOUNCEMENT_TOPIC_NAME, DELETE_ANNOUNCEMENT_TOPIC_NAME")
		return
	}

	data.BrokerHosts = strings.Split(strings.TrimSpace(brokerHostsStr), ",")

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
	host = os.Getenv("HOST")
	serviceName = os.Getenv("SERVICE_NAME")
	if host == "" || serviceName == "" {
		err = errors.New("GetServiceData(): one of the following variables is not specified in .env: GRPC_HOST, DB_PASSWORD")
		return
	}
	return
}

func GetServicesData() (data ServicesData, err error) {
	data.AuthGRPCHost = os.Getenv("AUTH_SERVICE_GRPC_HOST")
	data.ReaderGRPCHost = os.Getenv("READER_SERVICE_GRPC_HOST")

	if data.AuthGRPCHost == "" || data.ReaderGRPCHost == "" {
		err = errors.New("GetServicesData(): one of the following variables is not specified in .env: AUTH_SERVICE_GRPC_HOST, READER_SERVICE_GRPC_HOST")
		return
	}

	return
}
