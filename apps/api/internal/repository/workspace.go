package repository

import (
	"context"
	"notes-collab-api/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WorkspaceRepo struct {
	db *pgxpool.Pool
}

func NewWorkspaceRepo(db *pgxpool.Pool) *WorkspaceRepo {
	return &WorkspaceRepo{db: db}
}

func scanWorkspace(row pgx.Row) (domain.Workspace, error) {
	var w domain.Workspace

	err := row.Scan(
		&w.ID,
		&w.Name,
		&w.CreatedAt,
		&w.UpdatedAt,
	)
	if err != nil {
		return domain.Workspace{}, domain.ErrWorkspaceNotFound
	}

	return w, nil
}

func (r *WorkspaceRepo) Create(ctx context.Context, w domain.WorkspaceWithMembership) error {
	return pgx.BeginFunc(ctx, r.db, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `
            INSERT INTO workspaces (id, name, created_at, updated_at)
            VALUES ($1, $2, $3, $4)
        `,
			w.ID,
			w.Name,
			w.CreatedAt,
			w.UpdatedAt,
		)
		if err != nil {
			return err
		}

		_, err = tx.Exec(ctx, `
            INSERT INTO workspace_members (workspace_id, user_id, role, joined_at)
            VALUES ($1, $2, $3, $4)
        `,
			w.ID,
			w.Membership.UserID,
			w.Membership.Role,
			w.Membership.JoinedAt,
		)
		return err
	})
}

func (r *WorkspaceRepo) ListByMember(ctx context.Context, userID string) ([]domain.WorkspaceWithMembership, error) {
	rows, err := r.db.Query(ctx, `
        SELECT w.id, w.name, w.created_at, w.updated_at, m.role, m.joined_at
        FROM workspaces w
        JOIN workspace_members m ON m.workspace_id = w.id
        WHERE m.user_id = $1
        ORDER BY w.created_at
    `, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []domain.WorkspaceWithMembership
	for rows.Next() {
		item := domain.WorkspaceWithMembership{Membership: domain.Membership{UserID: userID}}
		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.CreatedAt,
			&item.UpdatedAt,
			&item.Membership.Role,
			&item.Membership.JoinedAt,
		); err != nil {
			return nil, err
		}
		item.Membership.WorkspaceID = item.ID
		data = append(data, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

func (r *WorkspaceRepo) Update(ctx context.Context, workspaceID, name string) (domain.Workspace, error) {
	row := r.db.QueryRow(ctx, `
        UPDATE workspaces
        SET name = $1, updated_at = NOW()
        WHERE id = $2
        RETURNING id, name, created_at, updated_at
    `, name, workspaceID)

	return scanWorkspace(row)
}

func (r *WorkspaceRepo) Delete(ctx context.Context, workspaceID string) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM workspaces WHERE id = $1`, workspaceID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return domain.ErrWorkspaceNotFound
	}

	return nil
}
