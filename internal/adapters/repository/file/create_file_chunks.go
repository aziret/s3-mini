package file

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/aziret/s3-mini/internal/lib/logger/sl"
	"github.com/aziret/s3-mini/internal/model"
)

func (repo *repository) CreateFileChunksForFile(_ context.Context, file *model.File) error {
	const op = "repository.file.CreateFileChunksForFile"

	log := repo.log.With(
		slog.String("op", op),
	)

	var chunkSize int64

	chunkSizeVar := os.Getenv("CHUNK_SIZE")
	if chunkSizeVar == "" {
		chunkSize = 1_000_000
	} else {
		chunkSize, _ = strconv.ParseInt(chunkSizeVar, 10, 64)
	}
	chunksNumber := int(math.Ceil(float64(file.Size) / float64(chunkSize)))

	query := fmt.Sprintf(`
		INSERT INTO file_chunks (file_id, chunk_size, chunk_number)
        VALUES %s
        ON CONFLICT (file_id, chunk_number)
        DO UPDATE SET 
            chunk_size = EXCLUDED.chunk_size
	`, valuesPlaceholder(chunksNumber))

	_, err := repo.db.Exec(query, *fileChunkValues(file.ID, chunkSize, chunksNumber)...)
	if err != nil {
		log.Error("failed inserting file chunk values", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func valuesPlaceholder(count int) string {
	values := make([]string, count)

	fieldsNumber := 3
	for i := range count {
		values[i] = fmt.Sprintf("($%d, $%d, $%d)", i*fieldsNumber+1, i*fieldsNumber+2, i*fieldsNumber+3)
	}

	return strings.Join(values, ",")
}

func fileChunkValues(fileID int64, chunkSize int64, chunksNumber int) *[]interface{} {
	fieldsNumber := 3

	values := make([]interface{}, chunksNumber*fieldsNumber)

	for i := range chunksNumber {
		values[i*fieldsNumber] = fileID
		values[i*fieldsNumber+1] = chunkSize
		values[i*fieldsNumber+2] = i
	}

	return &values
}
