package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"writer_db_service/db/DAO/postgres"
	"writer_db_service/env"
	"writer_db_service/log/proxies"

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

	dao, err := postgres.CreateWriterMarketplaceDAO()
	if err != nil {
		panic(err)
	}
	daoWithLogs := proxies.CreateDAOWithLog(ctx, dao)
	defer daoWithLogs.Close()

	RunKafkaWorkers(ctx, daoWithLogs, WorkersConfig{
		brokerHosts: kafkaKonfig.BrokerHosts,
		topic:       kafkaKonfig.NewAnnouncementTopicName,
		count:       kafkaKonfig.NumOfWorkers,
		mod:         ModNewAnnouncement,
		groupId:     kafkaKonfig.GroupId,
	})

	RunKafkaWorkers(ctx, daoWithLogs, WorkersConfig{
		brokerHosts: kafkaKonfig.BrokerHosts,
		topic:       kafkaKonfig.UpdateAnnouncementTopicName,
		count:       kafkaKonfig.NumOfWorkers,
		mod:         ModUpdateAnnouncement,
		groupId:     kafkaKonfig.GroupId,
	})

	RunKafkaWorkers(ctx, daoWithLogs, WorkersConfig{
		brokerHosts: kafkaKonfig.BrokerHosts,
		topic:       kafkaKonfig.DeleteAnnouncementTopicName,
		count:       kafkaKonfig.NumOfWorkers,
		mod:         ModDeleteAnnouncement,
		groupId:     kafkaKonfig.GroupId,
	})

	go HealthListener()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	defer stop()
	<-ctx.Done()
}
