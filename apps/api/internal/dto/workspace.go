package dto

import "time"

type CreateWorkspaceInput struct {
	OwnerID string `json:"-" validate:"required"`
	Name    string `json:"name" mod:"trim" validate:"required,max=100"`
}

type UpdateWorkspaceInput struct {
	WorkspaceID string `json:"-" validate:"required"`
	OwnerID     string `json:"-" validate:"required"`
	Name        string `json:"name" mod:"trim" validate:"required,max=100"`
}

type WorkspaceOutput struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	OwnerID   string    `json:"ownerId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type WorkspaceWithNotes struct {
	WorkspaceOutput
	Notes []NoteOutput `json:"notes"`
}
