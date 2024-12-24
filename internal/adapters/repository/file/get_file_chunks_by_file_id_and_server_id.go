package file

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/aziret/s3-mini/internal/lib/logger/sl"
	"github.com/aziret/s3-mini/internal/model"
)

func (repo *repository) GetFileChunksByFileIDAndServerID(_ context.Context, fileID int64, serverID string) (*[]model.FileChunk, error) {
	const op = "repository.file.GetFileChunksByFileID"
	log := repo.log.With(
		slog.String("op", op),
	)

	stmt, err := repo.db.Prepare("SELECT id, chunk_size, chunk_number, server_id FROM file_chunks WHERE file_id = $1 AND server_id = $2")
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

	rows, err := stmt.Query(fileID, serverID)

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

	var fileChunks []model.FileChunk

	for rows.Next() {
		var fileChunk model.FileChunk
		if err := rows.Scan(&fileChunk.UUID, &fileChunk.ChunkSize, &fileChunk.ChunkNumber, &fileChunk.ServerID); err != nil {
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
