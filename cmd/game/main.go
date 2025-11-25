package main

import (
	"log/slog"
	"os"
)

func main() {
	slog.Info("Started")
	if err := runApplication(); err != nil {
		slog.Error("Crashed", slog.String("error", err.Error()))
		os.Exit(1)
	}
	slog.Info("Stopped")
}
