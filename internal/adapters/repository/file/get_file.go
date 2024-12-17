package file

import (
	"database/sql"
	"errors"
	"log/slog"

	"github.com/aziret/s3-mini/internal/lib/logger/sl"

	"github.com/aziret/s3-mini/internal/model"
)

func (repo *repository) GetFile(ID int64) (*model.File, error) {
	const op = "repository.file.GetFile"

	log := repo.log.With(
		slog.String("op", op),
	)

	query := `SELECT id, name, file_path, upload_id, "offset", filetype, size FROM files WHERE id = $1`
	row := repo.db.QueryRow(query, ID)

	var f model.File
	err := row.Scan(&f.ID, &f.Name, &f.FilePath, &f.UploadID, &f.Offset, &f.FileType, &f.Size)
	if errors.Is(err, sql.ErrNoRows) {
		log.Info("file not found", slog.Attr{Key: "ID", Value: slog.Int64Value(ID)})
		return nil, nil
	} else if err != nil {
		log.Error("failed to get file by ID", sl.Err(err))
		return nil, err
	}

	return &f, nil
}
