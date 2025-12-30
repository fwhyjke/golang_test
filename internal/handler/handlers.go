package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/fwhyjke/golang_test/internal/data"
)

type Handler struct {
	db *data.DataBase
}

func NewHandler(database *data.DataBase) *Handler {
	return &Handler{
		db: database,
	}
}

func (h *Handler) postNote(w http.ResponseWriter, r *http.Request) {
	var dto data.NoteDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "Invalid json", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(dto.Title) == "" {
		http.Error(w, "Note must have a title", http.StatusBadRequest)
		return
	}

	note := h.db.Create(dto)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(note)
}

func (h *Handler) getNotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.db.GetAll())
}

func (h *Handler) getNoteByID(w http.ResponseWriter, r *http.Request, id uint64) {
	note, ok := h.db.GetByID(id)
	if !ok {
		http.Error(w, "Invalid id", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(note)
}

func (h *Handler) putNoteByID(w http.ResponseWriter, r *http.Request, id uint64) {
	var dto data.NoteDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "Invalid json", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(dto.Title) == "" {
		http.Error(w, "Note must have a title", http.StatusBadRequest)
		return
	}

	note, ok := h.db.Update(id, dto)
	if !ok {
		http.Error(w, "Invalid id", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(note)
}

func (h *Handler) deleteNoteByID(w http.ResponseWriter, r *http.Request, id uint64) {
	if !h.db.Delete(id) {
		http.Error(w, "Invalid id", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
