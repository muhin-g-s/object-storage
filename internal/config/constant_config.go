package config

func NewConstantConfig() *Config {
	return &Config{
		HTTPServer: HTTPServer{
			Address: "0.0.0.0:8081",
		},
		Storage: Storage{
			StorageDir: "./storage",
		},
		Env: EnvLocal,
	}
}
