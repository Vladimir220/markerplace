package main

import (
	"context"
	"writer_db_service/env"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	kafkaKonfig, err := env.GetKafkaEnvData()
	if err != nil {
		panic(err)
	}

	RunKafkaWorkers(ctx, WorkersConfig{
		brokerHosts: kafkaKonfig.BrokerHosts,
		topic:       kafkaKonfig.WarningTopicName,
		count:       kafkaKonfig.NumOfWorkers,
		mod:         ModNewAnnouncement,
	})
}
