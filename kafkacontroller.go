package postgresql_generator

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Shopify/sarama"
)

// UploadEvent structure of kafka
type UploadEvent struct {
	Topic     string
	Payload   string
	Timestamp string
}

// Producer of kafka
type Producer struct {
	sarama.SyncProducer
}

func (e *UploadEvent) buildEvent() (*sarama.ProducerMessage, error) {
	if e.Payload == "" {
		return nil, errors.New("Missing Payload")
	}
	msg := &sarama.ProducerMessage{Topic: e.Topic, Key: sarama.StringEncoder(e.Payload), Value: sarama.StringEncoder(e.Payload)}
	return msg, nil
}

// NewProducer create kafka producer
func NewProducer() (*Producer, error) {
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if kafkaBrokers == "" {
		return nil, errors.New("Missing KAFKA_BROKERS env")
	}
	sp, err := sarama.NewSyncProducer(strings.Split(kafkaBrokers, ","), nil)
	if err != nil {
		return nil, err
	}
	return &Producer{sp}, nil
}

// Produce kafka event
func (p *Producer) Produce(e *UploadEvent) error {
	var partition int32
	var offset int64
	msg, err := e.buildEvent()
	if err != nil {
		return err
	}
	partition, offset, err = p.SyncProducer.SendMessage(msg)
	if err == nil {
		log.Printf("Produced: %s metadata: {partition:%d, offset:%d}\n", fmt.Sprint(e), partition, offset)
	}
	return err
}
