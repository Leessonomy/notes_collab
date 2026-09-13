package repository

import (
	"context"
	"errors"
	"notes-collab-api/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WorkspaceMemberRepo struct {
	db *pgxpool.Pool
}

func NewWorkspaceMemberRepo(db *pgxpool.Pool) *WorkspaceMemberRepo {
	return &WorkspaceMemberRepo{db: db}
}

func (r *WorkspaceMemberRepo) Get(ctx context.Context, workspaceID, userID string) (domain.WorkspaceMember, error) {
	var (
		m    = domain.WorkspaceMember{WorkspaceID: workspaceID, UserID: userID}
		role string
	)

	err := r.db.QueryRow(ctx,
		`SELECT role, joined_at FROM workspace_members WHERE workspace_id = $1 AND user_id = $2`,
		workspaceID,
		userID,
	).Scan(&role, &m.JoinedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return domain.WorkspaceMember{}, domain.ErrNotMember
	}
	if err != nil {
		return domain.WorkspaceMember{}, err
	}

	m.Role = domain.Role(role)
	return m, nil
}
