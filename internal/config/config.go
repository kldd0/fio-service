package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

const configFile = "data/config.yaml"

type Config struct {
	Env        string `yaml:"env"`
	StorageDSN string `yaml:"storage_dsn"`
	HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

type Service struct {
	config Config
}

func New() (*Service, error) {
	const op = "config.New"

	var s *Service = &Service{}

	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("%s: reading file: %w", op, err)
	}

	err = yaml.Unmarshal(data, &s.config)
	if err != nil {
		return nil, fmt.Errorf("%s: unmarshaling yaml: %w", op, err)
	}

	return s, nil
}

func (s Service) DSN() string {
	return s.config.StorageDSN
}

func (s Service) HTTPAddr() string {
	return s.config.HTTPServer.Address
}

func (s Service) Timeout() time.Duration {
	return s.config.HTTPServer.Timeout
}

func (s Service) IdleTimeout() time.Duration {
	return s.config.HTTPServer.IdleTimeout
}
