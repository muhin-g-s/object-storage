package config

func NewConstantConfig() *Config {
	return &Config{
		HTTPServer: HTTPServer{
			Address: "0.0.0.0:8080",
		},
		Storage: Storage{
			StorageDir: "./storage",
		},
	}
}
