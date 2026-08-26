package dto

type CreateNoteInput struct {
	WorkspaceID string
	Title       string
	Content     string
	OwnerID     string
}

type UpdateNoteInput struct {
	NoteID  string
	OwnerID string
	Title   string
	Content string
}
