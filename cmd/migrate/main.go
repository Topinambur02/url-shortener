package main

import (
	"context"
	"errors"
	"os"

	"github.com/golang-migrate/migrate/v4"
	pgxMigrate "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/stdlib"

	"github.com/topinambur02/url-shortener/internal/config"
	"github.com/topinambur02/url-shortener/internal/db"
	"github.com/topinambur02/url-shortener/pkg/logging"
)

func main() {
	logging.Init()
	logger := logging.GetLogger()
	logger.Info("Starting database migration process...")

	logger.Info("Loading configuration file...")
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		logger.Errorf("Failed to load config from .env: %v\n", err)
		return
	}
	logger.Info("Configuration loaded successfully")

	logger.Info("Initializing database connection...")
	ctx := context.Background()
	database, err := db.InitDB(ctx, cfg)
	if err != nil {
		logger.Errorf("Database initialization failed: %v\n", err)
		return
	}

	logger.Info("Creating database driver instance for postgres...")
	sqlDB := stdlib.OpenDBFromPool(database)
	driver, err := pgxMigrate.WithInstance(sqlDB, &pgxMigrate.Config{})
	if err != nil {
		logger.Errorf("Failed to create postgres driver instance: %v\n", err)
	}
	logger.Info("Database driver instance created successfully")

	logger.Info("Initializing migrate instance with source 'file://migrations'...")
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)
	if err != nil {
		logger.Errorf("Failed to initialize migrate instance: %v\n", err)
	}

	if len(os.Args) < 2 {
		logger.Errorln("Missing command argument. Please specify 'up' or 'down'")
	}

	cmd := os.Args[len(os.Args)-1]
	logger.Printf("Detected migration command: '%s'\n", cmd)

	switch cmd {
	case "up":
		logger.Info("Executing 'up' migrations...")
		if err := m.Up(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				logger.Info("No migration changes to apply. Schema is already up to date.")
			} else {
				logger.Errorf("Migration 'up' failed: %v\n", err)
			}
		} else {
			logger.Info("Migration 'up' completed successfully")
		}

	case "down":
		logger.Info("Executing 'down' migrations...")
		if err := m.Down(); err != nil {
			if errors.Is(err, migrate.ErrNoChange) {
				logger.Info("No migration changes to revert. Database is already at the base state.")
			} else {
				logger.Errorf("Migration 'down' failed: %v\n", err)
			}
		} else {
			logger.Info("Migration 'down' completed successfully")
		}

	default:
		logger.Errorf("Unknown migration command '%s'. Allowed commands are 'up' or 'down'\n", cmd)
	}

	logger.Info("Migration process finished successfully")
}
