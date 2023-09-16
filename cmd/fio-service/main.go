package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"

	"github.com/kldd0/fio-service/internal/config"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// setup logger [dev] -- debug
	log := log.Default()

	// config initialization
	config, err := config.New()
	if err != nil {
		log.Fatal("Error: failed initializing config: ", err)
	}

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

		fmt.Fprint(w, "Pong")
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

		log.Print("Stopping server")
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("HTTP Server Shutdown Error: %v", err)
		}
		close(done)
	}()

	log.Printf("Starting HTTP server on: %s", config.HTTPAddr())

	// start http server
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal("Error: HTTP server ListenAndServe error: ", err)
	}

	<-done
}
