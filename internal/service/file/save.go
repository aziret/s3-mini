package file

import (
	"context"
	"log/slog"

	"github.com/aziret/s3-mini/internal/lib/logger/sl"

	"github.com/aziret/s3-mini/internal/model"
)

func (s *Service) Save(ctx context.Context, file *model.File) (*model.File, error) {
	const op = "service.file.Save"

	logger := s.logger.With(
		slog.String("op", op),
	)

	id, err := s.fileRepo.Create(ctx, file)
	if err != nil {
		logger.Error("failed to save the file", sl.Err(err))
		return nil, err
	}

	file.ID = id

	return file, nil
}
