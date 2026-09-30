package bootstrap

import (
	"fmt"

	"github.com/dev-gopi/authhub/internal/config"
	"github.com/dev-gopi/authhub/internal/platform/database"
	"github.com/dev-gopi/authhub/internal/platform/rabbitmq"
	redisinfra "github.com/dev-gopi/authhub/internal/platform/redis"
	vaultinfra "github.com/dev-gopi/authhub/internal/platform/vault"
	sharedvalidator "github.com/dev-gopi/authhub/internal/shared/validator"

	"go.uber.org/zap"
)

type Dependencies struct {
	Config *config.Config
	Logger *zap.Logger

	Postgres *database.Postgres
	Redis    *redisinfra.Client
	RabbitMQ *rabbitmq.Connection
	Vault    *vaultinfra.Client

	Validator *sharedvalidator.Validator
}

func BuildDependencies(
	cfg *config.Config,
	logger *zap.Logger,
) (*Dependencies, error) {

	postgresClient, err := database.NewPostgres(
		cfg.Postgres.DSN,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"initialize postgres: %w",
			err,
		)
	}

	redisClient, err := redisinfra.NewClient(
		cfg.Redis.Addr,
		cfg.Redis.Password,
		cfg.Redis.DB,
	)
	if err != nil {
		_ = postgresClient.Close()

		return nil, fmt.Errorf(
			"initialize redis: %w",
			err,
		)
	}

	rabbitConnection, err :=
		rabbitmq.NewConnection(
			cfg.RabbitMQ.URL,
		)

	if err != nil {
		_ = redisClient.Close()
		_ = postgresClient.Close()

		return nil, fmt.Errorf(
			"initialize rabbitmq: %w",
			err,
		)
	}

	vaultClient, err := vaultinfra.NewClient(
		cfg.Vault.Address,
		cfg.Vault.Token,
	)

	if err != nil {
		_ = rabbitConnection.Close()
		_ = redisClient.Close()
		_ = postgresClient.Close()

		return nil, fmt.Errorf(
			"initialize vault: %w",
			err,
		)
	}

	return &Dependencies{
		Config: cfg,
		Logger: logger,

		Postgres: postgresClient,
		Redis:    redisClient,
		RabbitMQ: rabbitConnection,
		Vault:    vaultClient,

		Validator: sharedvalidator.New(),
	}, nil
}

func (d *Dependencies) Close() {
	if d.RabbitMQ != nil {
		_ = d.RabbitMQ.Close()
	}

	if d.Redis != nil {
		_ = d.Redis.Close()
	}

	if d.Postgres != nil {
		_ = d.Postgres.Close()
	}

	if d.Logger != nil {
		_ = d.Logger.Sync()
	}
}
