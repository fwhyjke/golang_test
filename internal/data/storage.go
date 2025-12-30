package data

import (
	"errors"
	"sync"
	"sync/atomic"
)

type Note struct {
	ID          uint64
	Title       string
	Description string
	Done        bool
}

type NoteDTO struct {
	Title       string  `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
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

var ErrNotFoundID error = errors.New("note by ID not found")
