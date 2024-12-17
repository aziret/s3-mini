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
}
