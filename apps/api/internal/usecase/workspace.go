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
	GetAll(ctx context.Context) ([]domain.Workspace, error)
}

type Workspace struct {
	workspaces WorkspaceStorage
}

func NewWorkspace(workspaces WorkspaceStorage) *Workspace {
	return &Workspace{workspaces: workspaces}
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

func (w *Workspace) GetAll(ctx context.Context) ([]domain.Workspace, error) {
	return w.workspaces.GetAll(ctx)
}
