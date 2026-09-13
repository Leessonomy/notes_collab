package usecase

import (
	"context"
	"errors"
	"notes-collab-api/internal/domain"
)

type WorkspaceMemberRepo interface {
	Get(ctx context.Context, workspaceID, userID string) (domain.WorkspaceMember, error)
}

func checkAccessRights(ctx context.Context, members WorkspaceMemberRepo, workspaceID, userID string) (domain.WorkspaceMember, domain.Rights, error) {
	m, err := members.Get(ctx, workspaceID, userID)

	if errors.Is(err, domain.ErrNotMember) {
		return domain.WorkspaceMember{}, domain.Rights{}, domain.ErrForbidden
	}
	if err != nil {
		return domain.WorkspaceMember{}, domain.Rights{}, err
	}

	return m, m.Role.Rights(), nil
}
