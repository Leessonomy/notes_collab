package dto

import "notes-collab-api/internal/domain"

type CreateWorkspaceInput struct {
	Name    string
	OwnerID string
}

type WorkspaceWithNotes struct {
	domain.Workspace
	Notes []domain.Note `json:"notes"`
}
