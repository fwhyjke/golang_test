package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/fwhyjke/golang_test/internal/repository"
)

type Handler struct {
	repo repository.NoteRepository
}

func NewHandler(repo repository.NoteRepository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) postNote(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if !(r.Header.Get("Content-Type") == "application/json") {
		http.Error(w, "invalid media-type, must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	var dto repository.NoteDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(dto.Title) == "" {
		http.Error(w, "note must have a title", http.StatusBadRequest)
		return
	}

	note, err := h.repo.Create(ctx, dto)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(note)
}

func (h *Handler) getNotes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	notes, err := h.repo.GetAll(ctx)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notes)
}

func (h *Handler) getNoteByID(w http.ResponseWriter, r *http.Request, id uint64) {
	ctx := r.Context()

	note, err := h.repo.GetByID(ctx, id)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(note)
}

func (h *Handler) putNoteByID(w http.ResponseWriter, r *http.Request, id uint64) {
	ctx := r.Context()
	if !(r.Header.Get("Content-Type") == "application/json") {
		http.Error(w, "invalid media-type, must be application/json", http.StatusUnsupportedMediaType)
	}

	var dto repository.NoteDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(dto.Title) == "" {
		http.Error(w, "note must have a title", http.StatusBadRequest)
		return
	}

	note, err := h.repo.Update(ctx, id, dto)
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(note)
}

func (h *Handler) deleteNoteByID(w http.ResponseWriter, r *http.Request, id uint64) {
	ctx := r.Context()

	if err := h.repo.Delete(ctx, id); err != nil {
		handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
