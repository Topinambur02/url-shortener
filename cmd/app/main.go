package main

import (
	"flag"
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
	inmemory "github.com/topinambur02/url-shortener/internal/repository/in-memory"
	"github.com/topinambur02/url-shortener/internal/repository/postgres"
	"github.com/topinambur02/url-shortener/internal/service"
	"github.com/topinambur02/url-shortener/pkg/shutdown"

	_ "github.com/topinambur02/url-shortener/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           URL shortener (Test Task)
// @version         1.0.0
// @description     Сервис сокращения ссылок
// @host            localhost:8080
// @BasePath        /api
// @schemes         http
func main() {
	storageFlag := flag.String("storage", "", "Тип хранилища: postgres или inmemory")
	flag.Parse()
	log.Println("Starting application...")

	cfg, err := config.LoadConfig(".env")

	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	} else {
		log.Print("Default configuration loaded successfully")
	}

	storageType := *storageFlag
	if storageType == "" {
		storageType = cfg.App.StorageType
	}

	var repo repository.UrlRepository

	switch storageType {
	case "postgres":
		log.Println("Initializing database...")
		database, err := db.InitDB(cfg)
		if err != nil {
			log.Fatalf("Error initializing database: %v", err)
		}
		repo = postgres.NewUrlRepository(database)
		log.Println("Database initialized successfully")

	case "inmemory":
		log.Println("Initializing in-memory storage...")
		repo = inmemory.NewUrlRepository()

	default:
		log.Fatalf("Unknown storage type: %s. Use 'postgres' or 'inmemory'", storageType)
	}

	s := service.NewUrlService(repo)
	h := handler.NewURLHandler(s)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api", h.Create)
	mux.HandleFunc("POST /api/{$}", h.Create)
	mux.HandleFunc("GET /api/{short}", h.GetByShortUrl)
	mux.Handle("/docs/", httpSwagger.WrapHandler)

	host := cfg.App.Host
	port := strconv.Itoa(cfg.App.Port)
	address := host + ":" + port

	server := &http.Server{
		Addr:              address,
		ReadHeaderTimeout: 10 * time.Second,
		Handler:           mux,
	}

	go func() {
		log.Printf("Starting HTTP server on %s", address)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Println("Setting up graceful shutdown")
	err = shutdown.GracefulShutdown([]os.Signal{syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM}, server.Shutdown)
	if err != nil {
		log.Fatalf("Shutdown error: %v", err)
	}
	log.Println("Server exited gracefully")
}
