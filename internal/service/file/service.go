package file

import (
	"context"
	"log/slog"

	"github.com/aziret/s3-mini/internal/model"
)

type FileRepository interface {
	Create(ctx context.Context, info *model.FileInfo) (int64, error)
}

type Service struct {
	logger   *slog.Logger
	fileRepo FileRepository
}

func NewService(fileRepo FileRepository) *Service {
	return &Service{
		fileRepo: fileRepo,
	}
}
