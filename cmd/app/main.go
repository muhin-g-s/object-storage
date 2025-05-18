package main

import (
	"context"
	"net/http"
	"object-storage/internal/config"
	"object-storage/internal/infrastructure/api"
	"object-storage/internal/infrastructure/logger"
	"object-storage/internal/services/storage"
	"object-storage/pkg/logger/sl"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.NewConstantConfig()
	logger := logger.SetupLogger(cfg.Env)

	logger.Info("starting server", "address", cfg.HTTPServer.Address)

	storage := storage.NewStorage(logger, cfg)

	storage.SetupStorage()

	api := api.NewApi(cfg, logger, storage)

	srv := &http.Server{
		Addr:    cfg.Address,
		Handler: api.Router,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logger.Error("cannot start server", sl.Err(err))
		}
	}()

	logger.Info("server started")

	<-done
	logger.Info("stopping server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("failed to stop server", sl.Err(err))

		return
	}

	storage.Shutdown()

	logger.Info("server stopped")
}
