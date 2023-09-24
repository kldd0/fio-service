package storage

import (
	"context"
	"fmt"

	"github.com/kldd0/fio-service/internal/model/domain_models"
)

type Storage interface {
	Get(ctx context.Context, filter string, target interface{}, limit, offset int) ([]domain_models.FioStruct, error)
	Save(ctx context.Context, fio_struct *domain_models.FioStruct) error
}

var (
	ErrEntryAlreadyExists = fmt.Errorf("entry already exists")
	ErrEntryDoesntExist   = fmt.Errorf("entry doesn't exists")
)
