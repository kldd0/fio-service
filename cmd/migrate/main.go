package main

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/kldd0/fio-service/internal/logs"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"
)

const (
	dbDriver   = "pgx"
	configFile = "config/pg.env"
)

var (
	flags = flag.NewFlagSet("migrate", flag.ExitOnError)
	dir   = flags.String("dir", "migrations", "directory with migration files")
)

func createDbString() string {
	type DbConfig struct {
		Db       string `env:"POSTGRES_DB"`
		User     string `env:"POSTGRES_USER"`
		Password string `env:"POSTGRES_PASSWORD"`
	}

	var cfg DbConfig

	err := cleanenv.ReadConfig(configFile, &cfg)
	if err != nil {
		panic(err)
	}

	dbString := "host=localhost user=" + cfg.User + " password=" + cfg.Password + " dbname=" + cfg.Db + " port=5432 sslmode=disable"
	return dbString
}

func main() {
	// setup logger
	logs.InitLogger(false)

	dbString := createDbString()

	flags.Parse(os.Args[1:])

	args := flags.Args()

	goose_command := args[1]

	db, err := goose.OpenDBWithDriver(dbDriver, dbString)
	if err != nil {
		logs.Logger.Fatal("Error: failed open db", zap.Error(err))
	}

	defer func() {
		if err := db.Close(); err != nil {
			logs.Logger.Fatal("Error: failed closing db", zap.Error(err))
		}
	}()

	if err := goose.Run(goose_command, db, *dir, args[1:]...); err != nil {
		logs.Logger.Fatal("Error: migrate "+goose_command, zap.Error(err))
	}
}
