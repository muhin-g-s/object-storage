package storage

import (
	"io/ioutil"
	"log/slog"
	"object-storage/internal/config"
	"sync"
)

type Storage struct {
	mu     sync.Mutex
	files  map[string][]byte
	logger *slog.Logger
	cfg    *config.Config
}

func NewStorage(logger *slog.Logger, cfg *config.Config) *Storage {
	return &Storage{
		files:  make(map[string][]byte),
		logger: logger,
		cfg:    cfg,
	}
}

func (s *Storage) Save(key string, data []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.files[key] = data

	err := ioutil.WriteFile(s.cfg.StorageDir+"/"+key, data, 0644)
	if err != nil {
		s.logger.Error("Ошибка при сохранении файла %s: %v", key, err)
	}
}

func (s *Storage) Load(key string) ([]byte, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, exists := s.files[key]
	if exists {
		return data, true
	}

	data, err := ioutil.ReadFile(s.cfg.StorageDir + "/" + key)
	if err != nil {
		return nil, false
	}

	s.files[key] = data
	return data, true
}

func (s *Storage) List() []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	keys := make([]string, 0, len(s.files))
	for key := range s.files {
		keys = append(keys, key)
	}

	return keys
}
