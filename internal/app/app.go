package app

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/aziret/s3-mini/internal/adapters/api/http/file"
	"github.com/aziret/s3-mini/pkg/api/filetransfer_v1"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/aziret/s3-mini/internal/config"
	"github.com/aziret/s3-mini/internal/lib/logger/sl"
)

type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
	fiberServer     *fiber.App
}

func NewApp(ctx context.Context) (*App, error) {
	app := &App{}

	err := app.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (a *App) Run(ctx context.Context) error {
	a.runCronTask(ctx)

	go func() {
		err := a.runGRPCServer()
		if err != nil {
			fmt.Println("Error running GRPC Server")
			panic(err)
		}
	}()

	go func() {
		if err := a.runHTTPServer(); err != nil {
			fmt.Println("Error running HTTP server")
			panic(err)
		}
	}()

	return a.runFiberServer(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initHTTPServer,
		a.initGRPCServer,
		a.initCronTask,
		a.initFiberServer,
		a.initDownloadFolder,
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

func (a *App) initFiberServer(ctx context.Context) error {
	a.fiberServer = fiber.New()
	return nil
}

func (a *App) initDownloadFolder(_ context.Context) error {
	const op = "app.initDownloadFolder"

	logger := a.serviceProvider.Logger()

	log := logger.With(
		slog.String("op", op),
	)

	err := os.MkdirAll("./downloads/", os.ModePerm)
	if err != nil {
		log.Error("could not create directory", sl.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (a *App) runCronTask(ctx context.Context) {
	a.serviceProvider.CronTask().Run(ctx)
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
	log.Info("running and serving HTTP server")

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

func (a *App) runFiberServer(_ context.Context) error {
	fileHandler := file.NewFileHandler(a.serviceProvider.FileService())
	a.fiberServer.Use(cors.New())
	a.fiberServer.Get("/files", fileHandler.GetFiles)
	a.fiberServer.Get("/files/:id", fileHandler.GetFile)

	port := os.Getenv("FIBER_SERVER_PORT")

	log.Printf("Starting server on http://localhost:%s", port)
	log.Fatal(a.fiberServer.Listen(":" + port))
	return nil
}
