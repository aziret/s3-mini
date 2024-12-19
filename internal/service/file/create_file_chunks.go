package file

import (
	"context"
	"log/slog"

	"github.com/aziret/s3-mini-internal/internal/lib/logger/sl"
	"github.com/aziret/s3-mini-internal/internal/model"
)

func (s *Service) CreateFileChunks() {
	ctx := context.Background()
	const op = "service.file.CreateFileChunks"

	log := s.logger.With(
		slog.String("op", op),
	)

	filesWithoutChunks, err := s.getFilesWithoutChunks(ctx)

	if err != nil {
		log.Error("failed to get files without chunks", sl.Err(err))
	}

	for _, fileWC := range *filesWithoutChunks {
		err := s.createFileChunksForSpecificFile(ctx, fileWC)
		if err != nil {
			log.Error("failed to create file chunks", sl.Err(err))
		}

		if err != nil {
			continue
		}

		err = s.markFileChunksCreated(ctx, fileWC.ID)
		if err != nil {
			log.Error("failed to mark file as chunks created", sl.Err(err))
		}
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
