package dto

type CreateNoteInput struct {
	OwnerID     string `json:"-" validate:"required"`
	WorkspaceID string `json:"workspaceId" validate:"required"`
	Title       string `json:"title" validate:"required,max=200"`
	Content     string `json:"content" validate:"max=100000"`
}

type UpdateNoteInput struct {
	NoteID  string `json:"-" validate:"required"`
	OwnerID string `json:"-" validate:"required"`
	Title   string `json:"title" validate:"required,max=200"`
	Content string `json:"content" validate:"max=100000"`
}
