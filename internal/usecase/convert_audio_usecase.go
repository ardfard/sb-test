package usecase

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/ardfard/sb-test/internal/domain/entity"
	"github.com/ardfard/sb-test/internal/domain/repository"
	"github.com/ardfard/sb-test/internal/domain/storage"
	"github.com/ardfard/sb-test/internal/infrastructure/converter"
)

type ConvertAudioUseCase struct {
	repo      repository.AudioRepository
	storage   storage.Storage
	converter *converter.AudioConverter
}

func NewConvertAudioUseCase(
	repo repository.AudioRepository,
	storage storage.Storage,
	converter *converter.AudioConverter,
) *ConvertAudioUseCase {
	return &ConvertAudioUseCase{
		repo:      repo,
		storage:   storage,
		converter: converter,
	}
}

func (uc *ConvertAudioUseCase) Convert(ctx context.Context, audioID uint) error {
	audio, err := uc.repo.GetByID(ctx, audioID)
	if err != nil {
		return fmt.Errorf("failed to get audio: %v", err)
	}

	// Update status to converting
	audio.Status = entity.AudioStatusConverting
	if err := uc.repo.Update(ctx, audio); err != nil {
		return fmt.Errorf("failed to update audio status: %v", err)
	}

	// Create temporary files
	inputPath := fmt.Sprintf("/tmp/%d%s", audio.ID, audio.OriginalFormat)
	outputPath := fmt.Sprintf("/tmp/%d.wav", audio.ID)
	defer os.Remove(inputPath)
	defer os.Remove(outputPath)

	// Download original file
	reader, err := uc.storage.Download(ctx, audio.StoragePath)
	if err != nil {
		return uc.handleError(ctx, audio, fmt.Sprintf("failed to download file: %v", err))
	}
	defer reader.Close()

	// Save to temporary file
	inputFile, err := os.Create(inputPath)
	if err != nil {
		return uc.handleError(ctx, audio, fmt.Sprintf("failed to create input file: %v", err))
	}
	defer inputFile.Close()

	if _, err := io.Copy(inputFile, reader); err != nil {
		return uc.handleError(ctx, audio, fmt.Sprintf("failed to save input file: %v", err))
	}

	// Convert the audio
	if err := uc.converter.ConvertToWAV(ctx, inputPath, outputPath); err != nil {
		return uc.handleError(ctx, audio, fmt.Sprintf("failed to convert audio: %v", err))
	}

	// Upload converted file
	outputFile, err := os.Open(outputPath)
	if err != nil {
		return uc.handleError(ctx, audio, fmt.Sprintf("failed to open converted file: %v", err))
	}
	defer outputFile.Close()

	wavPath := fmt.Sprintf("converted/%d.wav", audio.ID)
	if err := uc.storage.Upload(ctx, wavPath, outputFile); err != nil {
		return uc.handleError(ctx, audio, fmt.Sprintf("failed to upload converted file: %v", err))
	}

	// Update audio status to completed
	audio.Status = entity.AudioStatusCompleted
	audio.StoragePath = wavPath
	if err := uc.repo.Update(ctx, audio); err != nil {
		return fmt.Errorf("failed to update audio status: %v", err)
	}

	return nil
}

func (uc *ConvertAudioUseCase) handleError(ctx context.Context, audio *entity.Audio, errMsg string) error {
	audio.Status = entity.AudioStatusFailed
	audio.Error = errMsg
	uc.repo.Update(ctx, audio)
	return fmt.Errorf(errMsg)
}
