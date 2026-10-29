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
	ListByMember(ctx context.Context, userID string) ([]domain.Note, error)
	Delete(ctx context.Context, noteID string) error
	GetByID(ctx context.Context, noteID string) (domain.Note, error)
	Update(ctx context.Context, noteID string, title, content *string) (domain.Note, error)
}

type Note struct {
	noteRepo       NoteRepo
	membershipRepo MembershipRepo
}

func NewNote(noteRepo NoteRepo, membershipRepo MembershipRepo) *Note {
	return &Note{noteRepo: noteRepo, membershipRepo: membershipRepo}
}

func (n *Note) GetByID(ctx context.Context, noteID, userID string) (domain.Note, error) {
	note, err := n.noteRepo.GetByID(ctx, noteID)
	if err != nil {
		return domain.Note{}, err
	}

	membership, err := getMembership(ctx, n.membershipRepo, note.WorkspaceID, userID)
	if err != nil {
		return domain.Note{}, err
	}
	if !membership.Role.Rights().Read {
		return domain.Note{}, domain.ErrForbidden
	}

	return note, nil
}

func (n *Note) Create(ctx context.Context, input dto.CreateNoteInput) (domain.Note, error) {
	if err := validateInput(&input); err != nil {
		return domain.Note{}, err
	}

	membership, err := getMembership(ctx, n.membershipRepo, input.WorkspaceID, input.OwnerID)
	if err != nil {
		return domain.Note{}, err
	}
	if !membership.Role.Rights().Write {
		return domain.Note{}, domain.ErrForbidden
	}

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

func (n *Note) Update(ctx context.Context, input dto.UpdateNoteInput) (domain.Note, error) {
	if err := validateInput(&input); err != nil {
		return domain.Note{}, err
	}

	note, err := n.noteRepo.GetByID(ctx, input.NoteID)
	if err != nil {
		return domain.Note{}, err
	}

	membership, err := getMembership(ctx, n.membershipRepo, note.WorkspaceID, input.UserID)
	if err != nil {
		return domain.Note{}, err
	}
	if !membership.Role.Rights().Write {
		return domain.Note{}, domain.ErrForbidden
	}

	return n.noteRepo.Update(ctx, input.NoteID, input.Title, input.Content)
}

func (n *Note) Delete(ctx context.Context, noteID, userID string) error {
	note, err := n.noteRepo.GetByID(ctx, noteID)
	if err != nil {
		return err
	}

	membership, err := getMembership(ctx, n.membershipRepo, note.WorkspaceID, userID)
	if err != nil {
		return err
	}
	if !membership.Role.Rights().Write {
		return domain.ErrForbidden
	}

	return n.noteRepo.Delete(ctx, noteID)
}
