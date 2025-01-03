package repository

import (
	"context"
	"errors"

	"github.com/aziret/s3-mini/internal/model"
)

var (
	ErrUploadNotFound = errors.New("upload not found")
	ErrUploadExists   = errors.New("upload already exists")
	ErrServerExists   = errors.New("server already exists")
)

type FileRepository interface {
	Create(ctx context.Context, info *model.File) (int64, error)
	GetFile(ctx context.Context, ID int64) (*model.File, error)
	GetFilesWithoutChunks(ctx context.Context) (*[]model.File, error)
	CreateFileChunksForFile(ctx context.Context, file *model.File) error
	MarkFileChunksCreated(ctx context.Context, id int64) error
	RegisterClient(ctx context.Context, server *model.Server) error
	GetNotSentFileChunks(ctx context.Context) (*[]model.FileChunk, error)
	MarkFileChunkSuccessfullyUploaded(UUID string, serverID string) error
	MarkFilesAsUploadCompleted(ctx context.Context) error
	GetUploadCompletedFiles(ctx context.Context) (*[]model.File, error)
	MarkFilesAsReadyToDownload(ctx context.Context, deletedFileIds []int64) error
	GetFileChunksByFileIDAndServerID(ctx context.Context, fileID int64, serverID string) (*[]model.FileChunk, error)
	GetFileChunksServerIDs(ctx context.Context, fileID int64) (*[]string, error)
	GetFiles(ctx context.Context) (*[]model.FileInfo, error)
}
