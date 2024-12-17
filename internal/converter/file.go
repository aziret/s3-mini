package converter

import (
	"github.com/aziret/s3-mini/internal/model"
	tusd "github.com/tus/tusd/v2/pkg/handler"
)

func ToFileFromApi(info *tusd.FileInfo) *model.File {
	return &model.File{
		UploadID: info.ID,
		Size:     info.Size,
		Offset:   info.Offset,
		Name:     info.MetaData["filename"],
		FilePath: info.Storage["Path"],
		FileType: info.MetaData["filetype"],
	}
}
