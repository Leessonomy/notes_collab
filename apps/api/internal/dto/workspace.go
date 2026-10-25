package dto

import (
	"notes-collab-api/internal/domain"
	"time"
)

type CreateWorkspaceInput struct {
	OwnerID string `json:"-" validate:"required"`
	Name    string `json:"name" mod:"trim" validate:"required,max=100"`
}

type UpdateWorkspaceInput struct {
	WorkspaceID string `json:"-" validate:"required"`
	UserID      string `json:"-" validate:"required"`
	Name        string `json:"name" mod:"trim" validate:"required,max=100"`
}

type MembershipOutput struct {
	Role     domain.Role `json:"role"`
	JoinedAt time.Time   `json:"joinedAt"`
}

type WorkspaceOutput struct {
	ID         string           `json:"id"`
	Name       string           `json:"name"`
	OwnerID    string           `json:"ownerId"`
	Membership MembershipOutput `json:"membership"`
	CreatedAt  time.Time        `json:"createdAt"`
	UpdatedAt  time.Time        `json:"updatedAt"`
}

type WorkspaceWithNotes struct {
	WorkspaceOutput
	Notes []NoteOutput `json:"notes"`
}
