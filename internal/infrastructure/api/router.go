package api

import (
	"context"
	"log/slog"
	"object-storage/internal/config"
	hfw "object-storage/pkg/http"
)

type Storage interface {
	Save(ctx context.Context, key string, data []byte) error
	Load(ctx context.Context, key string) ([]byte, bool)
	List(ctx context.Context) []string
}

type Api struct {
	Router *hfw.Router
}

func NewApi(cfg *config.Config, logger *slog.Logger, storage Storage) *Api {
	router := hfw.NewRouter()

	handlers := NewHandlers(cfg, logger, storage)

	router.GET("/download/:key", handlers.HandleDownload)
	router.POST("/upload/:key", handlers.HandleUpload)
	router.GET("/list", handlers.HandleList)

	return &Api{
		Router: router,
	}
}
