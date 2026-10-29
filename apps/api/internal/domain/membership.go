package domain

import "time"

type Membership struct {
	WorkspaceID string
	UserID      string
	Role        Role
	JoinedAt    time.Time
}

type WorkspaceWithMembership struct {
	Workspace
	Membership Membership
}

type Member struct {
	UserID   string
	Name     string
	Email    string
	Role     Role
	JoinedAt time.Time
}
