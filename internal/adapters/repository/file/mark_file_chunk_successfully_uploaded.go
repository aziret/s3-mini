package file

import (
	"fmt"
	"github.com/aziret/s3-mini-internal/internal/lib/logger/sl"
	"log/slog"
)

func (repo *repository) MarkFileChunkSuccessfullyUploaded(UUID string) error {
	const op = "repository.file.MarkFileChunkSuccessfullyUploaded"
	log := repo.log.With(
		slog.String("op", op),
	)

	stmt, err := repo.db.Prepare("UPDATE file_chunks SET download_completed = TRUE WHERE id = $1")
	if err != nil {
		log.Error("failed to prepare statement", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(UUID)
	if err != nil {
		log.Error("failed to mark file chunk as uploaded", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("marked file chunk as uploaded", slog.String("UUID", UUID))

	return nil
}
