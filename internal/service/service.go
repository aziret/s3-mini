package service

import (
	"context"

	"github.com/aziret/s3-mini/internal/model"
)

type FileService interface {
	Save(ctx context.Context, info *model.File) (*model.File, error)
}
