package usecase

import (
	"context"
	"errors"
	"notes-collab-api/internal/domain"
	"notes-collab-api/internal/dto"
	"time"

	"github.com/google/uuid"
)

type InviteRepo interface {
	Create(ctx context.Context, inv domain.Invite, token domain.InviteToken) error
	ListPending(ctx context.Context, workspaceID string) ([]domain.Invite, error)
	GetPendingByToken(ctx context.Context, token domain.InviteToken) (domain.Invite, error)
	GetWithWorkspace(ctx context.Context, token domain.InviteToken) (domain.InviteWithWorkspace, error)
	Accept(ctx context.Context, inviteID string, m domain.Membership) error
	Delete(ctx context.Context, workspaceID, inviteID string) error
}

type Invite struct {
	inviteRepo     InviteRepo
	membershipRepo MembershipRepo
}

func NewInvite(inviteRepo InviteRepo, membershipRepo MembershipRepo) *Invite {
	return &Invite{inviteRepo: inviteRepo, membershipRepo: membershipRepo}
}

func (i *Invite) Create(ctx context.Context, input dto.CreateInviteInput) (domain.Invite, domain.InviteToken, error) {
	if err := validateInput(&input); err != nil {
		return domain.Invite{}, "", err
	}

	membership, err := getMembership(ctx, i.membershipRepo, input.WorkspaceID, input.UserID)
	if err != nil {
		return domain.Invite{}, "", err
	}
	if !membership.Role.Rights().Manage {
		return domain.Invite{}, "", domain.ErrForbidden
	}

	now := time.Now().UTC()

	invite := domain.Invite{
		ID:          uuid.NewString(),
		WorkspaceID: input.WorkspaceID,
		Role:        input.Role,
		InviterID:   input.UserID,
		ExpiresAt:   now.Add(domain.InviteTTL),
		CreatedAt:   now,
	}
	token := domain.NewInviteToken()

	if err := i.inviteRepo.Create(ctx, invite, token); err != nil {
		return domain.Invite{}, "", err
	}

	return invite, token, nil
}

func (i *Invite) ListPending(ctx context.Context, workspaceID, userID string) ([]domain.Invite, error) {
	membership, err := getMembership(ctx, i.membershipRepo, workspaceID, userID)
	if err != nil {
		return nil, err
	}
	if !membership.Role.Rights().Manage {
		return nil, domain.ErrForbidden
	}

	return i.inviteRepo.ListPending(ctx, workspaceID)
}

func (i *Invite) Revoke(ctx context.Context, workspaceID, inviteID, userID string) error {
	membership, err := getMembership(ctx, i.membershipRepo, workspaceID, userID)
	if err != nil {
		return err
	}
	if !membership.Role.Rights().Manage {
		return domain.ErrForbidden
	}

	return i.inviteRepo.Delete(ctx, workspaceID, inviteID)
}

func (i *Invite) GetByToken(ctx context.Context, token domain.InviteToken) (domain.InviteWithWorkspace, error) {
	return i.inviteRepo.GetWithWorkspace(ctx, token)
}

func (i *Invite) Accept(ctx context.Context, token domain.InviteToken, userID string) (domain.Membership, error) {
	invite, err := i.inviteRepo.GetPendingByToken(ctx, token)
	if err != nil {
		return domain.Membership{}, err
	}

	membership := domain.Membership{
		WorkspaceID: invite.WorkspaceID,
		UserID:      userID,
		Role:        invite.Role,
		JoinedAt:    time.Now().UTC(),
	}

	err = i.inviteRepo.Accept(ctx, invite.ID, membership)
	if err == nil {
		return membership, nil
	}
	if !errors.Is(err, domain.ErrAlreadyMember) {
		return domain.Membership{}, err
	}

	return i.membershipRepo.Get(ctx, invite.WorkspaceID, userID)
}
