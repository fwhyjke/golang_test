package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/fwhyjke/golang_test/internal/repository"
)

type MockRepository struct {
	CreateFunc  func(ctx context.Context, dto repository.NoteDTO) (repository.Note, error)
	GetByIDFunc func(ctx context.Context, id uint64) (repository.Note, error)
	GetAllFunc  func(ctx context.Context) ([]repository.Note, error)
	UpdateFunc  func(ctx context.Context, id uint64, dto repository.NoteDTO) (repository.Note, error)
	DeleteFunc  func(ctx context.Context, id uint64) error
}

func (m *MockRepository) Create(ctx context.Context, dto repository.NoteDTO) (repository.Note, error) {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, dto)
	}
	return repository.Note{}, nil
}

func (m *MockRepository) GetByID(ctx context.Context, id uint64) (repository.Note, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return repository.Note{}, nil
}

func (m *MockRepository) GetAll(ctx context.Context) ([]repository.Note, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx)
	}
	return []repository.Note{}, nil
}

func (m *MockRepository) Update(ctx context.Context, id uint64, dto repository.NoteDTO) (repository.Note, error) {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, id, dto)
	}
	return repository.Note{}, nil
}

func (m *MockRepository) Delete(ctx context.Context, id uint64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func TestPostNote(t *testing.T) {
	testTable := []struct {
		name        string
		req         string
		contentType string
		mockCreate  func(ctx context.Context, dto repository.NoteDTO) (repository.Note, error)
		expStatus   int
		expBody     string
	}{
		{
			name:        "success",
			req:         `{"title": "123", "description": "qwe", "done": true}`,
			contentType: "application/json",
			mockCreate: func(ctx context.Context, dto repository.NoteDTO) (repository.Note, error) {
				return repository.Note{
					ID:          1,
					Title:       dto.Title,
					Description: dto.Description,
					Done:        dto.Done,
				}, nil
			},
			expStatus: http.StatusCreated,
			expBody:   `{"id":1,"title":"123","description":"qwe","done":true}`,
		},
		{
			name:        "empty title",
			req:         `{"title": ""}`,
			contentType: "application/json",
			expStatus:   http.StatusBadRequest,
			expBody:     "note must have a title",
		},
		{
			name:        "invalid json",
			req:         `{asdad}`,
			contentType: "application/json",
			expStatus:   http.StatusBadRequest,
			expBody:     "invalid json",
		},
		{
			name:        "wrong content type",
			req:         `{"title": "123"}`,
			contentType: "text/plain",
			expStatus:   http.StatusUnsupportedMediaType,
			expBody:     "invalid media-type, must be application/json",
		},
		{
			name:        "repository error",
			req:         `{"title": "123"}`,
			contentType: "application/json",
			mockCreate: func(ctx context.Context, dto repository.NoteDTO) (repository.Note, error) {
				return repository.Note{}, errors.New("error")
			},
			expStatus: http.StatusInternalServerError,
			expBody:   "internal server error",
		},
		{
			name:        "repository validation error",
			req:         `{"title": "123"}`,
			contentType: "application/json",
			mockCreate: func(ctx context.Context, dto repository.NoteDTO) (repository.Note, error) {
				return repository.Note{}, repository.ErrTitleNotDefined
			},
			expStatus: http.StatusBadRequest,
			expBody:   "bad request: title is required field",
		},
		{
			name:        "context timeout",
			req:         `{"title": "Test"}`,
			contentType: "application/json",
			mockCreate: func(ctx context.Context, dto repository.NoteDTO) (repository.Note, error) {
				return repository.Note{}, context.DeadlineExceeded
			},
			expStatus: http.StatusGatewayTimeout,
			expBody:   "time is out",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			mockRepo := &MockRepository{CreateFunc: testCase.mockCreate}

			handler := NewHandler(mockRepo)

			req := httptest.NewRequest("POST", "/todos", strings.NewReader(testCase.req))
			req.Header.Set("Content-Type", testCase.contentType)

			rec := httptest.NewRecorder()

			handler.postNote(rec, req)

			if status := rec.Code; status != testCase.expStatus {
				t.Errorf("status code: expected %v, got %v", testCase.expStatus, status)
			}

			body := strings.TrimSpace(rec.Body.String())
			if body != testCase.expBody {
				t.Errorf("body: expected %v, got %v", testCase.expBody, body)
			}
		})
	}
}

func TestGetNotes(t *testing.T) {
	testTable := []struct {
		name       string
		mockGetAll func(ctx context.Context) ([]repository.Note, error)
		expStatus  int
		expBody    string
	}{
		{
			name: "success couple",
			mockGetAll: func(ctx context.Context) ([]repository.Note, error) {
				return []repository.Note{
					{ID: 1, Title: "t1", Description: "d1", Done: false},
					{ID: 2, Title: "n2", Description: "d2", Done: true},
				}, nil
			},
			expStatus: http.StatusOK,
			expBody:   `[{"id":1,"title":"t1","description":"d1","done":false},{"id":2,"title":"n2","description":"d2","done":true}]`,
		},
		{
			name: "success empty",
			mockGetAll: func(ctx context.Context) ([]repository.Note, error) {
				return []repository.Note{}, nil
			},
			expStatus: http.StatusOK,
			expBody:   "[]",
		},
		{
			name: "context cancelled",
			mockGetAll: func(ctx context.Context) ([]repository.Note, error) {
				return nil, context.Canceled
			},
			expStatus: http.StatusGatewayTimeout,
			expBody:   "time is out",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			mockRepo := &MockRepository{GetAllFunc: testCase.mockGetAll}
			handler := NewHandler(mockRepo)

			req := httptest.NewRequest("GET", "/todos", nil)
			rec := httptest.NewRecorder()

			handler.getNotes(rec, req)

			if status := rec.Code; status != testCase.expStatus {
				t.Errorf("status code: expected %v, got %v", testCase.expStatus, status)
			}

			body := strings.TrimSpace(rec.Body.String())
			expectedBody := strings.TrimSpace(testCase.expBody)
			if body != expectedBody {
				t.Errorf("body: expected %v, got %v", expectedBody, body)
			}

		})
	}
}

