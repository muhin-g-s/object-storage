package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log/slog"
	"net/http"
	"object-storage/internal/config"
	hfw "object-storage/pkg/http"
)

type Handlers struct {
	config  *config.Config
	logger  *slog.Logger
	storage Storage
}

func NewHandlers(config *config.Config, logger *slog.Logger, storage Storage) *Handlers {
	return &Handlers{
		config:  config,
		logger:  logger,
		storage: storage,
	}
}

func (h *Handlers) HandleUpload(ctx *hfw.Context) {
	writer := ctx.Writer

	key := ctx.Param("key")

	if key == "" {
		http.Error(writer, "Неверное имя объекта", http.StatusBadRequest)
		return
	}

	data, err := ioutil.ReadAll(ctx.Request.Body)

	if err != nil {
		http.Error(writer, "Ошибка чтения данных", http.StatusInternalServerError)
		return
	}

	reqCtx := ctx.Request.Context()

	h.storage.Save(reqCtx, key, data)

	writer.WriteHeader(http.StatusOK)
	fmt.Fprintf(writer, "Объект %s успешно сохранен", key)
}

func (h *Handlers) HandleDownload(ctx *hfw.Context) {
	writer := ctx.Writer

	key := ctx.Param("key")

	if key == "" {
		http.Error(writer, "Неверное имя объекта", http.StatusBadRequest)
		return
	}

	reqCtx := ctx.Request.Context()

	data, exists := h.storage.Load(reqCtx, key)

	if !exists {
		http.Error(writer, "Объект не найден", http.StatusNotFound)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write(data)
}

func (h *Handlers) HandleList(ctx *hfw.Context) {
	writer := ctx.Writer
	reqCtx := ctx.Request.Context()

	keys := h.storage.List(reqCtx)

	writer.Header().Set("Content-Type", "application/json")
	json.NewEncoder(writer).Encode(keys)
}
