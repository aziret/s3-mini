package file

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"

	"github.com/aziret/s3-mini/internal/lib/logger/sl"
)

func (repo *repository) MarkFilesAsReadyToDownload(ctx context.Context, deletedFileIds []int64) error {
	const op = "repository.file.MarkFilesAsReadyToDownload"
	log := repo.log.With(
		slog.String("op", op),
	)

	deletedFilesPlaceholders := make([]string, 0, len(deletedFileIds))
	args := make([]interface{}, len(deletedFileIds))
	for i, fileID := range deletedFileIds {
		deletedFilesPlaceholders = append(deletedFilesPlaceholders, fmt.Sprintf("$%d", i+1))
		args[i] = fileID
	}

	stmt, err := repo.db.Prepare("UPDATE files SET ready_to_download = TRUE WHERE id IN (" + strings.Join(deletedFilesPlaceholders, ",") + ")")
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

	_, err = stmt.Exec(args...)
	if err != nil {
		log.Error("failed to mark file as ready to download", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
