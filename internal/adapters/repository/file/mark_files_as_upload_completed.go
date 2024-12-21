package file

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/aziret/s3-mini-internal/internal/lib/logger/sl"
)

func (repo *repository) MarkFilesAsUploadCompleted(_ context.Context) error {
	const op = "repository.file.GetFilesWithUploads"

	log := repo.log.With(
		slog.String("op", op),
	)

	stmt, err := repo.db.Prepare(`
        UPDATE files
        SET download_completed = TRUE
        WHERE NOT EXISTS (
            SELECT 1
            FROM file_chunks fc
            WHERE fc.file_id = files.id
            AND fc.download_completed = FALSE
        )
        AND EXISTS (
            SELECT 1
            FROM file_chunks fc
            WHERE fc.file_id = files.id
        );
    `)

	if err != nil {
		log.Error("failed to prepare statement", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Error("failed to close statement", sl.Err(err))
		}
	}(stmt)

	_, err = stmt.Exec()
	if err != nil {
		log.Error("failed to mark file as uploaded", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("marked files as uploaded")
	return nil
}
