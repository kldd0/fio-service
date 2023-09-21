package services

import (
	"github.com/IBM/sarama"
	"github.com/kldd0/fio-service/internal/storage"
)

type ServiceProvider struct {
	Db   storage.Storage
	Prod sarama.SyncProducer
	// FillService *
}
