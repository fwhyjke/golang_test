package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/fwhyjke/golang_test/internal/repository"
)

func handleError(w http.ResponseWriter, err error) {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		http.Error(w, err.Error(), http.StatusGatewayTimeout)
	} else if errors.Is(err, repository.ErrNotFoundID) {
		http.Error(w, err.Error(), http.StatusNotFound)
	} else {
		http.Error(w, "failed: "+err.Error(), http.StatusInternalServerError)
	}
}
