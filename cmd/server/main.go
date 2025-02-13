package main

import (
	"audio-processor/internal/delivery/http/handler"
	"audio-processor/internal/infrastructure/converter"
	"audio-processor/internal/infrastructure/storage"
	"audio-processor/internal/usecase"
	"audio-processor/pkg/worker"
	"context"
	"log"
	"net/http"
	"os"

	"audio-processor/internal/infrastructure/repository"

	"cloud.google.com/go/storage"
)

func main() {
	ctx := context.Background()

	// Initialize SQLite repository
	repo, err := repository.NewSQLiteAudioRepository("./audio.db")
	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}
	defer repo.Close()

	// Initialize GCS client
	gcsClient, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create GCS client: %v", err)
	}
	defer gcsClient.Close()

	bucketName := os.Getenv("GCS_BUCKET_NAME")
	gcsStorage := storage.NewGCSStorage(gcsClient, bucketName)

	// Initialize components
	converter := converter.NewAudioConverter()
	worker := worker.NewWorker(5) // 5 concurrent workers

	// Initialize use case
	useCase := usecase.NewAudioUseCase(repo, gcsStorage, converter, worker)

	// Initialize handler
	audioHandler := handler.NewAudioHandler(useCase)

	// Setup routes
	http.HandleFunc("/upload", audioHandler.UploadAudio)

	// Start server
	log.Fatal(http.ListenAndServe(":8080", nil))
}
