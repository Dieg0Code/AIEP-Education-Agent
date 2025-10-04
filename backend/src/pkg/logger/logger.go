package logger

import (
	"log/slog"
	"os"
)

// New crea una nueva instancia de logger con configuración JSON.
func New() *slog.Logger {
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true, // Incluye archivo y línea en los logs
	}))
}

// NewDevelopment crea un logger más legible para desarrollo.
func NewDevelopment() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))
}
