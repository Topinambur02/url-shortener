package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"

	"github.com/topinambur02/url-shortener/internal/config"
)

func InitDB(cfg *config.Config) (*sql.DB, error) {
	dsn := cfg.DSN
	db, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}
