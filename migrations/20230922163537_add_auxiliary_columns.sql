-- +goose Up
-- +goose StatementBegin
ALTER TABLE fio_table
ADD COLUMN age INTEGER NOT NULL,
ADD COLUMN gender VARCHAR(10) NOT NULL,
ADD COLUMN nationality VARCHAR(10) NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd