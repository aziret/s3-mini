package file

import (
	"context"
	"errors"
	"log/slog"

	"github.com/aziret/s3-mini/internal/adapters/repository"
	"github.com/aziret/s3-mini/internal/converter"
	"github.com/aziret/s3-mini/internal/lib/logger/sl"
)

func (impl *Implementation) ListenUpdates(ctx context.Context) {
	const op = "api.file.ListenUpdates"

	log := impl.logger.With(
		slog.String("op", op),
	)
	go func() {
		for {
			event := <-impl.Handler.CompleteUploads

			log.Info("Upload finished", slog.Attr{Key: "ID", Value: slog.StringValue(event.Upload.ID)})

			_, err := impl.fileService.Save(ctx, converter.ToFileFromApi(&event.Upload))

			if err != nil {
				if errors.Is(err, repository.ErrUploadExists) {
					log.Info("File already exists", slog.Attr{Key: "Upload ID", Value: slog.StringValue(event.Upload.ID)})
				}

				log.Error("Error saving file", sl.Err(err))
			}

		}
	}()
}
