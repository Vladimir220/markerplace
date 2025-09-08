package remote

import (
	"context"
	"encoding/json"
	"fmt"
	"marketplace/env"
	"marketplace/models"

	"github.com/segmentio/kafka-go"
)

type IWriter interface {
	NewAnnouncement(announcement models.ExtendedAnnouncement) (err error)
	UpdateAnnouncement(updatedAnnouncement models.ExtendedAnnouncement) (err error)
	DeleteAnnouncement(announcementId uint) (err error)
}

func CreateKafkaWriter(ctx context.Context) (writer IWriter, err error) {
	data, err := env.GetKafkaEnvData()
	if err != nil {
		err = fmt.Errorf("%s %v", "CreateKafkaWriter():", err)
		return
	}

	writer = &KafkaWriter{
		data: data,
		ctx:  ctx,
	}
	return
}

type KafkaWriter struct {
	data env.KafkaEnvData
	ctx  context.Context
}

func (kw KafkaWriter) NewAnnouncement(announcement models.ExtendedAnnouncement) (err error) {
	logLabel := "KafkaWriter:NewAnnouncement():"
	writer := &kafka.Writer{
		Addr:         kafka.TCP(kw.data.BrokerHosts...),
		Topic:        kw.data.WriterData.NewAnnouncementTopicName,
		RequiredAcks: kafka.RequireOne,
		MaxAttempts:  10,
		Async:        false,
	}
	defer writer.Close()

	msg, err := json.Marshal(announcement)
	if err != nil {
		err = fmt.Errorf("%s %v", logLabel, err)
		return
	}

	err = writer.WriteMessages(kw.ctx, kafka.Message{
		Value: msg,
	})
	if err != nil {
		err = fmt.Errorf("%s %v", logLabel, err)
		return
	}
	return
}

func (kw KafkaWriter) UpdateAnnouncement(updatedAnnouncement models.ExtendedAnnouncement) (err error) {
	logLabel := "KafkaWriter:UpdateAnnouncement():"
	writer := &kafka.Writer{
		Addr:         kafka.TCP(kw.data.BrokerHosts...),
		Topic:        kw.data.WriterData.UpdateAnnouncementTopicName,
		RequiredAcks: kafka.RequireOne,
		MaxAttempts:  10,
		Async:        false,
	}
	defer writer.Close()

	msg, err := json.Marshal(updatedAnnouncement)
	if err != nil {
		err = fmt.Errorf("%s %v", logLabel, err)
		return
	}

	err = writer.WriteMessages(kw.ctx, kafka.Message{
		Value: msg,
	})
	if err != nil {
		err = fmt.Errorf("%s %v", logLabel, err)
		return
	}
	return
}

func (kw KafkaWriter) DeleteAnnouncement(announcementId uint) (err error) {
	logLabel := "KafkaWriter:DeleteAnnouncement():"
	writer := &kafka.Writer{
		Addr:         kafka.TCP(kw.data.BrokerHosts...),
		Topic:        kw.data.WriterData.DeleteAnnouncementTopicName,
		RequiredAcks: kafka.RequireOne,
		MaxAttempts:  10,
		Async:        false,
	}
	defer writer.Close()

	err = writer.WriteMessages(kw.ctx, kafka.Message{
		Value: []byte(fmt.Sprint(announcementId)),
	})
	if err != nil {
		err = fmt.Errorf("%s %v", logLabel, err)
		return
	}
	return
}
