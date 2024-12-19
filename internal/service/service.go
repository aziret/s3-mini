package service

import (
	"context"

	"github.com/aziret/s3-mini/internal/model"
)

type FileService interface {
	Save(ctx context.Context, file *model.File) (*model.File, error)
	CreateFileChunks()
	RegisterClient(ctx context.Context, server *model.Server) error
}
