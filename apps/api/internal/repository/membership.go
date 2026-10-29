package repository

import (
	"context"
	"errors"
	"notes-collab-api/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MembershipRepo struct {
	db *pgxpool.Pool
}

func NewMembershipRepo(db *pgxpool.Pool) *MembershipRepo {
	return &MembershipRepo{db: db}
}

func (r *MembershipRepo) Get(ctx context.Context, workspaceID, userID string) (domain.Membership, error) {
	m := domain.Membership{WorkspaceID: workspaceID, UserID: userID}

	err := r.db.QueryRow(ctx,
		`SELECT role, joined_at FROM workspace_members WHERE workspace_id = $1 AND user_id = $2`,
		workspaceID,
		userID,
	).Scan(&m.Role, &m.JoinedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Membership{}, domain.ErrNotMember
	}
	if err != nil {
		return domain.Membership{}, err
	}

	return m, nil
}

func (r *MembershipRepo) List(ctx context.Context, workspaceID string) ([]domain.Member, error) {
	rows, err := r.db.Query(ctx, `
        SELECT m.user_id, u.name, u.email, m.role, m.joined_at
        FROM workspace_members m
        JOIN users u ON u.id = m.user_id
        WHERE m.workspace_id = $1
        ORDER BY m.joined_at
    `, workspaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []domain.Member
	for rows.Next() {
		var m domain.Member
		if err := rows.Scan(
			&m.UserID,
			&m.Name,
			&m.Email,
			&m.Role,
			&m.JoinedAt,
		); err != nil {
			return nil, err
		}
		data = append(data, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

func (r *MembershipRepo) SetRole(ctx context.Context, workspaceID, userID string, role domain.Role) (domain.Membership, error) {
	m := domain.Membership{WorkspaceID: workspaceID, UserID: userID}

	err := r.db.QueryRow(ctx, `
        UPDATE workspace_members
        SET role = $1
        WHERE workspace_id = $2 AND user_id = $3
        RETURNING role, joined_at
    `, role, workspaceID, userID).Scan(&m.Role, &m.JoinedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Membership{}, domain.ErrNotMember
	}
	if err != nil {
		return domain.Membership{}, err
	}

	return m, nil
}

func (r *MembershipRepo) Delete(ctx context.Context, workspaceID, userID string) error {
	tag, err := r.db.Exec(ctx,
		`DELETE FROM workspace_members WHERE workspace_id = $1 AND user_id = $2`,
		workspaceID,
		userID,
	)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return domain.ErrNotMember
	}

	return nil
}
