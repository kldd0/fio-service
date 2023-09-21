package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	timeout = time.Second * 15
)

type configGetter interface {
	RedisUri() string
	RedisPass() string
}

type Client struct {
	rdb *redis.Client
}

func New(ctx context.Context, config configGetter) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         config.RedisUri(),
		Password:     "",
		ReadTimeout:  timeout,
		WriteTimeout: timeout,
		DB:           0, // use default DB
	})

	return &Client{rdb: rdb}, rdb.Ping(ctx).Err()
}
