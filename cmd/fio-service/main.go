package main

import (
	"context"
	"encoding/json"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/kldd0/fio-service/internal/clients/redis"
	"github.com/kldd0/fio-service/internal/config"
	"github.com/kldd0/fio-service/internal/kafka"
	"github.com/kldd0/fio-service/internal/logs"
	"github.com/kldd0/fio-service/internal/model/api"
	"github.com/kldd0/fio-service/internal/services"
	"github.com/kldd0/fio-service/internal/storage/postgres"
)

var (
	develMode = flag.Bool("devel", false, "development mode")
)

func main() {
	flag.Parse()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// setup logger
	logs.InitLogger(*develMode)

	// setup config
	config, err := config.New()
	if err != nil {
		logs.Logger.Fatal("Error: config init failed:", zap.Error(err))
	}

	// setup db
	db, err := postgres.New(config.DbUri())
	if err != nil {
		logs.Logger.Fatal("Error: failed connecting to database:", zap.Error(err))
	}
	defer db.Close()

	if err := db.InitDB(ctx); err != nil {
		logs.Logger.Fatal("Error: storage init failed:", zap.Error(err))
	}

	// setup redis
	_, err = redis.New(ctx, config)
	if err != nil {
		logs.Logger.Fatal("Error: redis init failed:", zap.Error(err))
	}

	// setup producer for responding
	err = kafka.NewSyncProducer(config.KafkaBrokers())
	if err != nil {
		logs.Logger.Fatal("Error: sync producer init failed", zap.Error(err))
	}

	provider := services.ServiceProvider{
		Db:          db,
		Prod:        kafka.Producer,
		APIServices: api.FioAPIClient{},
	}

	// setup consumer group
	go func() {
		err = kafka.StartConsumerGroup(ctx, config.KafkaBrokers(), provider)
		if err != nil {
			logs.Logger.Fatal("Error: consumer group is failed:", zap.Error(err))
		}
	}()

	// http router
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// healthcheck route
	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		w.WriteHeader(http.StatusOK)

		_ = json.NewEncoder(w).Encode(map[string]bool{
			"pong": true,
		})
	})

	// server configuration
	srv := &http.Server{
		Addr:         config.HTTPAddr(),
		Handler:      router,
		ReadTimeout:  config.Timeout(),
		WriteTimeout: config.Timeout(),
		IdleTimeout:  config.IdleTimeout(),
	}

	// listen to OS signals and gracefully shutdown HTTP server
	done := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint

		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		logs.Info("Stopping server")
		if err := srv.Shutdown(ctx); err != nil {
			logs.Info("HTTP Server Shutdown Error:", zap.Error(err))
		}
		close(done)
	}()

	logs.Info("Info: Starting HTTP server on " + config.HTTPAddr())

	// start http server
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		logs.Logger.Fatal("Error: HTTP server ListenAndServe error:", zap.Error(err))
	}

	<-done
}
