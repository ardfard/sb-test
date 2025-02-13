package usecase

import (
	"audio-processor/internal/domain/entity"
	"audio-processor/internal/domain/repository"
	"audio-processor/internal/infrastructure/converter"
	"audio-processor/internal/infrastructure/storage"
	"audio-processor/pkg/worker"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type AudioUseCase struct {
	repo      repository.AudioRepository
	storage   *storage.GCSStorage
	converter *converter.AudioConverter
	worker    *worker.Worker
}

func NewAudioUseCase(
	repo repository.AudioRepository,
	storage *storage.GCSStorage,
	converter *converter.AudioConverter,
	worker *worker.Worker,
) *AudioUseCase {
	return &AudioUseCase{
		repo:      repo,
		storage:   storage,
		converter: converter,
		worker:    worker,
	}
}

func (uc *AudioUseCase) UploadAudio(ctx context.Context, filename string, content io.Reader) (*entity.Audio, error) {
	audio := &entity.Audio{
		ID:             generateID(),
		OriginalName:   filename,
		OriginalFormat: filepath.Ext(filename),
		Status:         entity.AudioStatusPending,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := uc.repo.Store(ctx, audio); err != nil {
		return nil, fmt.Errorf("failed to store audio metadata: %v", err)
	}

	// Upload original file to GCS
	originalPath := fmt.Sprintf("original/%s%s", audio.ID, audio.OriginalFormat)
	if err := uc.storage.Upload(ctx, originalPath, content); err != nil {
		return nil, fmt.Errorf("failed to upload original file: %v", err)
	}

	// Schedule conversion job
	uc.worker.EnqueueJob(func() {
		uc.convertAudio(context.Background(), audio)
	})

	return audio, nil
}

func (uc *AudioUseCase) convertAudio(ctx context.Context, audio *entity.Audio) {
	audio.Status = entity.AudioStatusConverting
	uc.repo.Update(ctx, audio)

	// Download, convert and upload logic here
	// This is simplified - you'd need to implement temporary file handling
	inputPath := fmt.Sprintf("/tmp/%s%s", audio.ID, audio.OriginalFormat)
	outputPath := fmt.Sprintf("/tmp/%s.wav", audio.ID)

	if err := uc.converter.ConvertToWAV(ctx, inputPath, outputPath); err != nil {
		audio.Status = entity.AudioStatusFailed
		audio.Error = err.Error()
		uc.repo.Update(ctx, audio)
		return
	}

	wavPath := fmt.Sprintf("converted/%s.wav", audio.ID)
	// Upload converted WAV file
	// Update audio status to completed
	audio.Status = entity.AudioStatusCompleted
	audio.StoragePath = wavPath
	uc.repo.Update(ctx, audio)
}

func generateID() string {
	return uuid.New().String()
}
