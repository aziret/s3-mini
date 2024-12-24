package file

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aziret/s3-mini/internal/lib/logger/sl"
	"github.com/aziret/s3-mini/internal/model"
)

func (repo *repository) GetNotSentFileChunks(_ context.Context) (*[]model.FileChunk, error) {
	const op = "repository.file.GetNotSentFileChunks"

	log := repo.log.With(
		slog.String("op", op),
	)

	stmt, err := repo.db.Prepare("SELECT file_chunks.id, file_chunks.chunk_size, file_chunks.chunk_number, files.file_path FROM file_chunks INNER JOIN files ON file_chunks.file_id = files.id WHERE file_chunks.download_completed = FALSE")
	if err != nil {
		log.Error("failed to prepare statement", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
			log.Error("failed to close statement", sl.Err(err))
		}
	}()

	rows, err := stmt.Query()
	if err != nil {
		log.Error("failed to query statement", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Error("failed to close row", sl.Err(err))
		}
	}()

	var fileChunks []model.FileChunk

	for rows.Next() {
		var fileChunk model.FileChunk
		if err := rows.Scan(&fileChunk.UUID, &fileChunk.ChunkSize, &fileChunk.ChunkNumber, &fileChunk.FilePath); err != nil {
			log.Error("failed to scan row", sl.Err(err))
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		fileChunks = append(fileChunks, fileChunk)
	}

	if err := rows.Err(); err != nil {
		log.Error("failed to iterate row", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &fileChunks, nil
}
