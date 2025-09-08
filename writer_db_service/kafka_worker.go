package main

import (
	"context"
	"encoding/json"
	"strconv"
	"time"
	"writer_db_service/db/DAO/postgres"
	"writer_db_service/models"

	"github.com/Vladimir220/markerplace/logger_lib"
	"github.com/segmentio/kafka-go"
)

var ModNewAnnouncement = 0
var ModDeleteAnnouncement = 1
var ModUpdateAnnouncement = 2

type WorkersConfig struct {
	brokerHosts []string
	topic       string
	count       uint
	mod         int
	groupId     string
}

func RunKafkaWorkers(ctx context.Context, dao postgres.IWriterMarketplaceDAO, config WorkersConfig) {
	for range config.count {
		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers:          config.brokerHosts,
			Topic:            config.topic,
			GroupID:          config.groupId,
			StartOffset:      kafka.LastOffset,
			SessionTimeout:   time.Minute * 30,
			RebalanceTimeout: time.Second * 30,
		})

		go KafkaWorker(ctx, reader, dao, config.mod)
	}
}

func KafkaWorker(ctx context.Context, reader *kafka.Reader, dao postgres.IWriterMarketplaceDAO, mod int) {
	logger := logger_lib.CreateLoggerGateway(ctx, "KafkaWorker()")
	defer reader.Close()

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

		switch mod {
		case ModNewAnnouncement:
			var announcement models.ExtendedAnnouncement
			err = json.Unmarshal(msg.Value, &announcement)
			if err != nil {
				logger.WriteError(err.Error())

				err = reader.CommitMessages(ctx, msg)
				if err != nil {
					logger.WriteError(err.Error())
					time.Sleep(time.Second)
				}
				continue
			}

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
		case ModDeleteAnnouncement:
			id, err := strconv.ParseUint(string(msg.Value), 10, 64)
			if err != nil {
				logger.WriteError(err.Error())

				err = reader.CommitMessages(ctx, msg)
				if err != nil {
					logger.WriteError(err.Error())
					time.Sleep(time.Second)
				}
				continue
			}

			err = dao.DeleteAnnouncement(uint(id))
			if err != nil {
				logger.WriteError(err.Error())

				err = reader.CommitMessages(ctx, msg)
				if err != nil {
					logger.WriteError(err.Error())
					time.Sleep(time.Second)
				}
				continue
			}

		case ModUpdateAnnouncement:
			var updatedAnnouncement models.ExtendedAnnouncement
			err = json.Unmarshal(msg.Value, &updatedAnnouncement)
			if err != nil {
				logger.WriteError(err.Error())

				err = reader.CommitMessages(ctx, msg)
				if err != nil {
					logger.WriteError(err.Error())
					time.Sleep(time.Second)
				}
				continue
			}

			err = dao.UpdateAnnouncement(updatedAnnouncement)
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
