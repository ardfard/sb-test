package usecase

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/ardfard/sb-test/internal/domain/entity"
	"github.com/ardfard/sb-test/internal/domain/queue"
	"github.com/ardfard/sb-test/internal/domain/repository"
	"github.com/ardfard/sb-test/internal/domain/storage"
)

type UploadAudioUseCase struct {
	repo    repository.AudioRepository
	storage storage.Storage
	queue   queue.TaskQueue
}

func NewUploadAudioUseCase(
	repo repository.AudioRepository,
	storage storage.Storage,
	queue queue.TaskQueue,
) *UploadAudioUseCase {
	return &UploadAudioUseCase{
		repo:    repo,
		storage: storage,
		queue:   queue,
	}
}

func (uc *UploadAudioUseCase) Upload(ctx context.Context, filename string, content io.Reader) (*entity.Audio, error) {
	audio := &entity.Audio{
		OriginalName:   filename,
		OriginalFormat: filepath.Ext(filename),
		Status:         entity.AudioStatusPending,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		UserID:         123, // TODO: Get from context
		PhraseID:       456, // TODO: Get from context
	}

	if err := uc.repo.Store(ctx, audio); err != nil {
		return nil, fmt.Errorf("failed to store audio metadata: %v", err)
	}

	// Upload original file to storage
	originalPath := fmt.Sprintf("original/%d%s", audio.ID, audio.OriginalFormat)
	if err := uc.storage.Upload(ctx, originalPath, content); err != nil {
		return nil, fmt.Errorf("failed to upload original file: %v", err)
	}

	// Enqueue conversion task
	if err := uc.queue.Enqueue(ctx, "audio_conversion", audio.ID); err != nil {
		return nil, fmt.Errorf("failed to enqueue conversion: %v", err)
	}

	return audio, nil
}
