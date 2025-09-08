package main

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

var ModWarning = 0
var ModError = 1
var ModInfo = 2

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
	logger := CreateLogger()
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
		case ModWarning:
			logger.WriteWarning(string(msg.Value))
		case ModInfo:
			logger.WriteInfo(string(msg.Value))
		case ModError:
			logger.WriteError(string(msg.Value))
		}

		err = reader.CommitMessages(context.Background(), msg)
		if err != nil {
			logger.WriteError(err.Error())
			time.Sleep(time.Second)
			continue
		}
	}
}
