package main

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
	NumOfWorkers     uint
	GroupId          string
}

func GetKafkaEnvData() (data KafkaEnvData, err error) {
	logLabel := "getKafkaEnvData()"

	brokerHostsStr := os.Getenv("KAFKA_BROKER_HOSTS")
	countStr := os.Getenv("NUM_OF_WORKERS_FOR_EACH_TOPIC")
	data.GroupId = os.Getenv("GROUP_ID_FOR_ALL_TOPICS")
	data.InfoTopicName = os.Getenv("INFO_TOPIC_NAME")
	data.ErrorTopicName = os.Getenv("ERROR_TOPIC_NAME")
	data.WarningTopicName = os.Getenv("WARNING_TOPIC_NAME")

	if brokerHostsStr == "" || countStr == "" || data.GroupId == "" || data.InfoTopicName == "" || data.ErrorTopicName == "" || data.WarningTopicName == "" {
		err = errors.New("GetKafkaEnvData(): one of the following variables is not specified in .env: KAFKA_BROKER_HOSTS, NUM_OF_WORKERS_FOR_EACH_TOPIC, GROUP_ID_FOR_ALL_TOPICS, INFO_TOPIC_NAME, ERROR_TOPIC_NAME, WARNING_TOPIC_NAME")
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
