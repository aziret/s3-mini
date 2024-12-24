package file

import (
	"context"
	"fmt"
	"log/slog"

	repoPackage "github.com/aziret/s3-mini/internal/adapters/repository"
	"github.com/aziret/s3-mini/internal/lib/logger/sl"
	"github.com/aziret/s3-mini/internal/model"
	"github.com/lib/pq"
)

func (repo *repository) RegisterClient(_ context.Context, server *model.Server) error {
	const op = "repository.file.RegisterClient"

	log := repo.log.With(
		slog.String("op", op),
	)

	stmt, err := repo.db.Prepare("INSERT INTO servers (id) VALUES ($1) ON CONFLICT DO NOTHING")
	if err != nil {
		log.Error("failed to prepare statement", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(server.UUID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code.Name() == "unique_violation" {
			return fmt.Errorf("%s: %w", op, repoPackage.ErrServerExists)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
