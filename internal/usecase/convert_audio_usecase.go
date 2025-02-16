package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/ardfard/sb-test/internal/domain/converter"
	"github.com/ardfard/sb-test/internal/domain/entity"
	"github.com/ardfard/sb-test/internal/domain/repository"
	"github.com/ardfard/sb-test/internal/domain/storage"
)

type ConvertAudioUseCase struct {
	repo      repository.AudioRepository
	storage   storage.Storage
	converter converter.AudioConverter
}

func NewConvertAudioUseCase(
	repo repository.AudioRepository,
	storage storage.Storage,
	converter converter.AudioConverter,
) *ConvertAudioUseCase {
	return &ConvertAudioUseCase{
		repo:      repo,
		storage:   storage,
		converter: converter,
	}
}

const targetFormat = "wav"

func (uc *ConvertAudioUseCase) Convert(ctx context.Context, audioID uint) error {
	audio, err := uc.repo.GetByID(ctx, audioID)
	if err != nil {
		return fmt.Errorf("failed to get audio: %v", err)
	}

	originalPath := audio.StoragePath
	// Update status to converting
	audio.Status = entity.AudioStatusConverting
	if err := uc.repo.Update(ctx, audio); err != nil {
		return fmt.Errorf("failed to update audio status: %v", err)
	}

	// Download original file
	reader, err := uc.storage.Download(ctx, originalPath)
	if err != nil {
		return uc.handleError(ctx, audio, fmt.Sprintf("failed to download file: %v", err))
	}
	defer reader.Close()

	// Convert the audio
	output, err := uc.converter.ConvertFromReader(ctx, reader, audio.CurrentFormat, targetFormat)
	if err != nil {
		return uc.handleError(ctx, audio, fmt.Sprintf("failed to convert audio: %v", err))
	}
	defer output.Close()

	// Upload converted file
	convertedPath := fmt.Sprintf("%s/converted/%d.%s", basePath, audio.ID, targetFormat)
	if err := uc.storage.Upload(ctx, convertedPath, output); err != nil {
		return uc.handleError(ctx, audio, fmt.Sprintf("failed to upload converted file: %v", err))
	}

	// Update audio status to completed
	audio.Status = entity.AudioStatusCompleted
	audio.StoragePath = convertedPath
	audio.CurrentFormat = targetFormat
	if err := uc.repo.Update(ctx, audio); err != nil {
		return fmt.Errorf("failed to update audio status: %v", err)
	}

	// delete the original file
	if err := uc.storage.Delete(ctx, originalPath); err != nil {
		return uc.handleError(ctx, audio, fmt.Sprintf("failed to delete original file: %v", err))
	}

	return nil
}

func (uc *ConvertAudioUseCase) handleError(ctx context.Context, audio *entity.Audio, errMsg string) error {
	audio.Status = entity.AudioStatusFailed
	audio.Error = errMsg
	if err := uc.repo.Update(ctx, audio); err != nil {
		return fmt.Errorf("failed to update audio status: %v", err)
	}
	return errors.New(errMsg)
}
