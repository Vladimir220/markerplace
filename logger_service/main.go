package main

import "context"

func main() {
	envData, err := GetKafkaEnvData()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	RunKafkaWorkers(ctx, WorkersConfig{
		brokerHosts: envData.BrokerHosts,
		topic:       envData.WarningTopicName,
		count:       envData.NumOfWorkers,
		mod:         ModWarning,
	})

	RunKafkaWorkers(ctx, WorkersConfig{
		brokerHosts: envData.BrokerHosts,
		topic:       envData.InfoTopicName,
		count:       envData.NumOfWorkers,
		mod:         ModInfo,
	})

	RunKafkaWorkers(ctx, WorkersConfig{
		brokerHosts: envData.BrokerHosts,
		topic:       envData.ErrorTopicName,
		count:       envData.NumOfWorkers,
		mod:         ModError,
	})
}
