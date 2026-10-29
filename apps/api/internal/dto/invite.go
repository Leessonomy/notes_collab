package dto

import (
	"notes-collab-api/internal/domain"
	"time"
)

type CreateInviteInput struct {
	WorkspaceID string      `json:"-" validate:"required"`
	UserID      string      `json:"-" validate:"required"`
	Role        domain.Role `json:"role" validate:"required,oneof=collaborator guest"`
}

type InviteOutput struct {
	ID        string      `json:"id"`
	Role      domain.Role `json:"role"`
	ExpiresAt time.Time   `json:"expiresAt"`
	CreatedAt time.Time   `json:"createdAt"`
}

type CreatedInviteOutput struct {
	InviteOutput
	Token domain.InviteToken `json:"token"`
}

type InviteWithWorkspaceOutput struct {
	InviteOutput
	WorkspaceName string `json:"workspaceName"`
	InviterName   string `json:"inviterName"`
}

type AcceptInviteOutput struct {
	WorkspaceID string      `json:"workspaceId"`
	Role        domain.Role `json:"role"`
	JoinedAt    time.Time   `json:"joinedAt"`
}
