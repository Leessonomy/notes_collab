package controller

import (
	"encoding/json"
	"net/http"
	"notes-collab-api/internal/domain"
	"notes-collab-api/internal/repository"
	"notes-collab-api/internal/utils"
	"strings"
	"time"

	"github.com/google/uuid"
)

type NoteController struct {
	repo *repository.NoteRepo
}

func NewNoteController(repo *repository.NoteRepo) *NoteController {
	return &NoteController{repo: repo}
}

func (c *NoteController) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var b struct {
		Title       string `json:"title"`
		Content     string `json:"content"`
		WorkspaceID string `json:"workspaceId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(b.WorkspaceID) == "" {
		http.Error(w, "workspaceId is required", http.StatusBadRequest)
		return
	}

	now := time.Now().UTC()

	note := domain.Note{
		ID:          uuid.NewString(),
		WorkspaceID: b.WorkspaceID,
		Title:       b.Title,
		Content:     b.Content,
		AuthorID:    userID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := c.repo.Create(r.Context(), note); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, note)
}
