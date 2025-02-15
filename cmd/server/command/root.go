package command

import (
	"fmt"
	"log"

	"github.com/ardfard/sb-test/config"
	"github.com/ardfard/sb-test/internal/delivery/http/handler"
	"github.com/ardfard/sb-test/internal/delivery/http/router"
	"github.com/ardfard/sb-test/internal/infrastructure/converter"
	"github.com/ardfard/sb-test/internal/infrastructure/repository"
	"github.com/ardfard/sb-test/internal/infrastructure/storage"
	"github.com/ardfard/sb-test/internal/usecase"
	"github.com/ardfard/sb-test/pkg/worker"
	"github.com/spf13/cobra"

	"context"
	"net/http"
	"os"
	"os/signal"
)

var (
	configPath string
	rootCmd    = &cobra.Command{
		Use:   "sb-test",
		Short: "sb-test is an audio processing server",
		Long:  `A server that processes audio files, converts them to WAV format and stores them`,
		RunE:  run,
	}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "config.yaml", "config file path")
}

func Execute() error {
	return rootCmd.Execute()
}

func run(cmd *cobra.Command, args []string) error {
	// Load configuration using Viper with the specified config path
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("error loading config: %v", err)
	}

	// Initialize SQLite repository using configuration
	repo, err := repository.NewSQLiteAudioRepository(cfg.SQLite.DBPath)
	if err != nil {
		return fmt.Errorf("failed to create repository: %v", err)
	}
	defer repo.Close()

	// Initialize storage using configuration
	storageInstance, err := storage.NewStorage(&cfg.Storage)
	if err != nil {
		return fmt.Errorf("failed to create storage: %v", err)
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
		return fmt.Errorf("HTTP server error: %v", err)
	}

	return nil
}
