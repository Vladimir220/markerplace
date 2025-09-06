package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"writer_db_service/env"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	kafkaKonfig, err := env.GetKafkaEnvData()
	if err != nil {
		panic(err)
	}

	RunKafkaWorkers(ctx, WorkersConfig{
		brokerHosts: kafkaKonfig.BrokerHosts,
		topic:       kafkaKonfig.NewAnnouncementTopicName,
		count:       kafkaKonfig.NumOfWorkers,
		mod:         ModNewAnnouncement,
		groupId:     kafkaKonfig.GroupId,
	})

	fmt.Println("еее, здоровье!")
	go HealthListener()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	defer stop()
	<-ctx.Done()
}
