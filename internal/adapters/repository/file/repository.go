package file

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	repoPackage "github.com/aziret/s3-mini/internal/adapters/repository"
	"github.com/aziret/s3-mini/internal/lib/logger/sl"
	"github.com/lib/pq"

	"github.com/aziret/s3-mini/internal/model"
)

type repository struct {
	db  *sql.DB
	log *slog.Logger
}

func NewRepository(log *slog.Logger) (*repository, error) {
	const op = "repository.file.NewRepository"

	logger := log.With(
		slog.String("op", op),
	)

	db, err := sql.Open("postgres", connectionString())
	if err != nil {
		logger.Error("failed to open database connection", sl.Err(err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = db.Ping()
	if err != nil {
		logger.Error("Error pinging database: ", sl.Err(err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	logger.Info("Successfully connected to database")

	return &repository{
		db:  db,
		log: log,
	}, nil
}

func connectionString() string {
	var (
		pgUser    = os.Getenv("PG_USER")
		pgPass    = os.Getenv("PG_PASS")
		pgHost    = os.Getenv("PG_HOST")
		pgPort    = os.Getenv("PG_PORT")
		pgDB      = os.Getenv("PG_DB")
		pgSSLMode = os.Getenv("PG_SSL_MODE")
	)
	return fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s", pgUser, pgPass, pgDB, pgHost, pgPort, pgSSLMode)
}

func (repo *repository) Create(ctx context.Context, info *model.FileInfo) (int64, error) {
	const op = "repository.file.Create"

	log := repo.log.With(
		slog.String("op", op),
	)

	stmt, err := repo.db.Prepare("INSERT INTO files(name, file_path, upload_id, \"offset\", filetype) VALUES($1, $2, $3, $4, $5)")
	if err != nil {
		log.Error("failed to prepare statement", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(info.MetaData["filename"], info.Storage["Path"], info.ID, info.Offset, info.MetaData["filetype"])
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name() == "unique_violation" {
			return 0, fmt.Errorf("%s: %w", op, repoPackage.ErrUploadExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}
	_ = res

	return 0, nil
}
