package api

import (
	"log/slog"
	"object-storage/internal/config"
	hfw "object-storage/pkg/http"
)

type Api struct {
	Router *hfw.Router
	cfg    *config.Config
	logger *slog.Logger
}

func NewApi(cfg *config.Config, logger *slog.Logger) *Api {
	router := hfw.NewRouter()

	handlers := NewHandlers(cfg, logger)
	router.GET("/upload/:key", handlers.HandleUpload)

	return &Api{
		Router: router,
		cfg:    cfg,
		logger: logger,
	}
}
