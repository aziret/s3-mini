package file

import (
	"context"
	"fmt"
	"github.com/aziret/s3-mini-internal/internal/lib/logger/sl"
	"log/slog"
)

func (s *Service) GetFile(ctx context.Context, fileID int64) (string, error) {
	const op = "service.file.GetFile"
	log := s.logger.With(
		slog.String("op", op),
	)

	file, err := s.fileRepo.GetFile(ctx, fileID)

	if err != nil {
		log.Error("failed to get file", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}
	if file == nil {
		log.Info("file not found")
		return "", fmt.Errorf("file not found")
	}

	if file.ReadyToDownload {
		return s.DownloadFile(ctx, fileID)
	}

	return file.FilePath, nil
}
