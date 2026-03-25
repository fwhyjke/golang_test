package inmemory

import (
	"context"
	"sync"
	"sync/atomic"

	"github.com/fwhyjke/golang_test/internal/repository"
)

type InMemoryDataBase struct {
	mu    sync.RWMutex
	notes map[uint64]repository.Note
	idGen atomic.Uint64
}

func NewInMemoryDataBase() repository.NoteRepository {
	return &InMemoryDataBase{
		notes: make(map[uint64]repository.Note),
	}
}

func (db *InMemoryDataBase) Delete(ctx context.Context, id uint64) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	if _, ok := db.notes[id]; !ok {
		return repository.ErrNotFoundID
	}

	delete(db.notes, id)
	return nil
}

func (db *InMemoryDataBase) GetByID(ctx context.Context, id uint64) (repository.Note, error) {
	select {
	case <-ctx.Done():
		return repository.Note{}, ctx.Err()
	default:
	}

	db.mu.RLock()
	defer db.mu.RUnlock()

	note, ok := db.notes[id]
	if !ok {
		return repository.Note{}, repository.ErrNotFoundID
	}

	return note, nil
}

func (db *InMemoryDataBase) GetAll(ctx context.Context) ([]repository.Note, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	db.mu.RLock()
	defer db.mu.RUnlock()

	res := make([]repository.Note, 0, len(db.notes))
	for _, n := range db.notes {
		res = append(res, n)
	}

	return res, nil
}

func (db *InMemoryDataBase) Update(ctx context.Context, id uint64, dto repository.NoteDTO) (repository.Note, error) {
	select {
	case <-ctx.Done():
		return repository.Note{}, ctx.Err()
	default:
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	n, ok := db.notes[id]
	if !ok {
		return repository.Note{}, repository.ErrNotFoundID
	}

	if dto.Title == "" {
		return repository.Note{}, repository.ErrTitleNotDefined
	}
	n.Title = dto.Title
	n.Description = dto.Description
	n.Done = dto.Done

	db.notes[id] = n
	return n, nil
}

func (db *InMemoryDataBase) Create(ctx context.Context, dto repository.NoteDTO) (repository.Note, error) {
	select {
	case <-ctx.Done():
		return repository.Note{}, ctx.Err()
	default:
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	if dto.Title == "" {
		return repository.Note{}, repository.ErrTitleNotDefined
	}

	note := repository.Note{
		ID:          db.idGen.Add(1),
		Title:       dto.Title,
		Description: dto.Description,
		Done:        dto.Done,
	}

	db.notes[note.ID] = note
	return note, nil
}
