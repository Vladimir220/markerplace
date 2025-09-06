package logger_lib

import (
	"errors"
	"os"
	"strings"
)

type KafkaEnvData struct {
	InfoTopicName    string
	ErrorTopicName   string
	WarningTopicName string
	BrokerHosts      []string
}

type LoggerConfig struct {
	ErrorsToStdOut   bool
	WarningsToStdOut bool
	InfoToStdOut     bool
}

func GetLoggerConfig() (data LoggerConfig) {
	printErrorsToStdOutStr := os.Getenv("PRINT_ERRORS_TO_STD_OUT")
	printWarningsToStdStr := os.Getenv("PRINT_WARNINGS_TO_STD_OUT")
	printInfoToStdOutStr := os.Getenv("PRINT_INFO_TO_STD_OUT")

	if printErrorsToStdOutStr != "" {
		printErrorsToStdOutStr = strings.ToUpper(printErrorsToStdOutStr)
		if printErrorsToStdOutStr == "TRUE" {
			data.ErrorsToStdOut = true
		}
	}

	if printWarningsToStdStr != "" {
		printWarningsToStdStr = strings.ToUpper(printWarningsToStdStr)
		if printWarningsToStdStr == "TRUE" {
			data.WarningsToStdOut = true
		}
	}

	if printInfoToStdOutStr != "" {
		printInfoToStdOutStr = strings.ToUpper(printInfoToStdOutStr)
		if printInfoToStdOutStr == "TRUE" {
			data.InfoToStdOut = true
		}
	}

	return
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
