package log

import (
	"auth_service/env"
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/segmentio/kafka-go"
)

// For LoggerAdapter
func CreateKafkaLogger(ctx context.Context, parentName string) (logger ILoggerRemote, err error) {
	data, err := env.GetKafkaEnvData()

	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		err = errors.New("CreateKafkaLogger(): env variable SERVICE_NAME expected")
		return
	}

	logger = &KafkaLogger{
		data:        data,
		parentName:  parentName,
		serviceName: serviceName,
		ctx:         ctx,
	}
	return
}

type KafkaLogger struct {
	serviceName string
	parentName  string
	data        env.KafkaEnvData
	ctx         context.Context
}

func (l KafkaLogger) WriteWarning(msg string) (err error) {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(l.data.BrokerHosts...),
		Topic:        l.data.WarningTopicName,
		RequiredAcks: kafka.RequireOne,
		MaxAttempts:  10,
		Async:        false,
	}
	defer writer.Close()

	err = writer.WriteMessages(l.ctx, kafka.Message{
		Value: []byte(fmt.Sprintf("%s:%s: %s", l.serviceName, l.parentName, msg)),
	})
	return
}

func (l KafkaLogger) WriteError(msg string) (err error) {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(l.data.BrokerHosts...),
		Topic:        l.data.ErrorTopicName,
		RequiredAcks: kafka.RequireOne,
		MaxAttempts:  10,
		Async:        false,
	}
	defer writer.Close()

	err = writer.WriteMessages(l.ctx, kafka.Message{
		Value: []byte(fmt.Sprintf("%s:%s: %s", l.serviceName, l.parentName, msg)),
	})
	return
}

func (l KafkaLogger) WriteInfo(msg string) (err error) {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(l.data.BrokerHosts...),
		Topic:        l.data.InfoTopicName,
		RequiredAcks: kafka.RequireOne,
		MaxAttempts:  10,
		Async:        false,
	}
	defer writer.Close()

	err = writer.WriteMessages(l.ctx, kafka.Message{
		Value: []byte(fmt.Sprintf("%s:%s: %s", l.serviceName, l.parentName, msg)),
	})
	return
}
