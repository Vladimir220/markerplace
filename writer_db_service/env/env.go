package env

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type KafkaEnvData struct {
	NewAnnouncementTopicName string
	InfoTopicName            string
	ErrorTopicName           string
	WarningTopicName         string
	BrokerHosts              []string
	NumOfWorkers             uint
	GroupId                  string
}

type PostgresEnvData struct {
	User, Password, DbName, Host string
}

func GetKafkaEnvData() (data KafkaEnvData, err error) {
	logLabel := "GetKafkaEnvData()"

	brokerHostsStr := os.Getenv("KAFKA_BROKER_HOSTS")
	countStr := os.Getenv("NUM_OF_WORKERS_FOR_EACH_TOPIC")
	data.GroupId = os.Getenv("GROUP_ID_FOR_ALL_TOPICS")
	data.NewAnnouncementTopicName = os.Getenv("NEW_ANNOUNCEMENT_TOPIC_NAME")
	data.InfoTopicName = os.Getenv("INFO_TOPIC_NAME")
	data.ErrorTopicName = os.Getenv("ERROR_TOPIC_NAME")
	data.WarningTopicName = os.Getenv("WARNING_TOPIC_NAME")

	if brokerHostsStr == "" || countStr == "" || data.GroupId == "" || data.NewAnnouncementTopicName == "" || data.InfoTopicName == "" || data.ErrorTopicName == "" || data.WarningTopicName == "" {
		err = errors.New("GetKafkaEnvData(): one of the following variables is not specified in .env: KAFKA_BROKER_HOSTS, NUM_OF_WORKERS_FOR_EACH_TOPIC, GROUP_ID_FOR_ALL_TOPICS, NEW_ANNOUNCEMENT_TOPIC_NAME, INFO_TOPIC_NAME, ERROR_TOPIC_NAME, WARNING_TOPIC_NAME")
		return
	}

	data.BrokerHosts = strings.Split(strings.TrimSpace(brokerHostsStr), ",")

	count, err := strconv.ParseUint(countStr, 10, 64)
	if err != nil {
		err = fmt.Errorf("%s %v", logLabel, err)
		return
	}
	data.NumOfWorkers = uint(count)

	return
}

func GetPostgresEnvData() (data PostgresEnvData, err error) {
	data.User = os.Getenv("DB_USER")
	data.Password = os.Getenv("DB_PASSWORD")
	data.DbName = os.Getenv("DB_NAME")
	data.Host = os.Getenv("DB_HOST")
	if data.User == "" || data.Password == "" || data.DbName == "" || data.Host == "" {
		err = errors.New("getEnvLoginData(): one of the following variables is not specified in .env: DB_USER, DB_PASSWORD, DB_NAME, DB_HOST")
		return
	}
	return
}
