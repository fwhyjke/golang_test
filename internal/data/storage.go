package data

import (
	"sync"
	"sync/atomic"
)

type Note struct {
	ID          uint64
	Title       string
	Description string
	Done        bool
}
type IDGenerator struct {
	counter atomic.Uint64
}

func (g *IDGenerator) NextID() uint64 {
	return g.counter.Add(1)
}

type DataBase struct {
	mu    sync.RWMutex
	notes map[uint64]Note
	idGen atomic.Uint64
}

func NewDataBase() *DataBase {
	return &DataBase{
		notes: make(map[uint64]Note),
	}
}

func (db *DataBase) Create(title, desc string) Note {
	note := Note{
		ID:          db.idGen.Add(1),
		Title:       title,
		Description: desc,
		Done:        false,
	}

	db.mu.Lock()
	defer db.mu.Unlock()
	db.notes[note.ID] = note
	return note
}

func (db *DataBase) GetAll() []Note {
	db.mu.RLock()
	defer db.mu.RUnlock()

	res := make([]Note, 0, len(db.notes))
	for _, n := range db.notes {
		res = append(res, n)
	}
	return res
}

func (db *DataBase) GetByID(id uint64) (Note, bool) {
	db.mu.RLock()
	note, ok := db.notes[id]
	db.mu.RUnlock()
	return note, ok
}

func (db *DataBase) Update(id uint64, title, desc string, done bool) (Note, bool) {
	db.mu.Lock()
	defer db.mu.Unlock()

	n, ok := db.notes[id]
	if !ok {
		return Note{}, false
	}

	n.Title = title
	n.Description = desc
	n.Done = done
	db.notes[id] = n

	return n, true
}

func (db *DataBase) Delete(id uint64) bool {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, ok := db.notes[id]; !ok {
		return false
	}
	delete(db.notes, id)
	return true
}
