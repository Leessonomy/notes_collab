package dto

import "time"

type CreateNoteInput struct {
	OwnerID     string `json:"-" validate:"required"`
	WorkspaceID string `json:"workspaceId" validate:"required"`
	Title       string `json:"title" mod:"trim" validate:"required,max=200"`
	Content     string `json:"content" validate:"max=100000"`
}

type UpdateNoteInput struct {
	NoteID  string `json:"-" validate:"required"`
	OwnerID string `json:"-" validate:"required"`
	Title   string `json:"title" mod:"trim" validate:"required,max=200"`
	Content string `json:"content" validate:"max=100000"`
}

type NoteOutput struct {
	ID          string    `json:"id"`
	WorkspaceID string    `json:"workspaceId"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	OwnerID     string    `json:"ownerId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
