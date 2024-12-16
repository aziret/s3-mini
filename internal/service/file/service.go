package file

import (
	"github.com/aziret/s3-mini/internal/adapters/repository"
	"log/slog"
)

type Service struct {
	logger   *slog.Logger
	fileRepo repository.FileRepository
}

func NewService(fileRepo repository.FileRepository) *Service {
	return &Service{
		fileRepo: fileRepo,
	}
}
