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

func (r *WorkspaceRepo) Get(ctx context.Context) ([]domain.Workspace, error) {

	var data []domain.Workspace

	rows, err := r.db.Query(ctx, `
    SELECT *
    FROM workspaces
 `)

	for rows.Next() {
		var w domain.Workspace

		err := rows.Scan(
			&w.ID,
			&w.Name,
			&w.OwnerID,
			&w.CreatedAt,
			&w.UpdatedAt,
		)

		if err != nil {
			return []domain.Workspace{}, err
		}

		data = append(data, w)
	}

	if err != nil {
		return []domain.Workspace{}, err
	}

	return data, nil
}
