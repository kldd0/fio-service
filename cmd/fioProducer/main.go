package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/kldd0/fio-service/internal/logs"
	message "github.com/kldd0/fio-service/internal/model/domain_models"
	"go.uber.org/zap"
)

var (
	develMode = flag.Bool("devel", false, "development mode")
)

var Producer sarama.SyncProducer

func main() {
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// setup logger
	logs.InitLogger(*develMode)

	KafkaBrokers := []string{"localhost:29092"}
	err := NewSyncProducer(KafkaBrokers)
	if err != nil {
		logs.Logger.Fatal("Error: sync producer init failed", zap.Error(err))
	}

	msg := message.Message{
		Name:       "Dmitriy",
		Surname:    "Ushakov",
		Patronymic: "Vasilevich",
	}

	result, _ := json.Marshal(msg)

	partition, offset, err := Producer.SendMessage(
		&sarama.ProducerMessage{
			Topic: "fio-topic",
			Value: sarama.ByteEncoder(result),
		})

	if err != nil {
		logs.Logger.Fatal("Error: marshaling error", zap.Error(err))
	}
	logs.Logger.Info("Sent message", zap.Int32("partition", partition), zap.Int64("offset", offset))

	<-ctx.Done()
}

func NewSyncProducer(brokerList []string) error {
	const op = "producer.NewSyncProducer"

	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 3
	config.Producer.Retry.Backoff = time.Millisecond * 250

	if config.Producer.Idempotent {
		config.Producer.Retry.Max = 1
		config.Net.MaxOpenRequests = 1
	}
	config.Producer.Return.Successes = true
	_ = config.Producer.Partitioner

	producer, err := sarama.NewSyncProducer(brokerList, config)

	if err != nil {
		return fmt.Errorf("%s: starting Sarama producer: %w", op, err)
	}

	Producer = producer

	return nil
}
