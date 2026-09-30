package bootstrap

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dev-gopi/authhub/internal/config"
	"github.com/dev-gopi/authhub/internal/platform/httpserver"
	"github.com/dev-gopi/authhub/internal/platform/logger"
	sharedmiddleware "github.com/dev-gopi/authhub/internal/shared/middleware"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func RunAuthAPI() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	appLogger, err := logger.New(
		cfg.App.Environment,
	)
	if err != nil {
		return fmt.Errorf("create logger: %w", err)
	}

	deps, err := BuildDependencies(
		cfg,
		appLogger,
	)
	if err != nil {
		_ = appLogger.Sync()

		return err
	}

	defer deps.Close()

	router := chi.NewRouter()

	router.Use(
		sharedmiddleware.RequestID,
	)

	router.Use(
		sharedmiddleware.CorrelationID,
	)

	router.Use(
		sharedmiddleware.Recovery(appLogger),
	)

	router.Use(
		sharedmiddleware.RequestLogger(appLogger),
	)

	router.Use(
		sharedmiddleware.Timeout(
			cfg.App.RequestTimeout,
		),
	)

	health := httpserver.NewHealthHandler(
		deps.Postgres,
		deps.Redis,
		deps.RabbitMQ,
		deps.Vault,
		"auth-api",
	)

	router.Get(
		"/health",
		health.Health,
	)

	router.Get(
		"/ready",
		health.Ready,
	)

	server := httpserver.New(
		cfg.App.AuthAPIAddr,
		router,
	)

	serverErrors := make(chan error, 1)

	go func() {
		appLogger.Info(
			"auth api started",
			zap.String(
				"address",
				cfg.App.AuthAPIAddr,
			),
		)

		serverErrors <- server.Run()
	}()

	shutdownSignal := make(
		chan os.Signal,
		1,
	)

	signal.Notify(
		shutdownSignal,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	select {
	case err := <-serverErrors:
		return err

	case signalReceived := <-shutdownSignal:

		appLogger.Info(
			"shutdown signal received",
			zap.String(
				"signal",
				signalReceived.String(),
			),
		)
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		cfg.App.ShutdownTimeout,
	)

	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf(
			"shutdown server: %w",
			err,
		)
	}

	appLogger.Info(
		"auth api stopped",
	)

	return nil
}
