package file

import (
	"log/slog"
	"sync"

	"github.com/aziret/s3-mini/internal/adapters/api/grpc_client/filetransfer"

	"github.com/aziret/s3-mini/internal/adapters/repository"
)

type Service struct {
	logger     *slog.Logger
	fileRepo   repository.FileRepository
	transactor repository.Transactor
	serversMap serversMap
}

type serversMap struct {
	mu      sync.RWMutex
	servers map[string]*filetransfer.Implementation
}

func NewService(fileRepo repository.FileRepository, logger *slog.Logger, transactor repository.Transactor) *Service {
	return &Service{
		fileRepo: fileRepo,
		logger:   logger,
		serversMap: serversMap{
			servers: make(map[string]*filetransfer.Implementation),
		},
		transactor: transactor,
	}
}
