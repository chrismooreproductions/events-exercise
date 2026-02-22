package logger

import (
	"log/slog"
	"os"
)

var LevelVar slog.LevelVar
var AddSourceVar = false

var Logger *slog.Logger

func init() {
	LevelVar.Set(slog.LevelInfo)

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: AddSourceVar,
		Level:     &LevelVar,
	})

	Logger = slog.New(handler)
}
