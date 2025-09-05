package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"writer_db_service/DAO/postgres"
	"writer_db_service/models"

	"github.com/Vladimir220/markerplace/logger_lib"
	"github.com/segmentio/kafka-go"
)

var ModNewAnnouncement = 0

type WorkersConfig struct {
	brokerHosts []string
	topic       string
	count       uint
	mod         int
	groupId     string
}

func RunKafkaWorkers(ctx context.Context, config WorkersConfig) {
	for range config.count {
		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers:          config.brokerHosts,
			Topic:            config.topic,
			GroupID:          config.groupId,
			StartOffset:      kafka.LastOffset,
			SessionTimeout:   time.Minute * 30,
			RebalanceTimeout: time.Second * 30,
		})

		go KafkaWorker(ctx, reader, config.mod)
	}
}

func KafkaWorker(ctx context.Context, reader *kafka.Reader, mod int) {
	logger := logger_lib.CreateLoggerAdapter(ctx, "KafkaWorker()")
	defer reader.Close()
	dao, err := postgres.CreateMarketplaceDAO()
	if err != nil {
		logger.WriteError(err.Error())
		panic(err)
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			logger.WriteError(err.Error())
			time.Sleep(time.Second)
			continue
		}

		fmt.Println("я отработал")

		var announcement models.ExtendedAnnouncement
		err = json.Unmarshal([]byte(msg.Value), &announcement)
		if err != nil {
			logger.WriteError(err.Error())

			err = reader.CommitMessages(ctx, msg)
			if err != nil {
				logger.WriteError(err.Error())
				time.Sleep(time.Second)
			}
			continue
		}

		switch mod {
		case ModNewAnnouncement:
			err = dao.NewAnnouncement(announcement)
			if err != nil {
				logger.WriteError(err.Error())

				err = reader.CommitMessages(ctx, msg)
				if err != nil {
					logger.WriteError(err.Error())
					time.Sleep(time.Second)
				}
				continue
			}
		}

		err = reader.CommitMessages(ctx, msg)
		if err != nil {
			logger.WriteError(err.Error())
			time.Sleep(time.Second)
			continue
		}
	}
}
