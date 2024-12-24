package file

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/aziret/s3-mini/internal/lib/logger/sl"
)

func (repo *repository) MarkFileChunkSuccessfullyUploaded(UUID string, serverID string) error {
	const op = "repository.file.MarkFileChunkSuccessfullyUploaded"
	log := repo.log.With(
		slog.String("op", op),
	)

	stmt, err := repo.db.Prepare("UPDATE file_chunks SET download_completed = TRUE, server_id = $1 WHERE id = $2")
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

	_, err = stmt.Exec(serverID, UUID)
	if err != nil {
		log.Error("failed to mark file chunk as uploaded", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("marked file chunk as uploaded", slog.String("UUID", UUID))

	return nil
}
