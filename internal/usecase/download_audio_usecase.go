package usecase

import (
	"context"
	"fmt"
	"io"

	"github.com/ardfard/sb-test/internal/domain/converter"
	"github.com/ardfard/sb-test/internal/domain/repository"
	"github.com/ardfard/sb-test/internal/domain/storage"
)

type DownloadAudioUseCase struct {
	repo             repository.AudioRepository
	storage          storage.Storage
	converter        converter.AudioConverter
	userRepository   repository.UserRepository
	phraseRepository repository.PhraseRepository
}

func NewDownloadAudioUseCase(
	repo repository.AudioRepository,
	storage storage.Storage,
	converter converter.AudioConverter,
	userRepository repository.UserRepository,
	phraseRepository repository.PhraseRepository,
) *DownloadAudioUseCase {
	return &DownloadAudioUseCase{
		repo:             repo,
		storage:          storage,
		converter:        converter,
		userRepository:   userRepository,
		phraseRepository: phraseRepository,
	}
}

func (uc *DownloadAudioUseCase) Download(ctx context.Context, userID uint, phraseID uint, outputFormat string) (io.ReadCloser, error) {
	// check user and phrase exist
	user, err := uc.userRepository.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	phrase, err := uc.phraseRepository.GetByID(ctx, phraseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get phrase: %v", err)
	}

	if user == nil || phrase == nil {
		return nil, fmt.Errorf("user or phrase not found")
	}

	audio, err := uc.repo.GetByUserIDAndPhraseID(ctx, userID, phraseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get audio: %v", err)
	}

	reader, err := uc.storage.Download(ctx, audio.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("failed to download original file: %v", err)
	}

	// If the requested format is the same as the stored format, return the file without conversion
	if outputFormat == audio.CurrentFormat {
		return reader, nil
	}

	output, err := uc.converter.ConvertFromReader(ctx, reader, audio.CurrentFormat, outputFormat)
	if err != nil {
		reader.Close()
		return nil, fmt.Errorf("failed to convert audio: %v", err)
	}

	return output, nil
}
