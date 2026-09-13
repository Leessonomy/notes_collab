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
		Member: domain.WorkspaceMember{
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
	for _, note := range notes {
		notesByWorkspace[note.WorkspaceID] = append(notesByWorkspace[note.WorkspaceID], dto.NoteOutput{
			ID:          note.ID,
			WorkspaceID: note.WorkspaceID,
			Title:       note.Title,
			Content:     note.Content,
			OwnerID:     note.OwnerID,
			CreatedAt:   note.CreatedAt,
			UpdatedAt:   note.UpdatedAt,
		})
	}

	result := make([]dto.WorkspaceWithNotes, 0, len(workspaces))

	for _, item := range workspaces {
		list := notesByWorkspace[item.Workspace.ID]
		if list == nil {
			list = make([]dto.NoteOutput, 0)
		}

		result = append(result, dto.WorkspaceWithNotes{
			WorkspaceOutput: dto.WorkspaceOutput{
				ID:      item.Workspace.ID,
				Name:    item.Workspace.Name,
				OwnerID: item.Workspace.OwnerID,
				Membership: dto.MembershipOutput{
					Role:     string(item.Member.Role),
					JoinedAt: item.Member.JoinedAt,
				},
				CreatedAt: item.Workspace.CreatedAt,
				UpdatedAt: item.Workspace.UpdatedAt,
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

	member, rights, err := checkAccessRights(ctx, w.workspaceMemberRepo, input.WorkspaceID, input.UserID)
	if err != nil {
		return domain.WorkspaceWithMember{}, err
	}
	if !rights.Manage {
		return domain.WorkspaceWithMember{}, domain.ErrForbidden
	}

	workspace, err := w.workspaceRepo.Update(ctx, input.WorkspaceID, input.Name)
	if err != nil {
		return domain.WorkspaceWithMember{}, err
	}

	return domain.WorkspaceWithMember{Workspace: workspace, Member: member}, nil
}

func (w *Workspace) Delete(ctx context.Context, workspaceID, userID string) error {
	_, rights, err := checkAccessRights(ctx, w.workspaceMemberRepo, workspaceID, userID)
	if err != nil {
		return err
	}
	if !rights.Manage {
		return domain.ErrForbidden
	}

	return w.workspaceRepo.Delete(ctx, workspaceID)
}
