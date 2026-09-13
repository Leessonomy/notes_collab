package domain

import "time"

type WorkspaceMember struct {
	WorkspaceID string
	UserID      string
	Role        Role
	JoinedAt    time.Time
}

type WorkspaceWithMember struct {
	Workspace Workspace
	Member    WorkspaceMember
}
