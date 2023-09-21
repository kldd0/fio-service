package storage

import (
	"fmt"
)

type Storage interface {
}

var (
	ErrEntryAlreadyExists = fmt.Errorf("entry already exists")
	ErrEntryDoesntExists  = fmt.Errorf("entry doesn't exists")
)
