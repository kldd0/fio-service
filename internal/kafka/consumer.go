package consumer

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/kldd0/fio-service/internal/logs"
	"github.com/kldd0/fio-service/internal/services"
	"go.uber.org/zap"
)

var (
	KafkaTopic         = "fio-topic"
	KafkaConsumerGroup = "fio-consumer-group"
	BrokersList        = []string{"kafka:9092"}
	Assignor           = "range"
)

type Consumer struct {
	// services services.ServiceProvider
}

func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	fmt.Println("consumer - setup")
	return nil
}

func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	fmt.Println("consumer - cleanup")
	return nil
}

func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	_ = context.Background()
	for message := range claim.Messages() {
		logs.Logger.Info("Message received", zap.String("msg", string(message.Value)))

		session.MarkMessage(message, "")
	}

	return nil
}

func StartConsumerGroup(ctx context.Context, brokerList []string, services services.ServiceProvider) error {
	const op = "kafka.consumer.StartConsumerGroup"

	consumerGroupHandler := Consumer{}

	config := sarama.NewConfig()
	config.Version = sarama.V2_5_0_0
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	switch Assignor {
	case "sticky":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategySticky()}
	case "round-robin":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	case "range":
		config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRange()}
	default:
		log.Panicf("Unrecognized consumer group partition assignor: %s", Assignor)
	}

	// Create consumer group
	consumerGroup, err := sarama.NewConsumerGroup(brokerList, KafkaConsumerGroup, config)
	if err != nil {
		return fmt.Errorf("%s: starting consumer group: %w", op, err)
	}

	err = consumerGroup.Consume(ctx, []string{KafkaTopic}, &consumerGroupHandler)
	if err != nil {
		return fmt.Errorf("%s: consuming via handler: %w", op, err)
	}
	return nil
}
