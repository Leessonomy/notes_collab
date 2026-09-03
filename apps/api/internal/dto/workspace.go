package dto

import "notes-collab-api/internal/domain"

type CreateWorkspaceInput struct {
	OwnerID string `json:"-" validate:"required"`
	Name    string `json:"name" validate:"required,max=100"`
}

type UpdateWorkspaceInput struct {
	WorkspaceID string `json:"-" validate:"required"`
	OwnerID     string `json:"-" validate:"required"`
	Name        string `json:"name" validate:"required,max=100"`
}

type WorkspaceWithNotes struct {
	domain.Workspace
	Notes []domain.Note `json:"notes"`
}
