package usecase

import (
	"context"

	"github.com/ardfard/sb-test/internal/domain/entity"
	"github.com/ardfard/sb-test/internal/domain/repository"
)

type CreatePhraseUseCase struct {
	phraseRepository repository.PhraseRepository
}

func NewCreatePhraseUseCase(phraseRepository repository.PhraseRepository) *CreatePhraseUseCase {
	return &CreatePhraseUseCase{
		phraseRepository: phraseRepository,
	}
}

func (uc *CreatePhraseUseCase) Create(ctx context.Context, text string, userID uint) (*entity.Phrase, error) {
	phrase := &entity.Phrase{
		Phrase: text,
		UserID: userID,
	}

	phrase, err := uc.phraseRepository.Create(ctx, phrase)
	if err != nil {
		return nil, err
	}

	return phrase, nil
}
