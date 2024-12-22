package service

import (
	"context"

	"github.com/aziret/s3-mini-internal/internal/model"
)

type FileService interface {
	Save(ctx context.Context, file *model.File) (*model.File, error)
	CreateFileChunks(ctx context.Context)
	RegisterClient(ctx context.Context, server *model.Server) error
	UploadFileChunks(ctx context.Context)
	MarkFilesAsUploadCompleted(ctx context.Context)
	RemoveUploadedFiles(ctx context.Context)
}
