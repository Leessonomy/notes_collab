package usecase

import (
	"context"
	"notes-collab-api/internal/domain"
	"notes-collab-api/internal/dto"
	"time"

	"github.com/google/uuid"
)

type WorkspaceStorage interface {
	Create(ctx context.Context, w domain.Workspace) error
	GetByOwner(ctx context.Context, ownerID string) ([]domain.Workspace, error)
}

type Workspace struct {
	workspaces WorkspaceStorage
	notes      NoteStorage
}

func NewWorkspace(workspaces WorkspaceStorage, notes NoteStorage) *Workspace {
	return &Workspace{workspaces: workspaces, notes: notes}
}

func (w *Workspace) Create(ctx context.Context, input dto.CreateWorkspaceInput) (domain.Workspace, error) {
	now := time.Now().UTC()

	workspace := domain.Workspace{
		ID:        uuid.NewString(),
		Name:      input.Name,
		OwnerID:   input.OwnerID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := w.workspaces.Create(ctx, workspace); err != nil {
		return domain.Workspace{}, err
	}

	return workspace, nil
}

func (w *Workspace) Get(ctx context.Context, ownerID string) ([]domain.Workspace, error) {
	return w.workspaces.GetByOwner(ctx, ownerID)
}

func (w *Workspace) GetWorkspacesWithNotes(ctx context.Context, ownerID string) ([]domain.WorkspaceWithNotes, error) {
	workspaces, err := w.workspaces.GetByOwner(ctx, ownerID)

	if err != nil {
		return nil, err
	}

	notes, err := w.notes.GetByOwner(ctx, ownerID)

	if err != nil {
		return nil, err
	}

	result := make([]domain.WorkspaceWithNotes, 0, len(workspaces))

	for _, workspace := range workspaces {
		list := make([]domain.Note, 0)
		for _, note := range notes {
			if note.WorkspaceID == workspace.ID {
				list = append(list, note)
			}
		}
		result = append(result, domain.WorkspaceWithNotes{
			Workspace: workspace,
			Notes:     list,
		})
	}

	return result, nil
}
