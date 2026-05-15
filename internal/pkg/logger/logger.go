package logger

import (
	"log/slog"
	"os"
	"sync"
)

var (
	once sync.Once
	log *slog.Logger
)

func Init(){
	once.Do(func() {
		var handler slog.Handler
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})

		log = slog.New(handler)
		slog.SetDefault(log)

	})
}

func Get() *slog.Logger {
	return log
}