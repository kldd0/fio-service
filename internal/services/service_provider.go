package services

import "github.com/kldd0/fio-service/internal/storage"

type ServiceProvider struct {
	Db storage.Storage
	// FillService *
}
