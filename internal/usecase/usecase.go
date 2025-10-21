package usecase

import (
	"log/slog"
	"xaxaton/internal/repository"
)

type UC interface{}

type uc struct {
	repo repository.Repository
	log  *slog.Logger
}

func New(repo repository.Repository, log *slog.Logger) *uc {
	return &uc{repo: repo}
}
