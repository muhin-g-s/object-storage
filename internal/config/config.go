package config

type EnvType string

const (
	EnvLocal EnvType = "local"
	EnvDev   EnvType = "dev"
	EnvProd  EnvType = "prod"
)

type Config struct {
	HTTPServer
	Storage
	Env EnvType
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
