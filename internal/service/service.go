package service

import (
	"context"
	"github.com/aziret/s3-mini/internal/model"
)

type FileService interface {
	Save(ctx context.Context, info *model.FileInfo) (*model.File, error)
}
