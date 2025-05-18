package config

type Config struct {
	HTTPServer
	Storage
}

type HTTPServer struct {
	Address string
}

type Storage struct {
	StorageDir string
}

type ConfigFactory interface {
	GetConfig() *Config
}
