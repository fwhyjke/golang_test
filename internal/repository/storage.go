package repository

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

type NoteRepository interface {
	Create(ctx context.Context, dto NoteDTO) (Note, error)
	GetByID(ctx context.Context, id uint64) (Note, error)
	GetAll(ctx context.Context) ([]Note, error)
	Update(ctx context.Context, id uint64, dto NoteDTO) (Note, error)
	Delete(ctx context.Context, id uint64) error
}

type Note struct {
	ID          uint64
	Title       string
	Description string
	Done        bool
}

type NoteDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

type InMemoryDataBase struct {
	mu    sync.RWMutex
	notes map[uint64]Note
	idGen atomic.Uint64
}

func NewInMemoryDataBase() *InMemoryDataBase {
	return &InMemoryDataBase{
		notes: make(map[uint64]Note),
	}
}

var ErrNotFoundID error = errors.New("note by ID not found")
var ErrTitleNotDefined error = errors.New("title is required field")
