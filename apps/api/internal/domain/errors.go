package domain

import "errors"

var (
	ErrValidation = errors.New("validation error")

	ErrNoteByWorkspaceNotFound = errors.New("not found such note by workspace id")
	ErrNoteNotFound            = errors.New("note not found")
	ErrWorkspaceNotFound       = errors.New("workspace not found")

	ErrNotMember     = errors.New("not a workspace member")
	ErrAlreadyMember = errors.New("already a workspace member")
	ErrForbidden     = errors.New("forbidden")

	ErrInviteNotFound = errors.New("invite not found")
	ErrInviteInvalid  = errors.New("invite is invalid or expired")

	ErrUserNotFound       = errors.New("user not found")
	ErrEmailTaken         = errors.New("email already taken")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrWrongPassword      = errors.New("current password is incorrect")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrTokenNotFound      = errors.New("refresh token not found")
)
