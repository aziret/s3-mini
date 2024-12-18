package file

import (
	"context"
	"log/slog"

	"github.com/aziret/s3-mini/internal/lib/logger/sl"
	"github.com/aziret/s3-mini/internal/model"
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
	}
}

func (s *Service) getFilesWithoutChunks(ctx context.Context) (*[]model.File, error) {
	return s.fileRepo.GetFilesWithoutChunks(ctx)
}

func (s *Service) createFileChunksForSpecificFile(ctx context.Context, file model.File) error {
	return s.fileRepo.CreateFileChunksForFile(ctx, &file)
}
