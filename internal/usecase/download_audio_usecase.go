package usecase

import (
	"context"
	"fmt"
	"io"

	"github.com/ardfard/sb-test/internal/domain/entity"
	"github.com/ardfard/sb-test/internal/domain/repository"
	"github.com/ardfard/sb-test/internal/domain/storage"
)

type DownloadAudioUseCase struct {
	repo    repository.AudioRepository
	storage storage.Storage
}

func NewDownloadAudioUseCase(
	repo repository.AudioRepository,
	storage storage.Storage,
) *DownloadAudioUseCase {
	return &DownloadAudioUseCase{
		repo:    repo,
		storage: storage,
	}
}

func (uc *DownloadAudioUseCase) Download(ctx context.Context, audioID uint) (io.ReadCloser, error) {
	audio, err := uc.repo.GetByID(ctx, audioID)
	if err != nil {
		return nil, fmt.Errorf("failed to get audio: %v", err)
	}

	if audio.Status != entity.AudioStatusCompleted {
		return nil, fmt.Errorf("audio is not ready for download, current status: %s", audio.Status)
	}

	reader, err := uc.storage.Download(ctx, audio.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %v", err)
	}

	return reader, nil
}
