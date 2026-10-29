package repository

import (
	"context"
	"errors"
	"notes-collab-api/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InviteRepo struct {
	db *pgxpool.Pool
}

func NewInviteRepo(db *pgxpool.Pool) *InviteRepo {
	return &InviteRepo{db: db}
}

func (r *InviteRepo) Create(ctx context.Context, inv domain.Invite, token domain.InviteToken) error {
	_, err := r.db.Exec(ctx, `
        INSERT INTO invites (
            id, workspace_id, token_hash, role, created_by, expires_at, created_at
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `,
		inv.ID,
		inv.WorkspaceID,
		token.Hash(),
		inv.Role,
		inv.InviterID,
		inv.ExpiresAt,
		inv.CreatedAt,
	)
	return err
}

func (r *InviteRepo) ListPending(ctx context.Context, workspaceID string) ([]domain.Invite, error) {
	rows, err := r.db.Query(ctx, `
        SELECT id, workspace_id, role, created_by, expires_at, created_at
        FROM invites
        WHERE workspace_id = $1 AND accepted_at IS NULL AND expires_at > now()
        ORDER BY created_at
    `, workspaceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []domain.Invite
	for rows.Next() {
		inv, err := scanInvite(rows)
		if err != nil {
			return nil, err
		}
		data = append(data, inv)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

func (r *InviteRepo) GetPendingByToken(ctx context.Context, token domain.InviteToken) (domain.Invite, error) {
	inv, err := scanInvite(r.db.QueryRow(ctx, `
        SELECT id, workspace_id, role, created_by, expires_at, created_at
        FROM invites
        WHERE token_hash = $1 AND accepted_at IS NULL AND expires_at > now()
    `, token.Hash()))

	if errors.Is(err, pgx.ErrNoRows) {
		return domain.Invite{}, domain.ErrInviteInvalid
	}
	if err != nil {
		return domain.Invite{}, err
	}

	return inv, nil
}

func (r *InviteRepo) GetWithWorkspace(ctx context.Context, token domain.InviteToken) (domain.InviteWithWorkspace, error) {
	var inv domain.InviteWithWorkspace

	err := r.db.QueryRow(ctx, `
        SELECT i.id, i.workspace_id, i.role, i.created_by, i.expires_at, i.created_at, w.name, u.name
        FROM invites i
        JOIN workspaces w ON w.id = i.workspace_id
        JOIN users u ON u.id = i.created_by
        WHERE i.token_hash = $1 AND i.accepted_at IS NULL AND i.expires_at > now()
    `, token.Hash()).Scan(
		&inv.ID,
		&inv.WorkspaceID,
		&inv.Role,
		&inv.InviterID,
		&inv.ExpiresAt,
		&inv.CreatedAt,
		&inv.WorkspaceName,
		&inv.InviterName,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return domain.InviteWithWorkspace{}, domain.ErrInviteInvalid
	}
	if err != nil {
		return domain.InviteWithWorkspace{}, err
	}

	return inv, nil
}

func (r *InviteRepo) Accept(ctx context.Context, inviteID string, m domain.Membership) error {
	return pgx.BeginFunc(ctx, r.db, func(tx pgx.Tx) error {
		tag, err := tx.Exec(ctx, `
            INSERT INTO workspace_members (workspace_id, user_id, role, joined_at)
            VALUES ($1, $2, $3, $4)
            ON CONFLICT (workspace_id, user_id) DO NOTHING
        `,
			m.WorkspaceID,
			m.UserID,
			m.Role,
			m.JoinedAt,
		)
		if err != nil {
			return err
		}
		if tag.RowsAffected() == 0 {
			return domain.ErrAlreadyMember
		}

		tag, err = tx.Exec(ctx, `
            UPDATE invites
            SET accepted_by = $1, accepted_at = $2
            WHERE id = $3 AND accepted_at IS NULL AND expires_at > now()
        `, m.UserID, m.JoinedAt, inviteID)
		if err != nil {
			return err
		}
		if tag.RowsAffected() == 0 {
			return domain.ErrInviteInvalid
		}

		return nil
	})
}

func (r *InviteRepo) Delete(ctx context.Context, workspaceID, inviteID string) error {
	tag, err := r.db.Exec(ctx, `
        DELETE FROM invites
        WHERE id = $1 AND workspace_id = $2 AND accepted_at IS NULL
    `, inviteID, workspaceID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return domain.ErrInviteNotFound
	}

	return nil
}

func scanInvite(row pgx.Row) (domain.Invite, error) {
	var inv domain.Invite

	err := row.Scan(
		&inv.ID,
		&inv.WorkspaceID,
		&inv.Role,
		&inv.InviterID,
		&inv.ExpiresAt,
		&inv.CreatedAt,
	)
	if err != nil {
		return domain.Invite{}, err
	}

	return inv, nil
}
