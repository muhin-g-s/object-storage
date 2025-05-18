package storage

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"log/slog"
	"object-storage/internal/config"
	"object-storage/pkg/logger/sl"
	"os"
	"path/filepath"
	"sync"
)

var (
	ErrInvalidKey      = errors.New("invalid key")
	ErrStorageNotReady = errors.New("storage directory not ready")
	ErrDataCorrupted   = errors.New("data integrity check failed")
	ErrCompression     = errors.New("compression/decompression failed")
)

type writeRequest struct {
	key  string
	data []byte
	hash string
}

type Storage struct {
	mu         sync.RWMutex
	files      map[string][]byte
	logger     *slog.Logger
	cfg        *config.Config
	ready      bool
	writeQueue chan writeRequest
	shutdown   chan struct{}
}

func NewStorage(logger *slog.Logger, cfg *config.Config) *Storage {
	return &Storage{
		files:      make(map[string][]byte),
		logger:     logger,
		cfg:        cfg,
		writeQueue: make(chan writeRequest, 100),
		shutdown:   make(chan struct{}),
	}
}

func (s *Storage) Save(ctx context.Context, key string, data []byte) error {
	if !s.ready {
		return ErrStorageNotReady
	}

	if err := s.validateKey(key); err != nil {
		return err
	}

	compressedData, err := s.compress(data)
	if err != nil {
		s.logger.Error("Failed to compress data",
			slog.String("key", key),
			slog.Any("error", sl.Err(err)))
		return ErrCompression
	}

	hash := s.computeHash(compressedData)

	s.mu.Lock()
	s.files[key] = append([]byte(nil), compressedData...)
	s.mu.Unlock()

	select {
	case s.writeQueue <- writeRequest{key: key, data: compressedData, hash: hash}:
		return nil
	case <-ctx.Done():
		s.mu.Lock()
		delete(s.files, key)
		s.mu.Unlock()
		return ctx.Err()
	}
}

func (s *Storage) Load(ctx context.Context, key string) ([]byte, bool) {
	if !s.ready {
		return nil, false
	}
	if err := s.validateKey(key); err != nil {
		s.logger.Warn("Invalid key", slog.String("key", key))
		return nil, false
	}

	s.mu.RLock()
	compressedData, exists := s.files[key]
	s.mu.RUnlock()

	if exists {
		data, err := s.decompress(compressedData)
		if err != nil {
			s.logger.Error("Failed to decompress data",
				slog.String("key", key),
				slog.Any("error", sl.Err(err)))
			return nil, false
		}
		return data, true
	}

	compressedData, err := os.ReadFile(filepath.Join(s.cfg.StorageDir, key))
	if err != nil {
		if !os.IsNotExist(err) {
			s.logger.Warn("Failed to read file from disk",
				slog.String("key", key),
				slog.Any("error", sl.Err(err)))
		}
		return nil, false
	}

	hash := s.computeHash(compressedData)
	if !s.verifyFileHash(key, hash) {
		s.logger.Error("Data integrity check failed",
			slog.String("key", key))
		return nil, false
	}

	data, err := s.decompress(compressedData)
	if err != nil {
		s.logger.Error("Failed to decompress data",
			slog.String("key", key),
			slog.Any("error", sl.Err(err)))
		return nil, false
	}

	s.mu.Lock()
	s.files[key] = append([]byte(nil), compressedData...)
	s.mu.Unlock()

	return data, true
}

func (s *Storage) List(ctx context.Context) []string {
	if !s.ready {
		return nil
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	keys := make([]string, 0, len(s.files))
	for key := range s.files {
		keys = append(keys, key)
	}
	return keys
}

func (s *Storage) Shutdown() {
	close(s.shutdown)
}

func (s *Storage) SetupStorage() {
	if s.ready {
		return
	}

	s.createDir()
	s.setFileInMap()

	s.ready = true

	go s.processWriteQueue()
}

func (s *Storage) createDir() {
	storageDir := s.cfg.StorageDir

	if _, err := os.Stat(storageDir); os.IsNotExist(err) {
		err := os.Mkdir(storageDir, 0755)
		if err != nil {
			s.logger.Error("Ошибка создания директории %s: %v", storageDir, sl.Err(err))
		}
	}
}

func (s *Storage) setFileInMap() {
	storageDir := s.cfg.StorageDir

	entries, err := os.ReadDir(storageDir)
	if err != nil {
		s.logger.Error("Failed to read storage directory",
			slog.String("dir", storageDir),
			slog.Any("error", sl.Err(err)))
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			data, err := os.ReadFile(filepath.Join(storageDir, entry.Name()))
			if err != nil {
				s.logger.Warn("Failed to load file",
					slog.String("file", entry.Name()),
					slog.Any("error", sl.Err(err)))
				continue
			}
			s.mu.Lock()
			s.files[entry.Name()] = data
			s.mu.Unlock()
		}
	}
}

func (s *Storage) processWriteQueue() {
	for {
		select {
		case req := <-s.writeQueue:
			if err := s.writeToDisk(req.key, req.data); err != nil {
				s.logger.Error("Failed to write to disk",
					slog.String("key", req.key),
					slog.Any("error", sl.Err(err)))
				s.mu.Lock()
				delete(s.files, req.key)
				s.mu.Unlock()
			}
		case <-s.shutdown:
			return
		}
	}
}

func (s *Storage) writeToDisk(key string, data []byte) error {
	path := filepath.Join(s.cfg.StorageDir, key)

	if err := os.WriteFile(path, data, 0644); err != nil {
		return err
	}

	return nil
}

func (s *Storage) validateKey(key string) error {
	if key == "" || len(key) > 1024 || filepath.Base(key) != key {
		return ErrInvalidKey
	}
	return nil
}

func (s *Storage) compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	gz, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		return nil, err
	}
	if _, err := gz.Write(data); err != nil {
		gz.Close()
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *Storage) decompress(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return io.ReadAll(reader)
}

func (s *Storage) computeHash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func (s *Storage) verifyFileHash(path, expectedHash string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	return s.computeHash(data) == expectedHash
}
