package main

import (
	"errors"
	"log"
	"os"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/source/file"
	"github.com/golang-migrate/migrate/database/postgres"
	"github.com/topinambur02/url-shortener/internal/config"
	"github.com/topinambur02/url-shortener/internal/db"
)

func main() {
	cfg, err := config.LoadConfig(".env")

	if err != nil {
		log.Fatalln("ERROR: No .env file found")
	}

	log.Println("Starting database migration process")

	database, err := db.InitDB(cfg)

	if err != nil {
		log.Fatalf("ERROR: Database not initialized with error: %v", err)
	}

	if err != nil {
		log.Fatalln("Error getting sql db")
	}

	driver, err := postgres.WithInstance(database, &postgres.Config{})

	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Database connection established successfully")

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres",
		driver,
	)

	if err != nil {
		log.Fatal(err)
	}

	cmd := os.Args[len(os.Args)-1]

	if cmd == "up" {
		if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("Migration up failed: %v", err)
		}
		log.Println("Migration up completed successfully")
	}

	if cmd == "down" {
		if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatalf("Migration down failed: %v", err)
		}
		log.Println("Migration down completed successfully")
	}
}
