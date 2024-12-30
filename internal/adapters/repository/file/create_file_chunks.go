package file

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"math"
	"os"
	"strconv"

	"github.com/aziret/s3-mini/internal/lib/logger/sl"
	"github.com/aziret/s3-mini/internal/model"
)

type databaseHandler interface {
	Exec(query string, args ...any) (sql.Result, error)
}

func (repo *repository) CreateFileChunksForFile(ctx context.Context, file *model.File) error {
	const op = "repository.file.CreateFileChunksForFile"

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

	var chunkSize int64

	chunkSizeVar := os.Getenv("CHUNK_SIZE")
	if chunkSizeVar == "" {
		chunkSize = 1_000_000
	} else {
		chunkSize, _ = strconv.ParseInt(chunkSizeVar, 10, 64)
	}
	chunksNumber := int(math.Ceil(float64(file.Size) / float64(chunkSize)))

	query := `
		INSERT INTO file_chunks (file_id, chunk_size, chunk_number)
        SELECT $1 AS file_id, $2 as chunk_size, gs.num AS chunk_number
		FROM generate_series(0, $3 - 1) AS gs(num)
        ON CONFLICT (file_id, chunk_number)
        DO UPDATE SET 
            chunk_size = EXCLUDED.chunk_size;
	`

	_, err := dbHandler.Exec(query, file.ID, chunkSize, chunksNumber)
	if err != nil {
		log.Error("failed inserting file chunk values", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
