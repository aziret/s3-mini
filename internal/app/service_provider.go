package app

import (
	"log"
	"log/slog"

	"github.com/aziret/s3-mini/internal/adapters/api/grpc_server/filetransfer"

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
	log              *slog.Logger
	fileRepo         repository.FileRepository
	fileService      service.FileService
	fileImpl         *file.Implementation
	fileTransferImpl *filetransfer.Implementation
	cronTask         *crontask.CronTask
	grpcConfig       config.GRPCConfig
	transactor       repository.Transactor
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

func (s *serviceProvider) Transactor() repository.Transactor {
	const op = "serviceProvider.Transactor"

	log := s.Logger()

	logger := log.With(
		slog.String("op", op),
	)

	if s.transactor == nil {
		repo, err := fileRepository.NewRepository(log)
		if err != nil {
			logger.Error("failed to initialize repo", sl.Err(err))
		}

		s.transactor = repo
	}

	return s.transactor
}

func (s *serviceProvider) FileService() service.FileService {
	if s.fileService == nil {
		s.fileService = fileService.NewService(s.FileRepo(), s.Logger(), s.Transactor())
	}

	return s.fileService
}

func (s *serviceProvider) FileImpl() *file.Implementation {
	if s.fileImpl == nil {
		s.fileImpl = file.NewImplementation(s.Logger(), s.FileService())
	}

	return s.fileImpl
}

func (s *serviceProvider) FileTransferImpl() *filetransfer.Implementation {
	if s.fileTransferImpl == nil {
		s.fileTransferImpl = filetransfer.NewImplementation(s.FileService(), s.Logger())
	}

	return s.fileTransferImpl
}

func (s *serviceProvider) CronTask() *crontask.CronTask {
	if s.cronTask == nil {
		s.cronTask = crontask.NewCronTask(s.FileService(), s.Logger())
	}
	return s.cronTask
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get grpc config: %s", err.Error())
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}
