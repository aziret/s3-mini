package file

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aziret/s3-mini/internal/lib/logger/sl"
	"github.com/aziret/s3-mini/internal/model"
)

func (s *Service) CreateFileChunks(ctx context.Context) {
	const op = "service.file.CreateFileChunks"

	log := s.logger.With(
		slog.String("op", op),
	)

	filesWithoutChunks, err := s.getFilesWithoutChunks(ctx)

	if err != nil {
		log.Error("failed to get files without chunks", sl.Err(err))
	}

	for _, fileWC := range *filesWithoutChunks {
		s.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
			err := s.createFileChunksForSpecificFile(ctx, fileWC)
			if err != nil {
				log.Error("failed to create file chunks", sl.Err(err))
				return fmt.Errorf("%s: %w", op, err)
			}

			err = s.markFileChunksCreated(ctx, fileWC.ID)
			if err != nil {
				log.Error("failed to mark file as chunks created", sl.Err(err))
				return fmt.Errorf("%s: %w", op, err)
			}

			return nil
		})
	}

}

func (s *Service) getFilesWithoutChunks(ctx context.Context) (*[]model.File, error) {
	return s.fileRepo.GetFilesWithoutChunks(ctx)
}

func (s *Service) createFileChunksForSpecificFile(ctx context.Context, file model.File) error {
	return s.fileRepo.CreateFileChunksForFile(ctx, &file)
}

func (s *Service) markFileChunksCreated(ctx context.Context, fileID int64) error {
	return s.fileRepo.MarkFileChunksCreated(ctx, fileID)
}
