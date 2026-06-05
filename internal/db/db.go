package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/topinambur02/url-shortener/internal/config"
)

func InitDB(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	dsn := cfg.DSN
	pool, err := pgxpool.New(ctx, dsn)

	if err != nil {
        return nil, fmt.Errorf("unable to create connection pool: %w", err)
    }

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return pool, nil
}
