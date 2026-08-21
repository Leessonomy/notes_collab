package controller

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"notes-collab-api/internal/domain"
)

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func internalError(w http.ResponseWriter, err error) {
	log.Printf("internal error: %v", err)
	http.Error(w, "internal error", http.StatusInternalServerError)
}

func unauthorized(w http.ResponseWriter) {
	http.Error(w, "unauthorized", http.StatusUnauthorized)
}

func badRequest(w http.ResponseWriter, message string) {
	http.Error(w, message, http.StatusBadRequest)
}

func domainError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrEmailTaken):
		http.Error(w, err.Error(), http.StatusConflict)
	case errors.Is(err, domain.ErrNoteNotFound),
		errors.Is(err, domain.ErrNoteByWorkspaceNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, domain.ErrInvalidCredentials),
		errors.Is(err, domain.ErrTokenNotFound),
		errors.Is(err, domain.ErrUnauthorized):
		http.Error(w, err.Error(), http.StatusUnauthorized)
	default:
		internalError(w, err)
	}
}
