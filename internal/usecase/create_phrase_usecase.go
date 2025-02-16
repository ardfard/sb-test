package usecase

import (
	"context"

	"github.com/ardfard/sb-test/internal/domain/entity"
	"github.com/ardfard/sb-test/internal/domain/repository"
)

type CreatePhraseUsecase struct {
	phraseRepository repository.PhraseRepository
}

func NewCreatePhraseUsecase(phraseRepository repository.PhraseRepository) *CreatePhraseUsecase {
	return &CreatePhraseUsecase{phraseRepository: phraseRepository}
}

func (u *CreatePhraseUsecase) Execute(ctx context.Context, phrase *entity.Phrase) error {
	return u.phraseRepository.Create(ctx, phrase)
}
