package controller

import (
	"errors"
	"net/http"
	"notes-collab-api/internal/domain"
	"notes-collab-api/internal/utils"
)

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
		utils.InternalError(w, err)
	}
}
