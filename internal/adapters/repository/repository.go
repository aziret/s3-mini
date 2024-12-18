package repository

import (
	"context"
	"errors"

	"github.com/aziret/s3-mini/internal/model"
)

var (
	ErrUploadNotFound = errors.New("upload not found")
	ErrUploadExists   = errors.New("upload already exists")
)

type FileRepository interface {
	Create(ctx context.Context, info *model.File) (int64, error)
	GetFile(ID int64) (*model.File, error)
	GetFilesWithoutChunks(ctx context.Context) (*[]model.File, error)
	CreateFileChunksForFile(_ context.Context, file *model.File) error
	MarkFileChunksCreated(ctx context.Context, id int64) error
}
