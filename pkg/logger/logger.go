package logger

import (
	"context"
	"log/slog"
	"os"
	"sync"
)

var (
	once   sync.Once
	logger *slog.Logger
)

func getLevel() slog.Level {
	if os.Getenv("DEBUG") != "" {
		return slog.LevelDebug
	}
	return slog.LevelInfo
}

func initLogger() {
	once.Do(func() {
		logger = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
			Level: getLevel(),
		}))
	})
}

func GetLogger(ctx context.Context) *slog.Logger {
	initLogger()
	return logger
}
