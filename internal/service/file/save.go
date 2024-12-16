package file

import (
	"context"
	"github.com/aziret/s3-mini/internal/lib/logger/sl"
	"log/slog"

	"github.com/aziret/s3-mini/internal/model"
)

func (s *Service) Save(ctx context.Context, info *model.FileInfo) (*model.File, error) {
	const op = "service.file.Save"

	logger := s.logger.With(
		slog.String("op", op),
	)

	id, err := s.fileRepo.Create(ctx, info)
	if err != nil {
		logger.Error("failed to save the file", sl.Err(err))
		return nil, err
	}

	uploadedFile := model.File{
		ID:   id,
		Info: *info,
	}

	return &uploadedFile, nil
}
