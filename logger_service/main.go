package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

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
		groupId:     envData.GroupId,
	})

	RunKafkaWorkers(ctx, WorkersConfig{
		brokerHosts: envData.BrokerHosts,
		topic:       envData.InfoTopicName,
		count:       envData.NumOfWorkers,
		mod:         ModInfo,
		groupId:     envData.GroupId,
	})

	RunKafkaWorkers(ctx, WorkersConfig{
		brokerHosts: envData.BrokerHosts,
		topic:       envData.ErrorTopicName,
		count:       envData.NumOfWorkers,
		mod:         ModError,
		groupId:     envData.GroupId,
	})

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	defer stop()
	<-ctx.Done()
}
