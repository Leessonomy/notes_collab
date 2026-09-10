package dto

import (
	"notes-collab-api/internal/domain"
	"time"
)

type SignUpInput struct {
	Name     string `json:"name" mod:"trim" validate:"required,min=2"`
	Email    string `json:"email" mod:"trim" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LogInInput struct {
	Email    string `json:"email" mod:"trim" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type ChangePasswordInput struct {
	UserID          string `json:"-" validate:"required"`
	CurrentPassword string `json:"currentPassword" validate:"required"`
	NewPassword     string `json:"newPassword" validate:"required,min=8"`
}

type SessionOutput struct {
	AccessToken      string
	RefreshToken     string
	AccessExpiresAt  time.Time
	RefreshExpiresAt time.Time
	User             *domain.User
}

type UserOutput struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
