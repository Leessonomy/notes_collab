package repository

import (
	"context"
	"notes-collab-api/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type WorkspaceRepo struct {
	db *pgxpool.Pool
}

func NewWorkspaceRepo(db *pgxpool.Pool) *WorkspaceRepo {
	return &WorkspaceRepo{db: db}
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

func (r *WorkspaceRepo) GetByOwner(ctx context.Context, ownerID string) ([]domain.Workspace, error) {
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
