package file

import (
	"context"
	"log/slog"

	"github.com/aziret/s3-mini-internal/internal/lib/logger/sl"
)

func (s *Service) MarkFilesAsUploadCompleted(ctx context.Context) {
	const op = "service.file.MarkFilesAsUploadCompleted"

	log := s.logger.With(
		slog.String("op", op),
	)

	err := s.fileRepo.MarkFilesAsUploadCompleted(ctx)

	if err != nil {
		log.Error("failed to mark files as upload completed", sl.Err(err))
	}
}
