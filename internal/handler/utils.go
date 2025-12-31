package handler

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/fwhyjke/golang_test/internal/repository"
)

func handleError(w http.ResponseWriter, err error) {
	var statusCode int
	var message string
	var logMessage string

	switch {
	case errors.Is(err, context.Canceled):
		statusCode = http.StatusGatewayTimeout
		message = "time is out"
		logMessage = "timeout:" + err.Error()
	case errors.Is(err, context.DeadlineExceeded):
		statusCode = http.StatusGatewayTimeout
		message = "time is out"
		logMessage = "timeout:" + err.Error()

	case errors.Is(err, repository.ErrNotFoundID):
		statusCode = http.StatusNotFound
		message = err.Error()
		logMessage = err.Error()

	case errors.Is(err, repository.ErrTitleNotDefined):
		statusCode = http.StatusBadRequest
		message = "bad request: " + err.Error()
		logMessage = err.Error()

	default:
		statusCode = http.StatusInternalServerError
		message = "internal server error"
		logMessage = "internal error: " + err.Error()
	}

	log.Printf("error: code %d: %s", statusCode, logMessage)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(statusCode)
	w.Write([]byte(message + "\n"))
}
