package repository

import (
	"context"
	"notes-collab-api/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NoteRepo struct {
	db *pgxpool.Pool
}

func NewNoteRepo(db *pgxpool.Pool) *NoteRepo {
	return &NoteRepo{db: db}
}

func (r *NoteRepo) Create(ctx context.Context, n domain.Note) error {
	_, err := r.db.Exec(ctx, `
        INSERT INTO notes (
            id, workspace_id, title, content, owner_id, created_at, updated_at
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `,
		n.ID,
		n.WorkspaceID,
		n.Title,
		n.Content,
		n.OwnerID,
		n.CreatedAt,
		n.UpdatedAt,
	)

	return err
}

func scanNotes(rows pgx.Rows) ([]domain.Note, error) {
	defer rows.Close()

	data := make([]domain.Note, 0)
	for rows.Next() {
		var n domain.Note
		if err := rows.Scan(
			&n.ID,
			&n.WorkspaceID,
			&n.Title,
			&n.Content,
			&n.OwnerID,
			&n.CreatedAt,
			&n.UpdatedAt,
		); err != nil {
			return nil, err
		}
		data = append(data, n)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

func (r *NoteRepo) GetByOwner(ctx context.Context, ownerID string) ([]domain.Note, error) {
	rows, err := r.db.Query(ctx, `
        SELECT id, workspace_id, title, content, owner_id, created_at, updated_at
        FROM notes
        WHERE owner_id = $1
    `, ownerID)
	if err != nil {
		return nil, err
	}
	return scanNotes(rows)
}

func (r *NoteRepo) Delete(ctx context.Context, noteID, ownerID string) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM notes WHERE id = $1 AND owner_id = $2`, noteID, ownerID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return domain.ErrNoteNotFound
	}

	return nil
}
