package file

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/aziret/s3-mini/internal/lib/logger/sl"
)

func (repo *repository) GetFileChunksServerIDs(ctx context.Context, fileID int64) (*[]string, error) {
	const op = "repository.file.GetFileChunksServerIDs"

	log := repo.log.With(
		slog.String("op", op),
		slog.Int64("fileID", fileID),
	)

	stmt, err := repo.db.Prepare("SELECT DISTINCT  server_id FROM file_chunks WHERE file_id = $1")
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

	rows, err := stmt.Query(fileID)

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

	var serverIDs []string

	for rows.Next() {
		var serverID string
		err := rows.Scan(&serverID)
		if err != nil {
			log.Error("failed to scan row", sl.Err(err))
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		serverIDs = append(serverIDs, serverID)
	}

	if err := rows.Err(); err != nil {
		log.Error("failed to iterate row", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &serverIDs, nil
}
