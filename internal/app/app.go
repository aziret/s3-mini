package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"github.com/aziret/s3-mini/pkg/api/filetransfer_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/aziret/s3-mini/internal/config"
	"github.com/aziret/s3-mini/internal/lib/logger/sl"
)

type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
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

	err := a.runGRPCServer()
	if err != nil {
		return err
	}

	return a.runHTTPServer()
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initHTTPServer,
		a.initGRPCServer,
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

func (a *App) initGRPCServer(_ context.Context) error {
	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))

	reflection.Register(a.grpcServer)

	filetransfer_v1.RegisterFileTransferServiceV1Server(a.grpcServer, a.serviceProvider.FileTransferImpl())

	return nil
}

func (a *App) runCronTask() error {
	a.serviceProvider.CronTask().Run()
	return nil
}

func (a *App) runHTTPServer() error {
	const op = "app.runHTTPServer"

	logger := a.serviceProvider.Logger()

	log := logger.With(
		slog.String("op", op),
	)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Error("unable to listen", sl.Err(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) runGRPCServer() error {
	const op = "app.runGRPCServer"

	logger := a.serviceProvider.Logger()

	log := logger.With(
		slog.String("op", op),
	)
	log.Info("GRPC server is running on", slog.String("address", a.serviceProvider.GRPCConfig().Address()))

	list, err := net.Listen("tcp", a.serviceProvider.GRPCConfig().Address())
	if err != nil {
		log.Error("unable to listen grpc", sl.Err(err))

		return fmt.Errorf("%s: %w", op, err)
	}

	err = a.grpcServer.Serve(list)
	if err != nil {
		log.Error("unable to start grpc", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
