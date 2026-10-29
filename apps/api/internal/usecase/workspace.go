package usecase

import (
	"context"
	"notes-collab-api/internal/domain"
	"notes-collab-api/internal/dto"
	"time"

	"github.com/google/uuid"
)

type WorkspaceRepo interface {
	Create(ctx context.Context, w domain.WorkspaceWithMembership) error
	ListByMember(ctx context.Context, userID string) ([]domain.WorkspaceWithMembership, error)
	Update(ctx context.Context, workspaceID, name string) (domain.Workspace, error)
	Delete(ctx context.Context, workspaceID string) error
}

type Workspace struct {
	workspaceRepo  WorkspaceRepo
	noteRepo       NoteRepo
	membershipRepo MembershipRepo
}

func NewWorkspace(workspaceRepo WorkspaceRepo, noteRepo NoteRepo, membershipRepo MembershipRepo) *Workspace {
	return &Workspace{workspaceRepo: workspaceRepo, noteRepo: noteRepo, membershipRepo: membershipRepo}
}

func (w *Workspace) Create(ctx context.Context, input dto.CreateWorkspaceInput) (domain.WorkspaceWithMembership, error) {
	if err := validateInput(&input); err != nil {
		return domain.WorkspaceWithMembership{}, err
	}

	now := time.Now().UTC()

	workspaceID := uuid.NewString()

	workspace := domain.WorkspaceWithMembership{
		Workspace: domain.Workspace{
			ID:        workspaceID,
			Name:      input.Name,
			CreatedAt: now,
			UpdatedAt: now,
		},
		Membership: domain.Membership{
			WorkspaceID: workspaceID,
			UserID:      input.OwnerID,
			Role:        domain.Owner,
			JoinedAt:    now,
		},
	}

	if err := w.workspaceRepo.Create(ctx, workspace); err != nil {
		return domain.WorkspaceWithMembership{}, err
	}

	return workspace, nil
}

func (w *Workspace) ListWithNotes(ctx context.Context, userID string) ([]dto.WorkspaceWithNotes, error) {
	workspaces, err := w.workspaceRepo.ListByMember(ctx, userID)
	if err != nil {
		return nil, err
	}

	notes, err := w.noteRepo.ListByMember(ctx, userID)
	if err != nil {
		return nil, err
	}

	notesByWorkspace := make(map[string][]dto.NoteOutput, len(workspaces))

	for _, n := range notes {
		notesByWorkspace[n.WorkspaceID] = append(notesByWorkspace[n.WorkspaceID], dto.NoteOutput{
			ID:          n.ID,
			WorkspaceID: n.WorkspaceID,
			Title:       n.Title,
			Content:     n.Content,
			OwnerID:     n.OwnerID,
			CreatedAt:   n.CreatedAt,
			UpdatedAt:   n.UpdatedAt,
		})
	}

	result := make([]dto.WorkspaceWithNotes, 0, len(workspaces))

	for _, ws := range workspaces {
		list := notesByWorkspace[ws.ID]

		if len(list) == 0 {
			list = make([]dto.NoteOutput, 0)
		}

		result = append(result, dto.WorkspaceWithNotes{
			WorkspaceOutput: dto.WorkspaceOutput{
				ID:   ws.ID,
				Name: ws.Name,
				Membership: dto.MembershipOutput{
					Role:     ws.Membership.Role,
					JoinedAt: ws.Membership.JoinedAt,
				},
				CreatedAt: ws.CreatedAt,
				UpdatedAt: ws.UpdatedAt,
			},
			Notes: list,
		})
	}

	return result, nil
}

func (w *Workspace) Update(ctx context.Context, input dto.UpdateWorkspaceInput) (domain.WorkspaceWithMembership, error) {
	if err := validateInput(&input); err != nil {
		return domain.WorkspaceWithMembership{}, err
	}

	membership, err := getMembership(ctx, w.membershipRepo, input.WorkspaceID, input.UserID)
	if err != nil {
		return domain.WorkspaceWithMembership{}, err
	}
	if !membership.Role.Rights().Manage {
		return domain.WorkspaceWithMembership{}, domain.ErrForbidden
	}

	workspace, err := w.workspaceRepo.Update(ctx, input.WorkspaceID, input.Name)
	if err != nil {
		return domain.WorkspaceWithMembership{}, err
	}

	return domain.WorkspaceWithMembership{Workspace: workspace, Membership: membership}, nil
}

func (w *Workspace) Delete(ctx context.Context, workspaceID, userID string) error {
	membership, err := getMembership(ctx, w.membershipRepo, workspaceID, userID)
	if err != nil {
		return err
	}
	if !membership.Role.Rights().Manage {
		return domain.ErrForbidden
	}

	return w.workspaceRepo.Delete(ctx, workspaceID)
}
