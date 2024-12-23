package file

import "github.com/aziret/s3-mini-internal/internal/service"

type FileHandler struct {
	fileService service.FileService
}

func NewFileHandler(service service.FileService) *FileHandler {
	return &FileHandler{
		fileService: service,
	}
}
