package file

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/aziret/s3-mini/internal/lib/logger/sl"
)

func (repo *repository) WithinTransaction(ctx context.Context, tFunc func(tctx context.Context) error) error {
	const op = "repository.file.WithinTransaction"

	log := repo.log.With(
		slog.String("op", op),
	)

	tx, err := repo.db.Begin()

	if err != nil {
		log.Error("failed to begin transaction", sl.Err(err))

		return fmt.Errorf("%s, begin transaction: %w", op, err)
	}

	defer func() {
		if err == nil {
			return
		}
		if errRollback := tx.Rollback(); errRollback != nil {
			log.Error("%s, rollback transaction: %v", op, errRollback)
		}
		log.Error("transaction failed, %s: %w", op, err)
	}()

	err = tFunc(injectTx(ctx, tx))
	if err != nil {
		return fmt.Errorf("transaction failed, %s: %w", op, err)
	}

	if err = tx.Commit(); err != nil {
		log.Error("failed to commit transaction, %s: %w", op, err)
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

type txKey struct{}

func injectTx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func extractTx(ctx context.Context) *sql.Tx {
	if tx, ok := ctx.Value(txKey{}).(*sql.Tx); ok {
		return tx
	}
	return nil
}
