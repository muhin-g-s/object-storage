package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"object-storage/internal/config"
	hfw "object-storage/pkg/http"
)

type Handlers struct {
	config *config.Config
	logger *slog.Logger
}

func NewHandlers(config *config.Config, logger *slog.Logger) *Handlers {
	return &Handlers{config: config, logger: logger}
}

func (*Handlers) HandleUpload(ctx *hfw.Context) {
	key := ctx.Param("id")
	writer := ctx.Writer
	writer.WriteHeader(http.StatusOK)
	fmt.Fprintf(writer, "Объект %s успешно сохранен", key)
}

// func HandleUpload(w http.ResponseWriter, r *http.Request, storage *Storage) {
// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	key := r.URL.Path[UPLOAD_PREFIX_LEN:]

// 	data, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		http.Error(w, "Ошибка чтения данных", http.StatusInternalServerError)
// 		return
// 	}

// 	storage.Save(key, data)

// 	w.WriteHeader(http.StatusOK)
// 	fmt.Fprintf(w, "Объект %s успешно сохранен", key)
// }

// func HandleDownload(w http.ResponseWriter, r *http.Request, storage *Storage) {
// 	if r.Method != http.MethodGet {
// 		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	key := r.URL.Path[DOWNLOAD_PREFIX_LEN:]

// 	data, exists := storage.Load(key)
// 	if !exists {
// 		http.Error(w, "Объект не найден", http.StatusNotFound)
// 		return
// 	}

// 	w.WriteHeader(http.StatusOK)
// 	w.Write(data)
// }

// func HandleList(w http.ResponseWriter, r *http.Request, storage *Storage) {
// 	if r.Method != http.MethodGet {
// 		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	storage.mu.Lock()
// 	defer storage.mu.Unlock()

// 	keys := make([]string, 0, len(storage.files))
// 	for key := range storage.files {
// 		keys = append(keys, key)
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(keys)
// }
