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
	ListByMember(ctx context.Context, userID string) ([]domain.WorkspaceWithMember, error)
	Update(ctx context.Context, workspaceID, name string) (domain.Workspace, error)
	Delete(ctx context.Context, workspaceID string) error
}

type Workspace struct {
	workspaceRepo       WorkspaceRepo
	noteRepo            NoteRepo
	workspaceMemberRepo WorkspaceMemberRepo
}

func NewWorkspace(workspaceRepo WorkspaceRepo, noteRepo NoteRepo, workspaceMemberRepo WorkspaceMemberRepo) *Workspace {
	return &Workspace{workspaceRepo: workspaceRepo, noteRepo: noteRepo, workspaceMemberRepo: workspaceMemberRepo}
}

func (w *Workspace) Create(ctx context.Context, input dto.CreateWorkspaceInput) (domain.WorkspaceWithMember, error) {
	if err := validateInput(&input); err != nil {
		return domain.WorkspaceWithMember{}, err
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
		return domain.WorkspaceWithMember{}, err
	}

	return domain.WorkspaceWithMember{
		Workspace: workspace,
		Membership: domain.WorkspaceMember{
			WorkspaceID: workspace.ID,
			UserID:      input.OwnerID,
			Role:        domain.Owner,
			JoinedAt:    now,
		},
	}, nil
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
				ID:      ws.ID,
				Name:    ws.Name,
				OwnerID: ws.OwnerID,
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

func (w *Workspace) Update(ctx context.Context, input dto.UpdateWorkspaceInput) (domain.WorkspaceWithMember, error) {
	if err := validateInput(&input); err != nil {
		return domain.WorkspaceWithMember{}, err
	}

	member, err := getMember(ctx, w.workspaceMemberRepo, input.WorkspaceID, input.UserID)
	if err != nil {
		return domain.WorkspaceWithMember{}, err
	}
	if !member.Role.Rights().Manage {
		return domain.WorkspaceWithMember{}, domain.ErrForbidden
	}

	workspace, err := w.workspaceRepo.Update(ctx, input.WorkspaceID, input.Name)
	if err != nil {
		return domain.WorkspaceWithMember{}, err
	}

	return domain.WorkspaceWithMember{Workspace: workspace, Membership: member}, nil
}

func (w *Workspace) Delete(ctx context.Context, workspaceID, userID string) error {
	member, err := getMember(ctx, w.workspaceMemberRepo, workspaceID, userID)
	if err != nil {
		return err
	}
	if !member.Role.Rights().Manage {
		return domain.ErrForbidden
	}

	return w.workspaceRepo.Delete(ctx, workspaceID)
}
