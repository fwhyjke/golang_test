package handler

import (
	"net/http"
	"strconv"
	"strings"
)

func (h *Handler) HandleToDo() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodPost:
				h.postNote(w, r)
			case http.MethodGet:
				h.getNotes(w, r)
			default:
				w.Header().Set("Allow", "GET, POST")
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		},
	)
}

func (h *Handler) HandleToDoByID() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			idStr := strings.TrimPrefix(r.URL.Path, "/todos/")
			id, err := strconv.ParseUint(idStr, 10, 64)
			if err != nil {
				http.Error(w, "Invalid id in url", http.StatusBadRequest)
				return
			}

			switch r.Method {
			case http.MethodGet:
				h.getNoteByID(w, r, id)
			case http.MethodPut:
				h.putNoteByID(w, r, id)
			case http.MethodDelete:
				h.deleteNoteByID(w, r, id)
			default:
				w.Header().Set("Allow", "GET, PUT, DELETE")
				w.WriteHeader(http.StatusMethodNotAllowed)
			}
		},
	)
}
