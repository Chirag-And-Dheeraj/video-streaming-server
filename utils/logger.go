package utils

import (
	"log/slog"
	"os"
	"video-streaming-server/config"
)

var (
	Logger *slog.Logger
)

func InitLogger() {
	var handler slog.Handler
	if config.AppConfig.Debug == "true" {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	}
	Logger = slog.New(handler)
	Logger.Info("Logger initialized", "debug", config.AppConfig.Debug)
	slog.SetDefault(Logger)
}
