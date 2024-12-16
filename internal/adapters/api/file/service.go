package file

import (
	"github.com/aziret/s3-mini/internal/service"
	"net/http"

	"github.com/aziret/s3-mini/internal/lib/logger/sl"
	"github.com/tus/tusd/v2/pkg/filelocker"
	"github.com/tus/tusd/v2/pkg/filestore"
	tusd "github.com/tus/tusd/v2/pkg/handler"

	"log/slog"
)

type Implementation struct {
	handler     *tusd.Handler
	fileService service.FileService
	logger      *slog.Logger
}

func NewImplementation(logger *slog.Logger, fileService service.FileService) *Implementation {
	const op = "api.file.New"
	log := logger.With(
		slog.String("op", op),
	)

	store := filestore.New("./uploads")
	locker := filelocker.New("./uploads")
	composer := tusd.NewStoreComposer()
	store.UseIn(composer)
	locker.UseIn(composer)

	handler, err := tusd.NewHandler(tusd.Config{
		BasePath:              "/files/",
		StoreComposer:         composer,
		NotifyCompleteUploads: true,
	})
	if err != nil {
		log.Error("unable to create handler: %s", sl.Err(err))
	}

	http.Handle("/files/", http.StripPrefix("/files/", handler))
	http.Handle("/files", http.StripPrefix("/files", handler))
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Error("unable to listen: %s", sl.Err(err))
	}

	return &Implementation{
		handler:     handler,
		logger:      logger,
		fileService: fileService,
	}
}
