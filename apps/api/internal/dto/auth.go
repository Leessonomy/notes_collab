package dto

import (
	"notes-collab-api/internal/domain"
	"time"
)

type SignUpInput struct {
	Name     string
	Email    string
	Password string
}

type LogInInput struct {
	Email    string
	Password string
}

type SessionOutput struct {
	AccessToken      string
	RefreshToken     string
	AccessExpiresAt  time.Time
	RefreshExpiresAt time.Time
	User             *domain.User
}
