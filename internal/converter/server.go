package converter

import (
	"github.com/aziret/s3-mini-internal/internal/model"
	"github.com/aziret/s3-mini-internal/pkg/api/filetransfer_v1"
)

func ToServerFromApi(req *filetransfer_v1.PingRequest, address string) *model.Server {
	return &model.Server{
		UUID:    req.GetUuid(),
		Address: address,
	}
}
