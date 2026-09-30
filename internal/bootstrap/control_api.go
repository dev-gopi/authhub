package bootstrap

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	rootcontroller "github.com/dev-gopi/authhub/internal/modules/rootauth/controller"
	rootmiddleware "github.com/dev-gopi/authhub/internal/modules/rootauth/middleware"
	rootrepository "github.com/dev-gopi/authhub/internal/modules/rootauth/repository"
	rootrouter "github.com/dev-gopi/authhub/internal/modules/rootauth/router"
	rootservice "github.com/dev-gopi/authhub/internal/modules/rootauth/service"
	rootvalidator "github.com/dev-gopi/authhub/internal/modules/rootauth/validator"

	"github.com/dev-gopi/authhub/internal/config"
	"github.com/dev-gopi/authhub/internal/platform/httpserver"
	"github.com/dev-gopi/authhub/internal/platform/logger"
	sharedmiddleware "github.com/dev-gopi/authhub/internal/shared/middleware"
	"github.com/dev-gopi/authhub/internal/shared/security"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

func RunControlAPI() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf(
			"load config: %w",
			err,
		)
	}

	appLogger, err := logger.New(
		cfg.App.Environment,
	)
	if err != nil {
		return fmt.Errorf(
			"create logger: %w",
			err,
		)
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

	userRepo :=
		rootrepository.NewPlatformUserRepository(
			deps.Postgres.DB,
		)

	passwordRepo :=
		rootrepository.NewPasswordRepository(
			deps.Postgres.DB,
		)

	sessionRepo :=
		rootrepository.NewSessionRepository(
			deps.Postgres.DB,
		)

	loginAttemptRepo :=
		rootrepository.NewLoginAttemptRepository(
			deps.Postgres.DB,
		)

	passwordHasher :=
		security.NewPasswordHasher()

	rootService, err :=
		rootservice.New(
			deps.Postgres.DB,
			deps.Redis,
			userRepo,
			passwordRepo,
			sessionRepo,
			loginAttemptRepo,
			passwordHasher,
			appLogger,
		)

	if err != nil {
		return fmt.Errorf(
			"initialize root auth service: %w",
			err,
		)
	}

	rootValidator :=
		rootvalidator.New(
			deps.Validator,
		)

	rootController :=
		rootcontroller.New(
			rootService,
			rootValidator,
		)

	rootMiddleware :=
		rootmiddleware.New(
			rootService,
		)

	router := chi.NewRouter()

	router.Use(
		sharedmiddleware.RequestID,
	)

	router.Use(
		sharedmiddleware.CorrelationID,
	)

	router.Use(
		sharedmiddleware.Recovery(
			appLogger,
		),
	)

	router.Use(
		sharedmiddleware.RequestLogger(
			appLogger,
		),
	)

	router.Use(
		sharedmiddleware.Timeout(
			cfg.App.RequestTimeout,
		),
	)

	health :=
		httpserver.NewHealthHandler(
			deps.Postgres,
			deps.Redis,
			deps.RabbitMQ,
			deps.Vault,
			"control-api",
		)

	router.Get(
		"/health",
		health.Health,
	)

	router.Get(
		"/ready",
		health.Ready,
	)

	rootrouter.Register(
		router,
		rootController,
		rootMiddleware.Authenticate,
	)

	server :=
		httpserver.New(
			cfg.App.ControlAPIAddr,
			router,
		)

	serverErrors :=
		make(chan error, 1)

	go func() {
		appLogger.Info(
			"control api started",
			zap.String(
				"address",
				cfg.App.ControlAPIAddr,
			),
		)

		serverErrors <- server.Run()
	}()

	shutdownSignal :=
		make(chan os.Signal, 1)

	signal.Notify(
		shutdownSignal,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	select {
	case err := <-serverErrors:
		return err

	case sig := <-shutdownSignal:
		appLogger.Info(
			"shutdown signal received",
			zap.String(
				"signal",
				sig.String(),
			),
		)
	}

	ctx, cancel :=
		context.WithTimeout(
			context.Background(),
			cfg.App.ShutdownTimeout,
		)

	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf(
			"shutdown control api: %w",
			err,
		)
	}

	return nil
}
