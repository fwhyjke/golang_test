package router

import (
	"net/http"

	"github.com/fwhyjke/golang_test/internal/data"
	"github.com/fwhyjke/golang_test/internal/handler"
	"github.com/fwhyjke/golang_test/internal/middleware"
)

func NewToDoServerMux(db *data.DataBase) *http.ServeMux {
	mux := http.NewServeMux()
	h := handler.NewHandler(db)

	mux.Handle("/todos", middleware.Chain(h.HandleToDo(), middleware.LoggingMiddleware, middleware.TimeoutMiddleware))
	mux.Handle("/todos/", middleware.Chain(h.HandleToDoByID(), middleware.LoggingMiddleware, middleware.TimeoutMiddleware))

	return mux
}
