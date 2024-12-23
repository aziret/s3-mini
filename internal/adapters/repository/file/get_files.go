package file

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/aziret/s3-mini-internal/internal/lib/logger/sl"
	"github.com/aziret/s3-mini-internal/internal/model"
	"log/slog"
)

func (repo *repository) GetFiles(_ context.Context) (*[]model.FileInfo, error) {
	const op = "repository.file.GetFiles"

	log := repo.log.With(
		slog.String("op", op),
	)

	stmt, err := repo.db.Prepare("SELECT id, name FROM files")
	if err != nil {
		log.Error("failed to prepare statement", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Error("failed to close statement", sl.Err(err))
		}
	}(stmt)
	rows, err := stmt.Query()
	if err != nil {
		log.Error("failed to query files", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Error("failed to close row", sl.Err(err))
		}
	}()

	var files []model.FileInfo
	for rows.Next() {
		var file model.FileInfo

		if err := rows.Scan(&file.ID, &file.Name); err != nil {
			log.Error("failed to scan file", sl.Err(err))
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		files = append(files, file)
	}

	if err := rows.Err(); err != nil {
		log.Error("failed to iterate row", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &files, nil
}
