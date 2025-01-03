package filetransfer

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aziret/s3-mini/internal/lib/logger/sl"
	"github.com/aziret/s3-mini/internal/model"
	"github.com/aziret/s3-mini/pkg/api/filetransfer_v1"
)

func (i *Implementation) UploadFileChunks(ctx context.Context, fileChunksChan chan model.FileChunkUpload) error {
	const op = "grpcClient.fileTransfer.UploadFileChunks"
	log := i.logger.With(
		slog.String("op", op),
	)

	stream, err := i.filetransferClient.UploadFile(ctx)
	if err != nil {
		log.Error("failed to create stream for upload", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	for fileChunkUpload := range fileChunksChan {
		err = stream.Send(&filetransfer_v1.FileChunkUpload{
			Uuid: fileChunkUpload.UUID,
			Data: fileChunkUpload.Data,
		})
		if err != nil {
			log.Error("failed to send file chunk", sl.Err(err), slog.String("UUID", fileChunkUpload.UUID))
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	_, err = stream.CloseAndRecv()
	if err != nil {
		log.Error("failed to close stream", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
