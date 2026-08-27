package dto

import "notes-collab-api/internal/domain"

type CreateWorkspaceInput struct {
	Name    string
	OwnerID string
}

type UpdateWorkspaceInput struct {
	WorkspaceID string
	OwnerID     string
	Name        string
}

type WorkspaceWithNotes struct {
	domain.Workspace
	Notes []domain.Note `json:"notes"`
}
