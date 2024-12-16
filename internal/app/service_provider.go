package app

import (
	"github.com/aziret/s3-mini/internal/adapters/api/file"
	"github.com/aziret/s3-mini/internal/adapters/repository"
	"github.com/aziret/s3-mini/internal/service"
)

type serviceProvider struct {
	fileRepo    repository.FileRepository
	fileService service.FileService
	fileImpl    *file.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}
