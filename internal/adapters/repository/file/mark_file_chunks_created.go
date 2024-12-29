package file

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aziret/s3-mini/internal/lib/logger/sl"
)

func (repo *repository) MarkFileChunksCreated(ctx context.Context, id int64) error {
	const op = "repository.file.MarkFileChunksCreated"

	log := repo.log.With(
		slog.String("op", op),
	)

	var dbHandler databaseHandler

	tx := extractTx(ctx)
	if tx != nil {
		dbHandler = tx
	} else {
		dbHandler = repo.db
	}

	query := "UPDATE files SET file_chunks_created = TRUE WHERE id = $1"

	_, err := dbHandler.Exec(query, id)

	if err != nil {
		log.Error("failed to mark file chunks created", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Info("successfully marked file chunks created", slog.Attr{Key: "ID", Value: slog.Int64Value(id)})
	return nil
}
