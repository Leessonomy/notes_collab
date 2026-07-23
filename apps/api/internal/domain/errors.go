package domain

import "errors"

var (
	ErrNoteByWorkspaceNotFound = errors.New("not found such note by workspace id")

	ErrUserNotFound       = errors.New("user not found")
	ErrEmailTaken         = errors.New("email already taken")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrTokenNotFound      = errors.New("refresh token not found")
)
