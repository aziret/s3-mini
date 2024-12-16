package file

import (
	"context"
	"log/slog"
)

func (impl *Implementation) ListenUpdates(_ context.Context) {
	const op = "api.file.ListenUpdates"

	log := impl.logger.With(
		slog.String("op", op),
	)
	go func() {
		for {
			event := <-impl.handler.CompleteUploads

			log.Info("Upload finished", slog.Attr{Key: "ID", Value: slog.StringValue(event.Upload.ID)})
		}
	}()
}
