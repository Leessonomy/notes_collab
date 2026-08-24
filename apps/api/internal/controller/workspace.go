package controller

import (
	"encoding/json"
	"net/http"
	"notes-collab-api/internal/dto"
	"notes-collab-api/internal/usecase"
	"notes-collab-api/internal/utils"
)

type WorkspaceController struct {
	workspaceUseCase *usecase.Workspace
}

func NewWorkspaceController(workspace *usecase.Workspace) *WorkspaceController {
	return &WorkspaceController{workspaceUseCase: workspace}
}

func (c *WorkspaceController) GetWithNotes(w http.ResponseWriter, r *http.Request) {
	userID := utils.UserIDFromContext(r.Context())

	data, err := c.workspaceUseCase.GetWithNotes(r.Context(), userID)

	if err != nil {
		domainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, data)
}

func (c *WorkspaceController) Create(w http.ResponseWriter, r *http.Request) {
	userID := utils.UserIDFromContext(r.Context())

	var body struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		badRequest(w, "invalid body")
		return
	}

	workspace, err := c.workspaceUseCase.Create(r.Context(), dto.CreateWorkspaceInput{
		Name:    body.Name,
		OwnerID: userID,
	})
	if err != nil {
		domainError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, workspace)
}

func (c *WorkspaceController) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := c.workspaceUseCase.Delete(r.Context(), id)

	if err != nil {
		domainError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
