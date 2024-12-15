package main

import (
	"github.com/aziret/s3-mini/internal/adapters/repository/file"
	"github.com/aziret/s3-mini/internal/lib/logger/sl"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	var (
		env         = os.Getenv("ENV")
		storagePath = os.Getenv("STORAGE_PATH")
		log         = setupLogger(env)
	)

	log.Debug("debug messages enables")

	repo, err := file.NewRepository(storagePath)

	if err != nil {
		log.Error("failed to initialize storage", sl.Err(err))
		os.Exit(1)
	}

	_ = repo
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelDebug},
			),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelDebug},
			),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelInfo},
			),
		)
	}

	return log
}

func init() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
}
