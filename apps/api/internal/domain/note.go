package domain

import "time"

type Note struct {
	ID          string
	WorkspaceID string
	Title       string
	Content     string
	OwnerID     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
