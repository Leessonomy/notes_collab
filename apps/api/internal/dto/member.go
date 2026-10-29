package dto

import (
	"notes-collab-api/internal/domain"
	"time"
)

type ChangeMemberRoleInput struct {
	WorkspaceID string      `json:"-" validate:"required"`
	UserID      string      `json:"-" validate:"required"`
	MemberID    string      `json:"-" validate:"required"`
	Role        domain.Role `json:"role" validate:"required,oneof=collaborator guest"`
}

type MemberOutput struct {
	UserID   string      `json:"userId"`
	Name     string      `json:"name"`
	Email    string      `json:"email"`
	Role     domain.Role `json:"role"`
	JoinedAt time.Time   `json:"joinedAt"`
}
