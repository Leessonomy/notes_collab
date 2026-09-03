package usecase

import (
	"context"
	"notes-collab-api/internal/domain"
	"notes-collab-api/internal/dto"
	"time"

	"github.com/google/uuid"
)

type WorkspaceRepo interface {
	Create(ctx context.Context, w domain.Workspace) error
	ListByOwner(ctx context.Context, ownerID string) ([]domain.Workspace, error)
	Update(ctx context.Context, workspaceID, ownerID, name string) (domain.Workspace, error)
	Delete(ctx context.Context, workspaceID, ownerID string) error
}

type Workspace struct {
	workspaceRepo WorkspaceRepo
	noteRepo      NoteRepo
}

func NewWorkspace(workspaceRepo WorkspaceRepo, noteRepo NoteRepo) *Workspace {
	return &Workspace{workspaceRepo: workspaceRepo, noteRepo: noteRepo}
}

func (w *Workspace) Create(ctx context.Context, input dto.CreateWorkspaceInput) (domain.Workspace, error) {
	if err := validateInput(input); err != nil {
		return domain.Workspace{}, err
	}

	now := time.Now().UTC()

	workspace := domain.Workspace{
		ID:        uuid.NewString(),
		Name:      input.Name,
		OwnerID:   input.OwnerID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := w.workspaceRepo.Create(ctx, workspace); err != nil {
		return domain.Workspace{}, err
	}

	return workspace, nil
}

func (w *Workspace) ListWithNotes(ctx context.Context, ownerID string) ([]dto.WorkspaceWithNotes, error) {
	workspaces, err := w.workspaceRepo.ListByOwner(ctx, ownerID)

	if err != nil {
		return nil, err
	}

	notes, err := w.noteRepo.ListByOwner(ctx, ownerID)

	if err != nil {
		return nil, err
	}

	result := make([]dto.WorkspaceWithNotes, 0, len(workspaces))

	for _, workspace := range workspaces {
		list := make([]domain.Note, 0)
		for _, note := range notes {
			if note.WorkspaceID == workspace.ID {
				list = append(list, note)
			}
		}
		result = append(result, dto.WorkspaceWithNotes{
			Workspace: workspace,
			Notes:     list,
		})
	}

	return result, nil
}

func (w *Workspace) Update(ctx context.Context, input dto.UpdateWorkspaceInput) (domain.Workspace, error) {
	if err := validateInput(input); err != nil {
		return domain.Workspace{}, err
	}

	return w.workspaceRepo.Update(ctx, input.WorkspaceID, input.OwnerID, input.Name)
}

func (w *Workspace) Delete(ctx context.Context, workspaceID, ownerID string) error {
	return w.workspaceRepo.Delete(ctx, workspaceID, ownerID)
}
