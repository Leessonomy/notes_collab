package usecase

import (
	"context"
	"errors"
	"notes-collab-api/internal/domain"
	"notes-collab-api/internal/dto"
)

type MembershipRepo interface {
	Get(ctx context.Context, workspaceID, userID string) (domain.Membership, error)
	List(ctx context.Context, workspaceID string) ([]domain.Member, error)
	SetRole(ctx context.Context, workspaceID, userID string, role domain.Role) (domain.Membership, error)
	Delete(ctx context.Context, workspaceID, userID string) error
}

type Member struct {
	membershipRepo MembershipRepo
}

func NewMember(membershipRepo MembershipRepo) *Member {
	return &Member{membershipRepo: membershipRepo}
}

func (m *Member) List(ctx context.Context, workspaceID, userID string) ([]domain.Member, error) {
	membership, err := getMembership(ctx, m.membershipRepo, workspaceID, userID)
	if err != nil {
		return nil, err
	}
	if !membership.Role.Rights().Read {
		return nil, domain.ErrForbidden
	}

	return m.membershipRepo.List(ctx, workspaceID)
}

func (m *Member) ChangeRole(ctx context.Context, input dto.ChangeMemberRoleInput) (domain.Membership, error) {
	if err := validateInput(&input); err != nil {
		return domain.Membership{}, err
	}

	membership, err := getMembership(ctx, m.membershipRepo, input.WorkspaceID, input.UserID)
	if err != nil {
		return domain.Membership{}, err
	}
	if !membership.Role.Rights().Manage {
		return domain.Membership{}, domain.ErrForbidden
	}

	target, err := m.membershipRepo.Get(ctx, input.WorkspaceID, input.MemberID)
	if err != nil {
		return domain.Membership{}, err
	}
	if target.Role == domain.Owner {
		return domain.Membership{}, domain.ErrForbidden
	}

	return m.membershipRepo.SetRole(ctx, input.WorkspaceID, input.MemberID, input.Role)
}

func (m *Member) Remove(ctx context.Context, workspaceID, memberID, userID string) error {
	membership, err := getMembership(ctx, m.membershipRepo, workspaceID, userID)
	if err != nil {
		return err
	}
	if memberID != userID && !membership.Role.Rights().Manage {
		return domain.ErrForbidden
	}

	target, err := m.membershipRepo.Get(ctx, workspaceID, memberID)
	if err != nil {
		return err
	}
	if target.Role == domain.Owner {
		return domain.ErrForbidden
	}

	return m.membershipRepo.Delete(ctx, workspaceID, memberID)
}

func getMembership(ctx context.Context, repo MembershipRepo, workspaceID, userID string) (domain.Membership, error) {
	m, err := repo.Get(ctx, workspaceID, userID)

	if errors.Is(err, domain.ErrNotMember) {
		return domain.Membership{}, domain.ErrForbidden
	}
	if err != nil {
		return domain.Membership{}, err
	}

	return m, nil
}
