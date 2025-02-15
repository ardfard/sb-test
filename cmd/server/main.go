package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/ardfard/sb-test/config"
	"github.com/ardfard/sb-test/internal/delivery/http/handler"
	"github.com/ardfard/sb-test/internal/delivery/http/router"
	"github.com/ardfard/sb-test/internal/infrastructure/converter"
	"github.com/ardfard/sb-test/internal/infrastructure/repository"
	"github.com/ardfard/sb-test/internal/infrastructure/storage"
	"github.com/ardfard/sb-test/internal/usecase"
	"github.com/ardfard/sb-test/pkg/worker"
)

func main() {
	// Load configuration using Viper
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Initialize SQLite repository using configuration
	repo, err := repository.NewSQLiteAudioRepository(cfg.SQLite.DBPath)
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	// Initialize storage using configuration
	storageInstance, err := storage.NewStorage(&cfg.Storage)
	if err != nil {
		log.Fatalf("Failed to create storage: %v", err)
	}

	// Initialize other components using configuration values.
	converterInstance := converter.NewAudioConverter()
	workerInstance := worker.NewWorker(cfg.Worker.NumWorkers)

	// Initialize use case with our components.
	useCase := usecase.NewAudioUseCase(repo, storageInstance, converterInstance, workerInstance)

	// Initialize handler.
	audioHandler := handler.NewAudioHandler(useCase)

	// Initialize the router with all defined routes.
	r := router.SetupRoutes(audioHandler)

	// Create server
	srv := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: r,
	}

	// Handle graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt)
		<-sigChan
		log.Println("Received interrupt signal, shutting down...")
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}
	}()

	// Start server
	log.Printf("Starting server on %s", cfg.ServerAddress)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server error: %v", err)
	}
}
