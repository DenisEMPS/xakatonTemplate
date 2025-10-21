package usecase

import (
	"log/slog"
	"xaxaton/internal/repository"
	"xaxaton/internal/tokenizer"
)

type UC interface{}

type uc struct {
	repo      repository.Repository
	tokenizer tokenizer.Tokenizer
	log       *slog.Logger
}

func New(repo repository.Repository, tokenizer tokenizer.Tokenizer, log *slog.Logger) *uc {
	return &uc{
		repo:      repo,
		tokenizer: tokenizer,
		log:       log,
	}
}
