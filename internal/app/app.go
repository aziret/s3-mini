package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/aziret/s3-mini/internal/config"
	"github.com/aziret/s3-mini/internal/lib/logger/sl"
)

type App struct {
	serviceProvider *serviceProvider
}

func NewApp(ctx context.Context) (*App, error) {
	app := &App{}

	err := app.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (a *App) Run() error {
	a.serviceProvider.CronTask().Run()

	return a.runHTTPServer()
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initHTTPServer,
		a.initCronTask,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	err := config.Load(".env")
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	impl := a.serviceProvider.FileImpl()

	http.Handle("/files/", http.StripPrefix("/files/", impl.Handler))
	http.Handle("/files", http.StripPrefix("/files", impl.Handler))

	// TODO: run this somewhere else
	impl.ListenUpdates(ctx)
	return nil
}

func (a *App) initCronTask(_ context.Context) error {
	a.serviceProvider.CronTask()
	return nil
}

func (a *App) runCronTask() error {
	a.serviceProvider.CronTask().Run()
	return nil
}

func (a *App) runHTTPServer() error {
	const op = "app.runHTTPServer"

	log := a.serviceProvider.Logger()

	logger := log.With(
		slog.String("op", op),
	)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		logger.Error("unable to listen: %s", sl.Err(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
