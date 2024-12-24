package filetransfer

import (
	"log/slog"

	"github.com/aziret/s3-mini/internal/lib/logger/sl"
	"github.com/aziret/s3-mini/pkg/api/filetransfer_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Implementation struct {
	logger             *slog.Logger
	filetransferClient filetransfer_v1.FileTransferServiceV1Client
}

func NewImplementation(logger *slog.Logger, storageAddress string) *Implementation {
	const op = "grpcClient.filetransfer.NewImplementation"
	log := logger.With(
		slog.String("op", op),
	)

	conn, err := grpc.NewClient(storageAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("Failed to connect to storage", sl.Err(err))

		panic(err)
	}

	client := filetransfer_v1.NewFileTransferServiceV1Client(conn)

	return &Implementation{
		logger:             logger,
		filetransferClient: client,
	}
}
