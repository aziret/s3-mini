package crontask

import (
	"log/slog"

	"github.com/aziret/s3-mini/internal/service"
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

func (task *CronTask) Run() {
	task.cron.AddFunc("1-59/5 * * * *", task.fileService.CreateFileChunks)
}
