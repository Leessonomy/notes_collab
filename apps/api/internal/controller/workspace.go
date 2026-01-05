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

type WorkspaceController struct {
	repo *repository.WorkspaceRepo
}

func NewWorkspaceController(repo *repository.WorkspaceRepo) *WorkspaceController {
	return &WorkspaceController{repo: repo}
}

func (c *WorkspaceController) GetAll(w http.ResponseWriter, r *http.Request) {
	data, err := c.repo.Get(r.Context())
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, data)
}

func (c *WorkspaceController) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var b struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	b.Name = strings.TrimSpace(b.Name)
	if b.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	now := time.Now().UTC()

	workspace := domain.Workspace{
		ID:        uuid.NewString(),
		Name:      b.Name,
		OwnerID:   userID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := c.repo.Create(r.Context(), workspace); err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, workspace)
}

func (c *NoteController) GetByWorkspaceID(w http.ResponseWriter, r *http.Request) {
	notes, err := c.repo.GetByWorkspaceID(r.Context(), r.PathValue("workspaceId"))
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, notes)
}
