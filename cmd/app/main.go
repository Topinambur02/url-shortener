package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/topinambur02/url-shortener/internal/config"
	"github.com/topinambur02/url-shortener/internal/db"
	"github.com/topinambur02/url-shortener/internal/handler"
	"github.com/topinambur02/url-shortener/internal/middleware"
	"github.com/topinambur02/url-shortener/internal/repository"
	inmemory "github.com/topinambur02/url-shortener/internal/repository/in-memory"
	"github.com/topinambur02/url-shortener/internal/repository/postgres"
	"github.com/topinambur02/url-shortener/internal/service"
	"github.com/topinambur02/url-shortener/pkg/logging"
	"github.com/topinambur02/url-shortener/pkg/shutdown"

	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/topinambur02/url-shortener/docs"

	"github.com/rs/cors"
)

// @title           URL shortener (Test Task)
// @version         1.0.0
// @description     Сервис сокращения ссылок
// @host            localhost:8080
// @BasePath        /api
// @schemes         http
func main() {
	logging.Init()
	logger := logging.GetLogger()
	
	logger.Info("=== [START] URL Shortener Application Bootstrap ===")

	storageFlag := flag.String("storage", "", "Тип хранилища: postgres или inmemory")
	flag.Parse()
	logger.Info("Command-line flags parsed successfully")

	logger.Info("Loading configuration fields from .env file...")
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		logger.Fatalf("FATAL: Error loading config file: %v", err)
	}
	logger.Info("Application configuration loaded successfully")

	storageType := *storageFlag
	if storageType == "" {
		logger.Info("Storage flag is empty, falling back to configuration file setting")
		storageType = cfg.App.StorageType
	}
	logger.Infof("Resolved target storage type: '%s'", storageType)

	ctx := context.Background()
	var repo repository.UrlRepository

	switch storageType {
	case "postgres":
		logger.Info("Starting Postgres database initialization...")
		database, err := db.InitDB(ctx, cfg)
		if err != nil {
			logger.Fatalf("FATAL: Database initialization failed: %v", err)
		}
		logger.Info("Database connection established. Creating Postgres repository...")
		repo = postgres.NewUrlRepository(database)
		logger.Info("Postgres repository layer ready")

	case "inmemory":
		logger.Info("Creating In-Memory storage repository...")
		repo = inmemory.NewUrlRepository()
		logger.Info("In-Memory repository layer ready")

	default:
		logger.Fatalf("FATAL: Unknown storage type provided: '%s'. Allowed options: 'postgres' or 'inmemory'", storageType)
	}

	logger.Info("Wiring up application core layers (Service and Handler)...")
	s := service.NewUrlService(repo)
	h := handler.NewURLHandler(s)
	logger.Info("Application core layers wired successfully")

	logger.Info("Configuring HTTP ServeMux and registering API endpoints...")
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api", h.Create)
	mux.HandleFunc("POST /api/{$}", h.Create)
	mux.HandleFunc("GET /api/{short}", h.GetByShortUrl)
	mux.Handle("/docs/", httpSwagger.WrapHandler)
	logger.Info("HTTP routes registered: POST /api, GET /api/{short}, GET /docs/")

	logger.Info("Applying Cross-Origin Resource Sharing (CORS) rules...")
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodOptions},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})
	corsHandler := c.Handler(mux)
	logger.Info("Applying HTTP Logging middleware...")
	handlerStack := middleware.HTTPLoggerMiddleware(corsHandler)

	host := cfg.App.Host
	port := strconv.Itoa(cfg.App.Port)
	address := host + ":" + port

	server := &http.Server{
		Addr:              address,
		ReadHeaderTimeout: 10 * time.Second,
		Handler:           handlerStack,
	}

	logger.Infof("Spawning background routine for HTTP server on %s...", address)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorf("CRITICAL: HTTP server intercepted an error: %s", err)
		}
	}()

	logger.Info("Registering system signal listeners for Graceful Shutdown...")
	err = shutdown.GracefulShutdown(
		[]os.Signal{syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM}, 
		server.Shutdown,
	)
	
	if err != nil {
		logger.Errorf("WARNING: Graceful shutdown encountered an error during cleanup: %v", err)
	}
	
	logger.Info("=== [STOP] Server exited cleanly. Application terminated ===")
}
