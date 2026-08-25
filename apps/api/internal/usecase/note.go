package usecase

import (
	"context"
	"notes-collab-api/internal/domain"
	"notes-collab-api/internal/dto"
	"time"

	"github.com/google/uuid"
)

type NoteRepo interface {
	Create(ctx context.Context, n domain.Note) error
	GetByOwner(ctx context.Context, ownerID string) ([]domain.Note, error)
	Delete(ctx context.Context, noteID, ownerID string) error
}

type Note struct {
	noteRepo NoteRepo
}

func NewNote(noteRepo NoteRepo) *Note {
	return &Note{noteRepo: noteRepo}
}

func (n *Note) Create(ctx context.Context, input dto.CreateNoteInput) (domain.Note, error) {
	now := time.Now().UTC()

	note := domain.Note{
		ID:          uuid.NewString(),
		WorkspaceID: input.WorkspaceID,
		Title:       input.Title,
		Content:     input.Content,
		OwnerID:     input.OwnerID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := n.noteRepo.Create(ctx, note); err != nil {
		return domain.Note{}, err
	}

	return note, nil
}

func (n *Note) Delete(ctx context.Context, noteID, ownerID string) error {
	return n.noteRepo.Delete(ctx, noteID, ownerID)
}
