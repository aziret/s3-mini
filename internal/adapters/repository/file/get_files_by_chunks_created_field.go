package file

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aziret/s3-mini-internal/internal/lib/logger/sl"
	"github.com/aziret/s3-mini-internal/internal/model"
)

func (repo *repository) GetFilesWithoutChunks(ctx context.Context) (*[]model.File, error) {
	return repo.getFilesByChunksCreatedField(ctx, false)
}

func (repo *repository) getFilesByChunksCreatedField(_ context.Context, file_chunks_created bool) (*[]model.File, error) {
	const op = "repository.file.getFilesByChunksCreatedField"

	log := repo.log.With(
		slog.String("op", op),
		slog.Bool("file_chunks_created", file_chunks_created),
	)

	query := `SELECT id, name, file_path, upload_id, "offset", filetype, size FROM files WHERE file_chunks_created = $1`
	rows, err := repo.db.Query(query, file_chunks_created)
	if err != nil {
		log.Error("failed to run query for fetching files", sl.Err(err))
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Error("failed to close rows for fetching files", sl.Err(err))
		}
	}()

	var files []model.File

	for rows.Next() {
		var file model.File

		err := rows.Scan(&file.ID, &file.Name, &file.FilePath, &file.UploadID, &file.Offset, &file.FileType, &file.Size)
		if err != nil {
			log.Error("failed to scan file row", sl.Err(err))
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		files = append(files, file)
	}

	if err := rows.Err(); err != nil {
		log.Error("failed to fetch files", sl.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &files, nil
}

func (repo *repository) GetFilesWithChunks(ctx context.Context) (*[]model.File, error) {
	return repo.getFilesByChunksCreatedField(ctx, true)
}
