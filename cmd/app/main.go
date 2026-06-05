package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/topinambur02/url-shortener/internal/config"
	"github.com/topinambur02/url-shortener/internal/db"
	"github.com/topinambur02/url-shortener/internal/handler"
	"github.com/topinambur02/url-shortener/internal/repository"
	"github.com/topinambur02/url-shortener/internal/repository/postgres"
	"github.com/topinambur02/url-shortener/internal/service"
	"github.com/topinambur02/url-shortener/pkg/shutdown"
)

func main() {
	log.Println("Starting application...")

	cfg, err := config.LoadConfig(".env")

	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	} else {
		log.Print("Default configuration loaded successfully")
	}

	storage_type := cfg.App.StorageType
	var repo repository.UrlRepository

	if storage_type == "postgres" {
		log.Println("Initializing database...")
		database, err := db.InitDB(cfg)

		if err != nil {
			log.Fatalf("Error initializing database: %v", err)
		} else {
			log.Println("Database initialized successfully")
		}

		repo = postgres.NewUrlRepository(database)
	} else {
		// TODO: Add in-memory storage
	}

	s := service.NewUrlService(repo)
	h := handler.NewURLHandler(s)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api", h.Create)
	mux.HandleFunc("GET /api/{short}", h.GetByShortUrl)

	host := cfg.App.Host
	port := strconv.Itoa(cfg.App.Port)
	address := host + ":" + port

	server := &http.Server{
		Addr:              address,
		ReadHeaderTimeout: 10 * time.Second,
		Handler: mux,
	}

	go func() {
		log.Printf("Starting HTTP server on %s", address)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Println("Setting up graceful shutdown")
	shutdown.GracefulShutdown([]os.Signal{syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM}, server.Shutdown)
	log.Println("Server exited gracefully")
}
