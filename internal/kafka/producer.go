package kafka

import (
	"fmt"
	"time"

	"github.com/IBM/sarama"
)

var Producer sarama.SyncProducer

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