func TestGetNoteByID(t *testing.T) {
	testTable := []struct {
		name        string
		id          uint64
		mockGetByID func(ctx context.Context, id uint64) (repository.Note, error)
		expStatus   int
		expBody     string
	}{
		{
			name: "success",
			id:   1,
			mockGetByID: func(ctx context.Context, id uint64) (repository.Note, error) {
				return repository.Note{
					ID:          id,
					Title:       "qwe",
					Description: "qwe",
					Done:        false,
				}, nil
			},
			expStatus: http.StatusOK,
			expBody:   `{"id":1,"title":"qwe","description":"qwe","done":false}`,
		},
		{
			name: "invalid id",
			id:   12345,
			mockGetByID: func(ctx context.Context, id uint64) (repository.Note, error) {
				return repository.Note{}, repository.ErrNotFoundID
			},
			expStatus: http.StatusNotFound,
			expBody:   "note by ID not found",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			mockRepo := &MockRepository{GetByIDFunc: testCase.mockGetByID}
			handler := NewHandler(mockRepo)

			req := httptest.NewRequest("GET", "/todos/1", nil)
			rec := httptest.NewRecorder()

			handler.getNoteByID(rec, req, testCase.id)

			if status := rec.Code; status != testCase.expStatus {
				t.Errorf("status code: expected %v, got %v", testCase.expStatus, status)
			}

			body := strings.TrimSpace(rec.Body.String())
			expectedBody := strings.TrimSpace(testCase.expBody)
			if body != expectedBody {
				t.Errorf("body: expected %v, got %v", expectedBody, body)
			}
		})
	}
}

func TestPutNoteByID(t *testing.T) {
	testTable := []struct {
		name        string
		id          uint64
		req         string
		contentType string
		mockUpdate  func(ctx context.Context, id uint64, dto repository.NoteDTO) (repository.Note, error)
		expStatus   int
		expBody     string
	}{
		{
			name:        "success",
			id:          1,
			req:         `{"title": "t", "description": "d", "done": true}`,
			contentType: "application/json",
			mockUpdate: func(ctx context.Context, id uint64, dto repository.NoteDTO) (repository.Note, error) {
				return repository.Note{
					ID:          id,
					Title:       dto.Title,
					Description: dto.Description,
					Done:        dto.Done,
				}, nil
			},
			expStatus: http.StatusOK,
			expBody:   `{"id":1,"title":"t","description":"d","done":true}`,
		},
		{
			name:        "empty title",
			id:          1,
			req:         `{"title": ""}`,
			contentType: "application/json",
			expStatus:   http.StatusBadRequest,
			expBody:     "note must have a title",
		},
		{
			name:        "invalid id",
			id:          527892,
			req:         `{"title": "qwe"}`,
			contentType: "application/json",
			mockUpdate: func(ctx context.Context, id uint64, dto repository.NoteDTO) (repository.Note, error) {
				return repository.Note{}, repository.ErrNotFoundID
			},
			expStatus: http.StatusNotFound,
			expBody:   "note by ID not found",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			mockRepo := &MockRepository{UpdateFunc: testCase.mockUpdate}
			handler := NewHandler(mockRepo)

			req := httptest.NewRequest("PUT", "/todos/1", strings.NewReader(testCase.req))
			req.Header.Set("Content-Type", testCase.contentType)
			rec := httptest.NewRecorder()

			handler.putNoteByID(rec, req, testCase.id)

			if status := rec.Code; status != testCase.expStatus {
				t.Errorf("status code: expected %v, got %v", testCase.expStatus, status)
			}

			body := strings.TrimSpace(rec.Body.String())
			expectedBody := strings.TrimSpace(testCase.expBody)
			if body != expectedBody {
				t.Errorf("body: expected %v, got %v", expectedBody, body)
			}
		})
	}
}

func TestDeleteNoteByID(t *testing.T) {
	testTable := []struct {
		name       string
		id         uint64
		mockDelete func(ctx context.Context, id uint64) error
		expStatus  int
	}{
		{
			name: "success",
			id:   1,
			mockDelete: func(ctx context.Context, id uint64) error {
				return nil
			},
			expStatus: http.StatusNoContent,
		},
		{
			name: "invalid id",
			id:   12345,
			mockDelete: func(ctx context.Context, id uint64) error {
				return repository.ErrNotFoundID
			},
			expStatus: http.StatusNotFound,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			mockRepo := &MockRepository{DeleteFunc: testCase.mockDelete}
			handler := NewHandler(mockRepo)

			req := httptest.NewRequest("DELETE", "/todos/1", nil)
			rec := httptest.NewRecorder()

			handler.deleteNoteByID(rec, req, testCase.id)

			if status := rec.Code; status != testCase.expStatus {
				t.Errorf("status code: expected %v, got %v", testCase.expStatus, status)
			}
		})
	}
}
