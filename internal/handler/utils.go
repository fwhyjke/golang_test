package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/fwhyjke/golang_test/internal/data"
)

func handleError(w http.ResponseWriter, err error) {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		http.Error(w, err.Error(), http.StatusGatewayTimeout)
	} else if errors.Is(err, data.ErrNotFoundID) {
		http.Error(w, err.Error(), http.StatusNotFound)
	} else {
		http.Error(w, "failed: "+err.Error(), http.StatusInternalServerError)
	}
}
