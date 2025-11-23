package logger

import (
	"log/slog"
	"os"
)

var Log *slog.Logger

func init() {
	lofgile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		slog.Error("Ошибка открытия файла лога", "Ошибка", err)
		os.Exit(1)
	}

	handler := slog.NewJSONHandler(lofgile, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})
	Log = slog.New(handler)
}

func LogError(err error, message string) {
	if err != nil {
		Log.Error(message, slog.Any("Ошибка", err))
	}
}
