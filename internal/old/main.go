package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

const (
	STORAGE_DIR         = "./storage"
	UPLOAD_PREFIX_LEN   = len("/upload/")
	DOWNLOAD_PREFIX_LEN = len("/download/")
)

type Storage struct {
	mu    sync.Mutex
	files map[string][]byte
}

func NewStorage() *Storage {
	return &Storage{
		files: make(map[string][]byte),
	}
}

func (s *Storage) Save(key string, data []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.files[key] = data

	err := ioutil.WriteFile(STORAGE_DIR+"/"+key, data, 0644)
	if err != nil {
		log.Printf("Ошибка при сохранении файла %s: %v", key, err)
	}
}

func (s *Storage) Load(key string) ([]byte, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, exists := s.files[key]
	if exists {
		return data, true
	}

	data, err := ioutil.ReadFile(STORAGE_DIR + "/" + key)
	if err != nil {
		return nil, false
	}

	s.files[key] = data
	return data, true
}

func HandleUpload(w http.ResponseWriter, r *http.Request, storage *Storage) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	key := r.URL.Path[UPLOAD_PREFIX_LEN:]

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка чтения данных", http.StatusInternalServerError)
		return
	}

	storage.Save(key, data)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Объект %s успешно сохранен", key)
}

func HandleDownload(w http.ResponseWriter, r *http.Request, storage *Storage) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	key := r.URL.Path[DOWNLOAD_PREFIX_LEN:]

	data, exists := storage.Load(key)
	if !exists {
		http.Error(w, "Объект не найден", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func HandleList(w http.ResponseWriter, r *http.Request, storage *Storage) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	storage.mu.Lock()
	defer storage.mu.Unlock()

	keys := make([]string, 0, len(storage.files))
	for key := range storage.files {
		keys = append(keys, key)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(keys)
}

func main() {

	if _, err := os.Stat(STORAGE_DIR); os.IsNotExist(err) {
		err := os.Mkdir(STORAGE_DIR, 0755)
		if err != nil {
			log.Fatalf("Ошибка создания директории %s: %v", STORAGE_DIR, err)
		}
	}

	storage := NewStorage()

	http.HandleFunc("/upload/", func(w http.ResponseWriter, r *http.Request) {
		HandleUpload(w, r, storage)
	})
	http.HandleFunc("/download/", func(w http.ResponseWriter, r *http.Request) {
		HandleDownload(w, r, storage)
	})
	http.HandleFunc("/list", func(w http.ResponseWriter, r *http.Request) {
		HandleList(w, r, storage)
	})

	log.Println("Сервер запущен на порту 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
