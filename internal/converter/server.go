package converter

import (
	"net"
	"strings"

	"github.com/aziret/s3-mini-internal/internal/model"
	"github.com/aziret/s3-mini-internal/pkg/api/filetransfer_v1"
)

func ToServerFromApi(req *filetransfer_v1.PingRequest, address string) *model.Server {
	host := strings.Split(address, ":")[0]
	return &model.Server{
		UUID:    req.GetUuid(),
		Address: net.JoinHostPort(host, req.Port),
	}
}
