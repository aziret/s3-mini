package file

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/aziret/s3-mini-internal/internal/lib/logger/sl"
)

func (s *Service) RemoveUploadedFiles(ctx context.Context) {
	const op = "service.file.RemoveUploadedFiles"

	log := s.logger.With(
		slog.String("op", op),
	)

	uploadCompletedFiles, err := s.fileRepo.GetUploadCompletedFiles(ctx)
	if err != nil {
		log.Error("failed to get uploaded files", sl.Err(err))
	}

	if len(*uploadCompletedFiles) == 0 {
		return
	}

	deletedFileIds := make([]int64, 0, len(*uploadCompletedFiles))

	for _, file := range *uploadCompletedFiles {
		err = s.removeUploadedFile(ctx, file.FilePath)
		if err != nil {
			log.Error("failed to remove uploaded file", sl.Err(err))
			continue
		}
		deletedFileIds = append(deletedFileIds, file.ID)
	}

	if len(deletedFileIds) > 0 {
		err = s.fileRepo.MarkFilesAsReadyToDownload(ctx, deletedFileIds)

		if err != nil {
			log.Error("failed to mark files as ready to download", sl.Err(err))
		} else {
			log.Info("successfully marked files as ready to download")
		}
	}
}

func (s *Service) removeUploadedFile(_ context.Context, filePath string) error {
	const op = "service.file.removeUploadedFile"
	log := s.logger.With(slog.String("filePath", filePath))

	err := os.Remove(filePath)
	if err != nil {
		log.Error(
			"failed to remove file",
			slog.String("filePath", filePath),
			sl.Err(err),
		)

		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
