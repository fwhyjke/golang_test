package repository

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	repo := NewInMemoryDataBase()

	testTable := []struct {
		name     string
		ctx      context.Context
		dto      NoteDTO
		expErr   error
		expID    uint64
		expTitle string
		expDesc  string
		expDone  bool
	}{
		{
			name: "full dto",
			ctx:  context.Background(),
			dto: NoteDTO{
				Title:       "title",
				Description: "desc",
				Done:        true,
			},
			expErr:   nil,
			expID:    1,
			expTitle: "title",
			expDesc:  "desc",
			expDone:  true,
		},
		{
			name: "minimal dto",
			ctx:  context.Background(),
			dto: NoteDTO{
				Title: "title2",
			},
			expErr:   nil,
			expID:    2,
			expTitle: "title2",
			expDesc:  "",
			expDone:  false,
		},
		{
			name:   "empty title",
			ctx:    context.Background(),
			dto:    NoteDTO{},
			expErr: ErrTitleNotDefined,
		},
		{
			name: "context canceled",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			dto:    NoteDTO{Title: "123"},
			expErr: context.Canceled,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			note, err := repo.Create(testCase.ctx, testCase.dto)

			if !errors.Is(err, testCase.expErr) {
				t.Fatalf("expected error %v, got %v", testCase.expErr, err)
			}

			if testCase.expErr != nil {
				return
			}

			if note.ID != testCase.expID {
				t.Errorf("id: expected %d, got %d", testCase.expID, note.ID)
			}
			if note.Title != testCase.expTitle {
				t.Errorf("title: expected %q, got %q", testCase.expTitle, note.Title)
			}
			if note.Description != testCase.expDesc {
				t.Errorf("description: expected %q, got %q", testCase.expDesc, note.Description)
			}
			if note.Done != testCase.expDone {
				t.Errorf("done: expected %v, got %v", testCase.expDone, note.Done)
			}
		})
	}
}

func TestGetByID(t *testing.T) {
	repo := NewInMemoryDataBase()
	created, _ := repo.Create(context.Background(), NoteDTO{Title: "test"})

	testTable := []struct {
		name   string
		ctx    context.Context
		id     uint64
		expErr error
	}{
		{
			name:   "success",
			ctx:    context.Background(),
			id:     created.ID,
			expErr: nil,
		},
		{
			name:   "not found",
			ctx:    context.Background(),
			id:     123,
			expErr: ErrNotFoundID,
		},
		{
			name: "context canceled",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			id:     created.ID,
			expErr: context.Canceled,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := repo.GetByID(testCase.ctx, testCase.id)

			if !errors.Is(err, testCase.expErr) {
				t.Fatalf("expected error %v, got %v", testCase.expErr, err)
			}
		})
	}
}

func TestGetAll(t *testing.T) {
	repo := NewInMemoryDataBase()
	repo.Create(context.Background(), NoteDTO{Title: "t1"})
	repo.Create(context.Background(), NoteDTO{Title: "t2"})

	testTable := []struct {
		name      string
		ctx       context.Context
		expErr    error
		expLength int
	}{
		{
			name:      "success",
			ctx:       context.Background(),
			expErr:    nil,
			expLength: 2,
		},
		{
			name: "context canceled",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			expErr: context.Canceled,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			notes, err := repo.GetAll(testCase.ctx)

			if !errors.Is(err, testCase.expErr) {
				t.Fatalf("expected error %v, got %v", testCase.expErr, err)
			}

			if testCase.expErr == nil && len(notes) != testCase.expLength {
				t.Fatalf("expected %d notes, got %d", testCase.expLength, len(notes))
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	repo := NewInMemoryDataBase()
	created, _ := repo.Create(context.Background(), NoteDTO{
		Title:       "title",
		Description: "desc",
		Done:        false,
	})

	testTable := []struct {
		name     string
		ctx      context.Context
		id       uint64
		dto      NoteDTO
		expErr   error
		expTitle string
		expDesc  string
		expDone  bool
	}{
		{
			name: "success put all",
			ctx:  context.Background(),
			id:   created.ID,
			dto: NoteDTO{
				Title:       "new title",
				Description: "new desc",
				Done:        true,
			},
			expErr:   nil,
			expTitle: "new title",
			expDesc:  "new desc",
			expDone:  true,
		},
		{
			name: "success put title",
			ctx:  context.Background(),
			id:   created.ID,
			dto: NoteDTO{
				Title: "only title",
			},
			expErr:   nil,
			expTitle: "only title",
			expDesc:  "",
			expDone:  false,
		},
		{
			name:   "not found",
			ctx:    context.Background(),
			id:     123,
			dto:    NoteDTO{Title: "x"},
			expErr: ErrNotFoundID,
		},
		{
			name:   "empty title",
			ctx:    context.Background(),
			id:     created.ID,
			dto:    NoteDTO{},
			expErr: ErrTitleNotDefined,
		},
		{
			name: "context deadline",
			ctx: func() context.Context {
				ctx, _ := context.WithTimeout(context.Background(), time.Nanosecond)
				time.Sleep(time.Millisecond)
				return ctx
			}(),
			id:     created.ID,
			dto:    NoteDTO{Title: "x"},
			expErr: context.DeadlineExceeded,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := repo.Update(testCase.ctx, testCase.id, testCase.dto)

			if !errors.Is(err, testCase.expErr) {
				t.Fatalf("expected error %v, got %v", testCase.expErr, err)
			}

			if testCase.expErr != nil {
				return
			}

			note, err := repo.GetByID(context.Background(), testCase.id)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if note.Title != testCase.expTitle {
				t.Errorf("title: expected %q, got %q", testCase.expTitle, note.Title)
			}
			if note.Description != testCase.expDesc {
				t.Errorf("description: expected %q, got %q", testCase.expDesc, note.Description)
			}
			if note.Done != testCase.expDone {
				t.Errorf("done: expected %v, got %v", testCase.expDone, note.Done)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	repo := NewInMemoryDataBase()
	created, _ := repo.Create(context.Background(), NoteDTO{Title: "x"})

	testTable := []struct {
		name   string
		ctx    context.Context
		id     uint64
		expErr error
	}{
		{
			name:   "success",
			ctx:    context.Background(),
			id:     created.ID,
			expErr: nil,
		},
		{
			name:   "not found",
			ctx:    context.Background(),
			id:     123,
			expErr: ErrNotFoundID,
		},
		{
			name: "context canceled",
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			id:     created.ID,
			expErr: context.Canceled,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			err := repo.Delete(testCase.ctx, testCase.id)

			if !errors.Is(err, testCase.expErr) {
				t.Fatalf("expected error %v, got %v", testCase.expErr, err)
			}
		})
	}
}
