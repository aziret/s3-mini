package file

import (
	"net/http"

	"github.com/tus/tusd/v2/pkg/filelocker"
	"github.com/tus/tusd/v2/pkg/filestore"
	tusd "github.com/tus/tusd/v2/pkg/handler"

	"log/slog"
)

type Implementation struct {
}

func NewImplementation(logger *slog.Logger) *Implementation {

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
		logger.Error("unable to create handler: %s", err)
	}

	go func() {
		for {
			event := <-handler.CompleteUploads
			logger.Info("Upload %s finished\n", event.Upload.ID)
		}
	}()

	http.Handle("/files/", http.StripPrefix("/files/", handler))
	http.Handle("/files", http.StripPrefix("/files", handler))
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		logger.Error("unable to listen: %s", err)
	}

	return &Implementation{}
}
