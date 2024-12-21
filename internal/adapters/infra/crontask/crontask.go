package crontask

import (
	"context"
	"log/slog"

	"github.com/aziret/s3-mini-internal/internal/service"
	"github.com/robfig/cron/v3"
)

type CronTask struct {
	fileService service.FileService
	logger      *slog.Logger
	cron        *cron.Cron
}

func NewCronTask(fileService service.FileService, logger *slog.Logger) *CronTask {
	return &CronTask{
		fileService: fileService,
		logger:      logger,
		cron:        cron.New(),
	}
}

func (task *CronTask) Run(ctx context.Context) {
	// TODO: handle error returned by AddFunc
	task.cron.AddFunc("1-59/5 * * * *", wrapFunction(ctx, task.fileService.CreateFileChunks))
	task.cron.AddFunc("2-59/5 * * * *", wrapFunction(ctx, task.fileService.UploadFileChunks))
	task.cron.AddFunc("3-59/5 * * * *", wrapFunction(ctx, task.fileService.MarkFilesAsUploadCompleted))
	task.cron.Start()
}

func wrapFunction(ctx context.Context, f func(ctx context.Context)) func() {
	return func() {
		f(ctx)
	}
}
