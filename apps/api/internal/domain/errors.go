package domain

import "errors"

var (
	ErrValidation = errors.New("validation error")

	ErrNoteByWorkspaceNotFound = errors.New("not found such note by workspace id")
	ErrNoteNotFound            = errors.New("note not found")
	ErrWorkspaceNotFound       = errors.New("workspace not found")

	ErrUserNotFound       = errors.New("user not found")
	ErrEmailTaken         = errors.New("email already taken")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrTokenNotFound      = errors.New("refresh token not found")
)
