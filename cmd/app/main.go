package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/topinambur02/url-shortener/internal/config"
	"github.com/topinambur02/url-shortener/internal/handler"
	"github.com/topinambur02/url-shortener/internal/middleware"
	"github.com/topinambur02/url-shortener/internal/repository"
	"github.com/topinambur02/url-shortener/internal/service"
	"github.com/topinambur02/url-shortener/pkg/logging"
	"github.com/topinambur02/url-shortener/pkg/shutdown"

	_ "github.com/topinambur02/url-shortener/internal/repository/in-memory"
	_ "github.com/topinambur02/url-shortener/internal/repository/postgres"

	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/topinambur02/url-shortener/docs"

	"github.com/rs/cors"

	_ "net/http/pprof"
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

	if err := run(context.Background()); err != nil {
		logger.Fatalf("FATAL: Application failed: %v", err)
	}

	logger.Info("=== [STOP] Server exited cleanly. Application terminated ===")
}

func run(ctx context.Context) error {
	logger := logging.GetLogger()
	storageFlag := flag.String("storage", "", "Тип хранилища: postgres или inmemory")
	flag.Parse()
	logger.Info("Command-line flags parsed successfully")

	cfg, err := config.LoadConfig(".env")
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}
	logger.Info("Application configuration loaded successfully")

	storageType := *storageFlag
	if storageType == "" {
		storageType = cfg.App.StorageType
	}

	repo, err := repository.InitRepository(ctx, storageType, cfg)
	if err != nil {
		return fmt.Errorf("initializing repository: %w", err)
	}

	address := net.JoinHostPort(cfg.App.Host, strconv.Itoa(cfg.App.Port))
	s := service.NewURLService(repo)
	h := handler.NewURLHandler(s, address)

	handlerStack := setupRouter(h)

	server := &http.Server{
		Addr:              address,
		ReadHeaderTimeout: 10 * time.Second,
		Handler:           handlerStack,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	logger.Infof("Spawning background routine for HTTP server on %s...", address)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorf("CRITICAL: HTTP server intercepted an error: %s", err)
		}
	}()

	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()

	logger.Info("Registering system signal listeners for Graceful Shutdown...")
	return shutdown.GracefulShutdown(
		[]os.Signal{syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM},
		server.Shutdown,
	)
}

func setupRouter(h *handler.URLHandler) http.Handler {
	logger := logging.GetLogger()
	logger.Info("Configuring HTTP ServeMux and registering API endpoints...")

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api", h.Create)
	mux.HandleFunc("POST /api/{$}", h.Create)
	mux.HandleFunc("GET /{short}", h.GetByShortURL)
	mux.Handle("/docs/", httpSwagger.WrapHandler)
	logger.Info("HTTP routes registered: POST /api, GET /{short}, GET /docs/")

	logger.Info("Applying Cross-Origin Resource Sharing (CORS) rules...")
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodOptions},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	corsHandler := c.Handler(mux)
	logger.Info("Applying HTTP Logging middleware...")
	return middleware.HTTPLoggerMiddleware(corsHandler)
}
