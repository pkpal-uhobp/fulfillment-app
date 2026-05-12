package core_postgres_pool

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ConnectionPool struct {
	*pgxpool.Pool
	queryTimeout time.Duration
}

func NewConnectionPool(ctx context.Context, config Config) (*ConnectionPool, error) {
	connectionString := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
		config.SSLMode,
	)

	pgxConfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("parse pgx config: %w", err)
	}

	pgxConfig.MaxConns = config.MaxConns
	pgxConfig.MinConns = config.MinConns
	pgxConfig.MaxConnLifetime = config.MaxConnLifetime
	pgxConfig.MaxConnIdleTime = config.MaxConnIdleTime
	pgxConfig.HealthCheckPeriod = config.HealthCheckPeriod

	connectCtx, cancel := context.WithTimeout(ctx, config.ConnectTimeout)
	defer cancel()

	db, err := pgxpool.NewWithConfig(connectCtx, pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("create pgx pool: %w", err)
	}

	if err := db.Ping(connectCtx); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping pgx pool: %w", err)
	}

	return &ConnectionPool{
		Pool:         db,
		queryTimeout: config.QueryTimeout,
	}, nil
}

func (p *ConnectionPool) QueryTimeout() time.Duration {
	return p.queryTimeout
}
