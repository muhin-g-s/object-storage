package logger

import (
	"log/slog"
	"object-storage/internal/config"
	"object-storage/pkg/logger/handlers/slogpretty"
	"os"
)

func SetupLogger(env config.EnvType) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case config.EnvLocal:
		logger = setupPrettySlog()
	case config.EnvDev:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case config.EnvProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return logger
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
