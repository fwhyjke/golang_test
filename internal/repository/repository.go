package repository

import (
	"context"
	"errors"
)

type NoteRepository interface {
	Create(ctx context.Context, dto NoteDTO) (Note, error)
	GetByID(ctx context.Context, id uint64) (Note, error)
	GetAll(ctx context.Context) ([]Note, error)
	Update(ctx context.Context, id uint64, dto NoteDTO) (Note, error)
	Delete(ctx context.Context, id uint64) error
}

type Note struct {
	ID          uint64 `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

type NoteDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

var ErrNotFoundID error = errors.New("note by ID not found")
var ErrTitleNotDefined error = errors.New("title is required field")
