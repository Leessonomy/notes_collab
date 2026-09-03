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

func NewWorkspaceController(workspaceUC *usecase.Workspace) *WorkspaceController {
	return &WorkspaceController{workspaceUseCase: workspaceUC}
}

func (c *WorkspaceController) ListWorkspacesWithNotes(w http.ResponseWriter, r *http.Request) {
	userID := utils.UserIDFromContext(r.Context())

	data, err := c.workspaceUseCase.ListWithNotes(r.Context(), userID)

	if err != nil {
		domainError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, data)
}

func (c *WorkspaceController) CreateWorkspace(w http.ResponseWriter, r *http.Request) {
	input := dto.CreateWorkspaceInput{
		OwnerID: utils.UserIDFromContext(r.Context()),
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.BadRequest(w, "invalid body")
		return
	}

	workspace, err := c.workspaceUseCase.Create(r.Context(), input)
	if err != nil {
		domainError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, workspace)
}

func (c *WorkspaceController) UpdateWorkspace(w http.ResponseWriter, r *http.Request) {
	input := dto.UpdateWorkspaceInput{
		WorkspaceID: r.PathValue("id"),
		OwnerID:     utils.UserIDFromContext(r.Context()),
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.BadRequest(w, "invalid body")
		return
	}

	workspace, err := c.workspaceUseCase.Update(r.Context(), input)
	if err != nil {
		domainError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, workspace)
}

func (c *WorkspaceController) DeleteWorkspace(w http.ResponseWriter, r *http.Request) {
	userID := utils.UserIDFromContext(r.Context())
	id := r.PathValue("id")

	err := c.workspaceUseCase.Delete(r.Context(), id, userID)

	if err != nil {
		domainError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
