package data

import (
	"context"
)

func (db *DataBase) Delete(ctx context.Context, id uint64) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	if _, ok := db.notes[id]; !ok {
		return ErrNotFoundID
	}

	delete(db.notes, id)
	return nil
}

func (db *DataBase) GetByID(ctx context.Context, id uint64) (Note, error) {
	select {
	case <-ctx.Done():
		return Note{}, ctx.Err()
	default:
	}

	db.mu.RLock()
	defer db.mu.RUnlock()

	note, ok := db.notes[id]
	if !ok {
		return Note{}, ErrNotFoundID
	}

	return note, nil
}

func (db *DataBase) GetAll(ctx context.Context) ([]Note, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	db.mu.RLock()
	defer db.mu.RUnlock()

	res := make([]Note, 0, len(db.notes))
	for _, n := range db.notes {
		res = append(res, n)
	}

	return res, nil
}

func (db *DataBase) Update(ctx context.Context, id uint64, dto NoteDTO) (Note, error) {
	select {
	case <-ctx.Done():
		return Note{}, ctx.Err()
	default:
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	n, ok := db.notes[id]
	if !ok {
		return Note{}, ErrNotFoundID
	}

	n.Title = dto.Title
	n.Description = dto.Description
	n.Done = dto.Done

	db.notes[id] = n
	return n, nil
}

func (db *DataBase) Create(ctx context.Context, dto NoteDTO) (Note, error) {
	select {
	case <-ctx.Done():
		return Note{}, ctx.Err()
	default:
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	note := Note{
		ID:          db.idGen.Add(1),
		Title:       dto.Title,
		Description: dto.Description,
		Done:        dto.Done,
	}

	db.notes[note.ID] = note
	return note, nil
}
