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
		&w.OwnerID,
		&w.CreatedAt,
		&w.UpdatedAt,
	)
	if err != nil {
		return domain.Workspace{}, domain.ErrWorkspaceNotFound
	}

	return w, nil
}

func (r *WorkspaceRepo) Create(ctx context.Context, w domain.Workspace) error {
	_, err := r.db.Exec(ctx, `
        INSERT INTO workspaces (
            id, name, owner_id, created_at, updated_at
        )
        VALUES ($1, $2, $3, $4, $5)
    `,
		w.ID,
		w.Name,
		w.OwnerID,
		w.CreatedAt,
		w.UpdatedAt,
	)
	return err
}

func (r *WorkspaceRepo) ListByOwner(ctx context.Context, ownerID string) ([]domain.Workspace, error) {
	rows, err := r.db.Query(ctx, `
        SELECT id, name, owner_id, created_at, updated_at
        FROM workspaces
        WHERE owner_id = $1
    `, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := make([]domain.Workspace, 0)
	for rows.Next() {
		var w domain.Workspace
		if err := rows.Scan(
			&w.ID,
			&w.Name,
			&w.OwnerID,
			&w.CreatedAt,
			&w.UpdatedAt,
		); err != nil {
			return nil, err
		}
		data = append(data, w)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

func (r *WorkspaceRepo) Update(ctx context.Context, workspaceID, ownerID, name string) (domain.Workspace, error) {
	row := r.db.QueryRow(ctx, `
        UPDATE workspaces
        SET name = $1, updated_at = NOW()
        WHERE id = $2 AND owner_id = $3
        RETURNING id, name, owner_id, created_at, updated_at
    `, name, workspaceID, ownerID)

	return scanWorkspace(row)
}

func (r *WorkspaceRepo) Delete(ctx context.Context, workspaceID, ownerID string) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM workspaces WHERE id = $1 AND owner_id = $2`, workspaceID, ownerID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return domain.ErrWorkspaceNotFound
	}

	return nil
}
