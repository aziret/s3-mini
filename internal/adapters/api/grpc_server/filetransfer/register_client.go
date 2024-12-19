package filetransfer

import (
	"context"
	"log/slog"

	"github.com/aziret/s3-mini/internal/converter"
	"github.com/aziret/s3-mini/internal/lib/logger/sl"
	"github.com/aziret/s3-mini/pkg/api/filetransfer_v1"
	"github.com/pkg/errors"
	"google.golang.org/grpc/peer"
)

const peerExtractionFailed = "failed to extract peer from ctx"

func (i *Implementation) RegisterClient(ctx context.Context, req *filetransfer_v1.PingRequest) (*filetransfer_v1.PongResponse, error) {
	const op = "filetransfer.RegisterClient"

	log := i.logger.With(
		slog.String("op", op),
	)

	resp := &filetransfer_v1.PongResponse{}

	// Get peer info from the context
	p, ok := peer.FromContext(ctx)
	if !ok {
		log.Error(peerExtractionFailed, slog.String("UUID", req.GetUuid()))

		resp.Success = false
		resp.Message = peerExtractionFailed
		return resp, errors.New(peerExtractionFailed)
	}

	err := i.fileService.RegisterClient(ctx, converter.ToServerFromApi(req, p.Addr.String()))

	if err != nil {
		log.Error("failed to register client", sl.Err(err))

		resp.Success = false
		resp.Message = "invalid data provided"
		return resp, errors.New(resp.Message)
	}
	resp.Success = true
	resp.Message = "Client saved successfully"
	return resp, nil
}
