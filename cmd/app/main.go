package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/topinambur02/url-shortener/internal/config"
	"github.com/topinambur02/url-shortener/pkg/shutdown"
)

func main() {
	cfg, err := config.LoadConfig(".env")

	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	} else {
		log.Print("Default configuration loaded successfully")
	}

	host := cfg.App.Host
	port := strconv.Itoa(cfg.App.Port)
	address := host + ":" + port

	log.Printf("Server will listen on %s", address)

	server := &http.Server{
		Addr:              address,
		ReadHeaderTimeout: 10 * time.Second,
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

