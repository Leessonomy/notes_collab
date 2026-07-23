package usecase

import (
	"context"
	"notes-collab-api/internal/domain"
	"notes-collab-api/internal/dto"
	"time"

	"github.com/google/uuid"
)

type NoteStorage interface {
	Create(ctx context.Context, n domain.Note) error
	GetByWorkspaceID(ctx context.Context, workspaceID string) ([]domain.Note, error)
}

type Note struct {
	notes NoteStorage
}

func NewNote(notes NoteStorage) *Note {
	return &Note{notes: notes}
}

func (n *Note) Create(ctx context.Context, input dto.CreateNoteInput) (domain.Note, error) {
	now := time.Now().UTC()

	note := domain.Note{
		ID:          uuid.NewString(),
		WorkspaceID: input.WorkspaceID,
		Title:       input.Title,
		Content:     input.Content,
		AuthorID:    input.AuthorID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := n.notes.Create(ctx, note); err != nil {
		return domain.Note{}, err
	}

	return note, nil
}

func (n *Note) GetByWorkspaceID(ctx context.Context, workspaceID string) ([]domain.Note, error) {
	return n.notes.GetByWorkspaceID(ctx, workspaceID)
}
