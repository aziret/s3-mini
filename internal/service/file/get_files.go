package file

import (
	"context"
	"github.com/aziret/s3-mini-internal/internal/lib/logger/sl"
	"github.com/aziret/s3-mini-internal/internal/model"
	"log/slog"
)

func (s *Service) GetFiles(ctx context.Context) (*[]model.FileInfo, error) {
	const op = "service.file.GetFiles"

	log := s.logger.With(
		slog.String("op", op),
	)

	files, err := s.fileRepo.GetFiles(ctx)
	if err != nil {
		log.Error("failed to get files", sl.Err(err))
	}
	return files, err
}
