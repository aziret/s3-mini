package file

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	repoPackage "github.com/aziret/s3-mini/internal/adapters/repository"
	"github.com/aziret/s3-mini/internal/lib/logger/sl"
	"github.com/mattn/go-sqlite3"

	"github.com/aziret/s3-mini/internal/model"
)

type repository struct {
	db  *sql.DB
	log *slog.Logger
}

func NewRepository(storagePath string, log *slog.Logger) (*repository, error) {
	const op = "repository.file.NewRepository"

	logger := log.With(
		slog.String("op", op),
	)

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		logger.Error("failed to open database connection", sl.Err(err))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &repository{
		db:  db,
		log: log,
	}, nil
}

func (repo *repository) Create(ctx context.Context, info *model.FileInfo) (int64, error) {
	const op = "repository.file.Create"

	log := repo.log.With(
		slog.String("op", op),
	)

	stmt, err := repo.db.Prepare("INSERT INTO files(name, file_path, upload_id, offset, filetype) VALUES(?, ?)")
	if err != nil {
		log.Error("failed to prepare statement", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(info.MetaData["filename"], info.Storage["Path"], info.ID, info.Offset, info.MetaData["filetype"])
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, repoPackage.ErrUploadExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}
	_ = res

	return 0, nil
}
