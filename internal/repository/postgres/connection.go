package postgres

import (
	"context"
	"fmt"

	"github.com/fwhyjke/golang_test/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	Host     string `env:"POSTGRES_HOST" required:"true"`
	Port     string `env:"POSTGRES_PORT" required:"true"`
	User     string `env:"POSTGRES_USER" required:"true"`
	Password string `env:"POSTGRES_PASSWORD" required:"true"`
	Database string `env:"POSTGRES_DB" required:"true"`
}

func NewPostgresConfig() (Config, error) {
	var cfg Config
	if err := envconfig.Process(context.Background(), &cfg); err != nil {
		return Config{}, fmt.Errorf("envconfig for postgres: %w", err)
	}
	return cfg, nil
}

type ConnectionPool struct {
	*pgxpool.Pool
}

func NewPostgresConnection(cfg Config) repository.NoteRepository {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database,
	)

	poolCfg, err := pgxpool.ParseConfig(connString)

	if err != nil {
		panic(fmt.Errorf("pgxpool config: %w", err))
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolCfg)

	if err != nil {
		panic(fmt.Errorf("pgxpool connect: %w", err))
	}

	if err := pool.Ping(context.Background()); err != nil {
		panic(fmt.Errorf("pgxpool ping: %w", err))
	}

	return &ConnectionPool{pool}
}
