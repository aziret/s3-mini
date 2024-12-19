package file

import (
	"log/slog"
	"sync"

	"github.com/aziret/s3-mini/internal/adapters/repository"
	"github.com/aziret/s3-mini/internal/model"
)

type Service struct {
	logger     *slog.Logger
	fileRepo   repository.FileRepository
	serversMap serversMap
}

type serversMap struct {
	mu      sync.RWMutex
	servers map[string]model.Server
}

func NewService(fileRepo repository.FileRepository, logger *slog.Logger) *Service {
	return &Service{
		fileRepo: fileRepo,
		logger:   logger,
		serversMap: serversMap{
			servers: make(map[string]model.Server),
		},
	}
}
