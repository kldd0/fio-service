package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/kldd0/fio-service/internal/model/domain_models"
	"github.com/kldd0/fio-service/internal/storage"

	"github.com/jmoiron/sqlx"
)

const dbDriver = "pgx"

const initSchema = `
CREATE TABLE IF NOT EXISTS fio_table (
	id SERIAL PRIMARY KEY,
	first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    patronymic VARCHAR(50)
);
`

type Storage struct {
	db *sqlx.DB
}

func New(dbUri string) (*Storage, error) {
	const op = "storage.postgres.New"

	db, err := sqlx.Open(dbDriver, dbUri)
	if err != nil {
		return nil, fmt.Errorf("%s: open db connection: %w", op, err)
	}

	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) InitDB(ctx context.Context) error {
	const op = "storage.postgres.InitDB"

	_, err := s.db.ExecContext(ctx, initSchema)
	if err != nil {
		return fmt.Errorf("%s: creating table: %w", op, err)
	}

	return nil
}

func (s *Storage) Save(ctx context.Context, fio_struct *domain_models.FioStruct) error {
	const op = "storage.postgres.Save"

	q := `INSERT INTO fio_table (name, surname, patronymic, age, gender, nationality) VALUES ($1, $2, $3, $4, $5, $6)`

	stmt, err := s.db.PrepareContext(ctx, q)
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	if _, err := stmt.ExecContext(ctx, fio_struct.Name, fio_struct.Surname, fio_struct.Patronymic, fio_struct.Age, fio_struct.Gender, fio_struct.Nationality); err != nil {
		return fmt.Errorf("%s: saving entry: %w", op, err)
	}

	return nil
}

func (s *Storage) Get(ctx context.Context, name, surname string) ([]domain_models.FioStruct, error) {
	const op = "storage.postgres.Get"

	q := `SELECT * FROM fio_table WHERE name = $1 AND surname = $2`

	stmt, err := s.db.PrepareContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	var result []domain_models.FioStruct

	err = stmt.QueryRowContext(ctx, name, surname).Scan(&result)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrEntryDoesntExists
		}

		return nil, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return result, nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}
