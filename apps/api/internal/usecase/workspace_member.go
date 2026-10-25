package usecase

import (
	"context"
	"errors"
	"notes-collab-api/internal/domain"
)

type WorkspaceMemberRepo interface {
	Get(ctx context.Context, workspaceID, userID string) (domain.WorkspaceMember, error)
}

func getMember(ctx context.Context, members WorkspaceMemberRepo, workspaceID, userID string) (domain.WorkspaceMember, error) {
	m, err := members.Get(ctx, workspaceID, userID)

	if errors.Is(err, domain.ErrNotMember) {
		return domain.WorkspaceMember{}, domain.ErrForbidden
	}
	if err != nil {
		return domain.WorkspaceMember{}, err
	}

	return m, nil
}
