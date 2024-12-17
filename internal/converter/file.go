package converter

import (
	"github.com/aziret/s3-mini/internal/model"
	tusd "github.com/tus/tusd/v2/pkg/handler"
)

func ToFileInfoFromApi(info *tusd.FileInfo) *model.FileInfo {
	return &model.FileInfo{
		ID:             info.ID,
		Size:           info.Size,
		SizeIsDeferred: info.SizeIsDeferred,
		Offset:         info.Offset,
		MetaData:       info.MetaData,
		IsPartial:      info.IsPartial,
		IsFinal:        info.IsFinal,
		PartialUploads: info.PartialUploads,
		Storage:        info.Storage,
	}
}
