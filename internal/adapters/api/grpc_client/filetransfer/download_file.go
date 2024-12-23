package filetransfer

import (
	"context"
	"fmt"

	"github.com/aziret/s3-mini-internal/internal/lib/logger/sl"
	"github.com/aziret/s3-mini-internal/internal/model"
	"github.com/aziret/s3-mini-internal/pkg/api/filetransfer_v1"

	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) DownloadFileChunks(ctx context.Context, fileChunks <-chan model.FileChunk, downloadedFileChunks chan<- model.FileChunkDownload) error {
	const op = "grpcClient.fileTransfer.DownloadFileChunks"
	log := i.logger.With(
		slog.String("op", op),
	)

	stream, err := i.filetransferClient.DownloadFile(ctx)
	if err != nil {
		log.Error("failed to create stream for upload", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	go func() {
		for fileChunk := range fileChunks {
			err := stream.Send(&filetransfer_v1.FileChunkRequest{
				Uuid:        fileChunk.UUID,
				ChunkNumber: fileChunk.ChunkNumber,
				ChunkSize:   fileChunk.ChunkSize,
			})
			if err != nil {
				log.Error("failed to send file chunk", slog.String("UUID", fileChunk.UUID), sl.Err(err))
				return
			}
		}
		if err := stream.CloseSend(); err != nil {
			log.Error("error closing send stream", sl.Err(err))
		}
	}()

	for {
		resp, err := stream.Recv()
		if err != nil {
			if status.Code(err) == codes.Canceled {
				log.Error("file transfer stream canceled", sl.Err(err))
				break
			}
			log.Error("error receiving file chunk", sl.Err(err))
			break
		}

		downloadedFileChunks <- model.FileChunkDownload{
			UUID:        resp.GetUuid(),
			Data:        resp.GetData(),
			ChunkNumber: resp.GetChunkNumber(),
			ChunkSize:   resp.GetChunkSize(),
		}
	}

	return nil
}
