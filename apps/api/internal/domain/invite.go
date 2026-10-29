package domain

import (
	"crypto/rand"
	"crypto/sha256"
	"time"
)

const InviteTTL = 7 * 24 * time.Hour

type Invite struct {
	ID          string
	WorkspaceID string
	Role        Role
	InviterID   string
	ExpiresAt   time.Time
	CreatedAt   time.Time
}

type InviteWithWorkspace struct {
	Invite
	WorkspaceName string
	InviterName   string
}

type InviteToken string

func NewInviteToken() InviteToken {
	return InviteToken(rand.Text())
}

func (t InviteToken) Hash() []byte {
	sum := sha256.Sum256([]byte(t))
	return sum[:]
}
