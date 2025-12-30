package router

import (
	"net/http"

	"github.com/fwhyjke/golang_test/internal/handler"
	"github.com/fwhyjke/golang_test/internal/middleware"
	"github.com/fwhyjke/golang_test/internal/repository"
)

func NewToDoServerMux(db *repository.DataBase) *http.ServeMux {
	mux := http.NewServeMux()
	h := handler.NewHandler(db)

	mux.Handle("/todos", middleware.Chain(h.HandleToDo(), middleware.LoggingMiddleware, middleware.TimeoutMiddleware))
	mux.Handle("/todos/", middleware.Chain(h.HandleToDoByID(), middleware.LoggingMiddleware, middleware.TimeoutMiddleware))

	return mux
}
