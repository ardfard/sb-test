package main

import (
	"context"
	"log"
	"net/http"

	"github.com/ardfard/sb-test/config"
	"github.com/ardfard/sb-test/internal/delivery/http/handler"
	"github.com/ardfard/sb-test/internal/delivery/http/router"
	"github.com/ardfard/sb-test/internal/infrastructure/converter"
	"github.com/ardfard/sb-test/internal/infrastructure/repository"
	"github.com/ardfard/sb-test/internal/infrastructure/storage"
	"github.com/ardfard/sb-test/internal/usecase"
	"github.com/ardfard/sb-test/pkg/worker"

	gcsStorage "cloud.google.com/go/storage"
)

func main() {
	// Load configuration using Viper
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	ctx := context.Background()

	// Initialize SQLite repository using configuration
	repo, err := repository.NewSQLiteAudioRepository(cfg.SQLite.DBPath)
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	// Initialize GCS client from configuration
	gcsClient, err := gcsStorage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create GCS client: %v", err)
	}
	defer gcsClient.Close()
	gcsStorageInstance := storage.NewGCSStorage(gcsClient, cfg.GCS.Bucket)

	// Initialize components using configuration values
	converterInstance := converter.NewAudioConverter()
	workerInstance := worker.NewWorker(cfg.Worker.NumWorkers)

	// Initialize use case with our components
	useCase := usecase.NewAudioUseCase(repo, gcsStorageInstance, converterInstance, workerInstance)

	// Initialize handler
	audioHandler := handler.NewAudioHandler(useCase)

	// Initialize the router with all defined routes.
	r := router.SetupRoutes(audioHandler)

	// Start server using address from configuration
	log.Fatal(http.ListenAndServe(cfg.ServerAddress, r))
}
