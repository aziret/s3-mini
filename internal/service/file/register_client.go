package file

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/aziret/s3-mini/internal/adapters/repository"
	"github.com/aziret/s3-mini/internal/model"
)

func (s *Service) RegisterClient(ctx context.Context, server *model.Server) error {
	const op = "service.file.RegisterClient"

	log := s.logger.With(
		slog.String("op", op),
	)

	err := s.fileRepo.RegisterClient(ctx, server)
	if err != nil {
		if errors.Is(err, repository.ErrServerExists) {
			log.Info("server already created", slog.String("serverID", server.UUID))
		} else {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	s.updateServerData(server)

	return nil
}

func (s *Service) updateServerData(server *model.Server) {
	s.serversMap.mu.Lock()
	defer s.serversMap.mu.Unlock()
	s.serversMap.servers[server.UUID] = *server
}
