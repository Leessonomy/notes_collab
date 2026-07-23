package controller

import (
	"encoding/json"
	"net/http"
	"notes-collab-api/internal/dto"
	"notes-collab-api/internal/usecase"
	"notes-collab-api/internal/utils"
)

type WorkspaceController struct {
	workspace *usecase.Workspace
}

func NewWorkspaceController(workspace *usecase.Workspace) *WorkspaceController {
	return &WorkspaceController{workspace: workspace}
}

func (c *WorkspaceController) GetAll(w http.ResponseWriter, r *http.Request) {
	data, err := c.workspace.GetAll(r.Context())
	if err != nil {
		domainError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, data)
}

func (c *WorkspaceController) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.UserIDFromContext(r.Context())
	if !ok {
		unauthorized(w)
		return
	}

	var body struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		badRequest(w, "invalid body")
		return
	}

	workspace, err := c.workspace.Create(r.Context(), dto.CreateWorkspaceInput{
		Name:    body.Name,
		OwnerID: userID,
	})
	if err != nil {
		domainError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, workspace)
}
