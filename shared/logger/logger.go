package logger

import (
	"log/slog"
	"os"
)

var (
	Log *slog.Logger
)

func Init(debug bool) {
	var handler slog.Handler
	if debug {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: true,
		})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelInfo,
			AddSource: true,
		})
	}
	Log = slog.New(handler)
	Log.Info("Logger initialized", "debug", debug)
	slog.SetDefault(Log)
}
