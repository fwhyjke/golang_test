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

type NoteDTO struct {
	Title       string `json:"title"`
	Description *string `json:"description"`
	Done        *bool  `json:"done"`
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
