package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

const configFile = "config/config.yaml"

type Config struct {
	DBUri string `yaml:"db_uri"`

	HTTPServer `yaml:"http_server"`

	Kafka `yaml:"kafka"`

	RedisUri  string `yaml:"redis_uri"`
	RedisPass string `yaml:"redis_pass"`
}

type HTTPServer struct {
	Address     string        `yaml:"address"`
	Timeout     time.Duration `yaml:"timeout"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

type Kafka struct {
	Brokers []string
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

func (s Service) DbUri() string {
	return s.config.DBUri
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

func (s Service) RedisUri() string {
	return s.config.RedisUri
}

func (s Service) RedisPass() string {
	return s.config.RedisPass
}

func (s Service) KafkaBrokers() []string {
	return s.config.Kafka.Brokers
}
