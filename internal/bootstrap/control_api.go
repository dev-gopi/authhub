package bootstrap

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dev-gopi/authhub/internal/config"
	auditrepository "github.com/dev-gopi/authhub/internal/modules/audit/repository"
	keymanagementservice "github.com/dev-gopi/authhub/internal/modules/keymanagement/service"
	outboxrepository "github.com/dev-gopi/authhub/internal/modules/outbox/repository"
	permissionrepository "github.com/dev-gopi/authhub/internal/modules/permission/repository"
	platformuserrepository "github.com/dev-gopi/authhub/internal/modules/platformuser/repository"
	rolerepository "github.com/dev-gopi/authhub/internal/modules/role/repository"
	rootcontroller "github.com/dev-gopi/authhub/internal/modules/rootauth/controller"
	rootmiddleware "github.com/dev-gopi/authhub/internal/modules/rootauth/middleware"
	rootrepository "github.com/dev-gopi/authhub/internal/modules/rootauth/repository"
	rootrouter "github.com/dev-gopi/authhub/internal/modules/rootauth/router"
	rootservice "github.com/dev-gopi/authhub/internal/modules/rootauth/service"
	rootvalidator "github.com/dev-gopi/authhub/internal/modules/rootauth/validator"
	tenantcontroller "github.com/dev-gopi/authhub/internal/modules/tenant/controller"
	tenantrepository "github.com/dev-gopi/authhub/internal/modules/tenant/repository"
	tenantrouter "github.com/dev-gopi/authhub/internal/modules/tenant/router"
	tenantservice "github.com/dev-gopi/authhub/internal/modules/tenant/service"
	tenantvalidator "github.com/dev-gopi/authhub/internal/modules/tenant/validator"
	tenantmemberrepository "github.com/dev-gopi/authhub/internal/modules/tenantmember/repository"
	userpoolrepository "github.com/dev-gopi/authhub/internal/modules/userpool/repository"
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
		return fmt.Errorf("load config: %w", err)
	}

	appLogger, err := logger.New(cfg.App.Environment)
	if err != nil {
		return fmt.Errorf("create logger: %w", err)
	}

	deps, err := BuildDependencies(cfg, appLogger)
	if err != nil {
		_ = appLogger.Sync()
		return err
	}
	defer deps.Close()

	rootUserRepo := rootrepository.NewPlatformUserRepository(deps.Postgres.DB)
	passwordRepo := rootrepository.NewPasswordRepository(deps.Postgres.DB)
	sessionRepo := rootrepository.NewSessionRepository(deps.Postgres.DB)
	loginAttemptRepo := rootrepository.NewLoginAttemptRepository(deps.Postgres.DB)
	passwordHasher := security.NewPasswordHasher()

	rootService, err := rootservice.New(
		deps.Postgres.DB, deps.Redis, rootUserRepo, passwordRepo, sessionRepo,
		loginAttemptRepo, passwordHasher, appLogger,
	)
	if err != nil {
		return fmt.Errorf("initialize root auth service: %w", err)
	}

	rootValidator := rootvalidator.New(deps.Validator)
	rootController := rootcontroller.New(rootService, rootValidator)
	rootMiddleware := rootmiddleware.New(rootService)

	keyManagementService := keymanagementservice.New(deps.Vault, cfg.Vault.TransitMount, appLogger)
	tenantRepo := tenantrepository.NewPostgresRepository(deps.Postgres.DB)
	tenantProfileRepo := tenantrepository.NewProfileRepository()
	userPoolRepo := userpoolrepository.NewPostgresRepository()
	permissionRepo := permissionrepository.NewPostgresRepository()
	roleRepo := rolerepository.NewPostgresRepository()
	platformUserRepo := platformuserrepository.NewPostgresRepository()
	tenantMemberRepo := tenantmemberrepository.NewPostgresRepository()
	outboxRepo := outboxrepository.NewPostgresRepository()
	auditRepo := auditrepository.NewPostgresRepository()

	tenantService := tenantservice.New(
		deps.Postgres.DB,
		tenantRepo,
		tenantProfileRepo,
		userPoolRepo,
		permissionRepo,
		roleRepo,
		platformUserRepo,
		tenantMemberRepo,
		passwordRepo,
		outboxRepo,
		auditRepo,
		passwordHasher,
		keyManagementService,
		cfg.App.IssuerBaseURL,
		cfg.App.DefaultRoleTemplateVersion,
		appLogger,
	)
	tenantValidator := tenantvalidator.New(deps.Validator)
	tenantController := tenantcontroller.New(tenantService, tenantValidator)

	router := chi.NewRouter()
	router.Use(sharedmiddleware.RequestID)
	router.Use(sharedmiddleware.CorrelationID)
	router.Use(sharedmiddleware.Recovery(appLogger))
	router.Use(sharedmiddleware.RequestLogger(appLogger))
	router.Use(sharedmiddleware.SecurityHeadersMiddleware)
	router.Use(sharedmiddleware.Timeout(cfg.App.RequestTimeout))

	health := httpserver.NewHealthHandler(deps.Postgres, deps.Redis, deps.RabbitMQ, deps.Vault, "control-api")
	router.Get("/health", health.Health)
	router.Get("/ready", health.Ready)

	rootrouter.Register(router, rootController, rootMiddleware.Authenticate)
	tenantrouter.Register(router, tenantController, rootMiddleware.Authenticate)

	server := httpserver.New(cfg.App.ControlAPIAddr, router)
	serverErrors := make(chan error, 1)
	go func() {
		appLogger.Info("control api started", zap.String("address", cfg.App.ControlAPIAddr))
		serverErrors <- server.Run()
	}()

	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM)
	select {
	case err := <-serverErrors:
		return err
	case sig := <-shutdownSignal:
		appLogger.Info("shutdown signal received", zap.String("signal", sig.String()))
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.App.ShutdownTimeout)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown control api: %w", err)
	}
	return nil
}
