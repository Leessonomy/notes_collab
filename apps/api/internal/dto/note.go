package dto

type CreateNoteInput struct {
	WorkspaceID string
	Title       string
	Content     string
	OwnerID     string
}
