package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/fwhyjke/golang_test/internal/data"
)

type noteDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

func PostNote(db *data.DataBase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var dto noteDTO
		if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
			http.Error(w, "err", http.StatusBadRequest)
			return
		}

		if strings.TrimSpace(dto.Title) == "" {
			http.Error(w, "err", http.StatusBadRequest)
			return
		}

		if strings.TrimSpace(dto.Description) == "" {
			http.Error(w, "err", http.StatusBadRequest)
			return
		}

		note := db.Create(dto.Title, dto.Description)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(note)
	}
}

func GetNotes(db *data.DataBase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(db.GetAll())
	}
}

func NoteByID(db *data.DataBase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/todos/")
		id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			http.Error(w, "err", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			note, ok := db.GetByID(id)
			if !ok {
				http.Error(w, "err", http.StatusNotFound)
				return
			}
			json.NewEncoder(w).Encode(note)

		case http.MethodPut:
			var dto noteDTO
			if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
				http.Error(w, "err", http.StatusBadRequest)
				return
			}

			if strings.TrimSpace(dto.Title) == "" {
				http.Error(w, "err", http.StatusBadRequest)
				return
			}

			note, ok := db.Update(id, dto.Title, dto.Description, dto.Done)
			if !ok {
				http.Error(w, "err", http.StatusNotFound)
				return
			}

			json.NewEncoder(w).Encode(note)

		case http.MethodDelete:
			if !db.Delete(id) {
				http.Error(w, "err", http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusNoContent)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}
