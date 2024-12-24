package file

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/aziret/s3-mini-internal/internal/lib/logger/sl"
	"github.com/aziret/s3-mini-internal/internal/model"
)

func (repo *repository) GetUploadCompletedFiles(ctx context.Context) (*[]model.File, error) {
	const op = "repository.file.GetUploadCompletedFiles"

	log := repo.log.With(
		slog.String("op", op),
	)

	stmt, err := repo.db.Prepare(`SELECT id, file_path FROM files WHERE download_completed = TRUE AND ready_to_download = FALSE`)
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
		log.Error("failed to execute query", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Error("failed to close row", sl.Err(err))
		}
	}(rows)

	files := make([]model.File, 0)

	for rows.Next() {
		var file model.File
		err := rows.Scan(&file.ID, &file.FilePath)
		if err != nil {
			log.Error("failed to scan file", sl.Err(err))
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		files = append(files, file)
	}

	if err := rows.Err(); err != nil {
		log.Error("failed to fetch files", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &files, nil
}
