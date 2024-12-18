package app

import (
	"log/slog"

	"github.com/aziret/s3-mini/internal/adapters/api/file"
	"github.com/aziret/s3-mini/internal/adapters/infra/crontask"
	"github.com/aziret/s3-mini/internal/adapters/repository"
	"github.com/aziret/s3-mini/internal/config"
	"github.com/aziret/s3-mini/internal/lib/logger/sl"
	"github.com/aziret/s3-mini/internal/service"

	fileRepository "github.com/aziret/s3-mini/internal/adapters/repository/file"
	fileService "github.com/aziret/s3-mini/internal/service/file"
)

type serviceProvider struct {
	log         *slog.Logger
	fileRepo    repository.FileRepository
	fileService service.FileService
	fileImpl    *file.Implementation
	cronTask    *crontask.CronTask
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) Logger() *slog.Logger {
	if s.log == nil {
		s.log = config.NewLogger()
	}

	return s.log
}

func (s *serviceProvider) FileRepo() repository.FileRepository {
	const op = "serviceProvider.FileRepo"

	log := s.Logger()

	logger := log.With(
		slog.String("op", op),
	)

	if s.fileRepo == nil {
		repo, err := fileRepository.NewRepository(log)
		if err != nil {
			logger.Error("failed to initialize repo", sl.Err(err))
		}

		s.fileRepo = repo
	}

	return s.fileRepo
}

func (s *serviceProvider) FileService() service.FileService {
	if s.fileService == nil {
		s.fileService = fileService.NewService(s.FileRepo(), s.Logger())
	}

	return s.fileService
}

func (s *serviceProvider) FileImpl() *file.Implementation {
	if s.fileImpl == nil {
		s.fileImpl = file.NewImplementation(s.Logger(), s.FileService())
	}

	return s.fileImpl
}

func (s *serviceProvider) CronTask() *crontask.CronTask {
	if s.cronTask == nil {
		s.cronTask = crontask.NewCronTask(s.FileService(), s.Logger())
	}
	return s.cronTask
}
