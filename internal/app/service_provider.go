package app

import (
	"github.com/aziret/s3-mini/internal/adapters/api/file"
	"github.com/aziret/s3-mini/internal/adapters/repository"
	"github.com/aziret/s3-mini/internal/lib/logger/sl"
	"github.com/aziret/s3-mini/internal/service"
	"log/slog"

	fileRepository "github.com/aziret/s3-mini/internal/adapters/repository/file"
)

type serviceProvider struct {
	log         *slog.Logger
	fileRepo    repository.FileRepository
	fileService service.FileService
	fileImpl    *file.Implementation
}

func newServiceProvider(log *slog.Logger) *serviceProvider {
	return &serviceProvider{
		log: log,
	}
}

func (s *serviceProvider) FileRepo() repository.FileRepository {
	const op = "serviceProvider.FileRepo"

	logger := s.log.With(
		slog.String("op", op),
	)

	if s.fileRepo == nil {
		repo, err := fileRepository.NewRepository(s.log)
		if err != nil {
			logger.Error("failed to initialize repo", sl.Err(err))
		}

		s.fileRepo = repo
	}

	return s.fileRepo
}
