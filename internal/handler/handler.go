package handler

import (
	"xaxaton/internal/usecase"
)

type Handler struct {
	uc usecase.UC
}

func New(uc usecase.UC) *Handler {
	return &Handler{uc: uc}
}
