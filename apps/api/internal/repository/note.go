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

func scanNote(row pgx.Row) (domain.Note, error) {
	var n domain.Note

	err := row.Scan(
		&n.ID,
		&n.WorkspaceID,
		&n.Title,
		&n.Content,
		&n.OwnerID,
		&n.CreatedAt,
		&n.UpdatedAt,
	)
	if err != nil {
		return domain.Note{}, domain.ErrNoteNotFound
	}

	return n, nil
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

func (r *NoteRepo) Update(ctx context.Context, noteID string, title, content *string) (domain.Note, error) {
	row := r.db.QueryRow(ctx, `
        UPDATE notes
        SET title = COALESCE($1, title),
            content = COALESCE($2, content),
            updated_at = NOW()
        WHERE id = $3
        RETURNING id, workspace_id, title, content, owner_id, created_at, updated_at
    `, title, content, noteID)

	return scanNote(row)
}

func (r *NoteRepo) ListByMember(ctx context.Context, userID string) ([]domain.Note, error) {
	rows, err := r.db.Query(ctx, `
        SELECT n.id, n.workspace_id, n.title, n.content, n.owner_id, n.created_at, n.updated_at
        FROM notes n
        JOIN workspace_members m ON m.workspace_id = n.workspace_id
        WHERE m.user_id = $1
        ORDER BY n.created_at
    `, userID)

	if err != nil {
		return nil, err
	}
	return scanNotes(rows)
}

func (r *NoteRepo) GetByID(ctx context.Context, noteID string) (domain.Note, error) {
	row := r.db.QueryRow(ctx, `
        SELECT id, workspace_id, title, content, owner_id, created_at, updated_at
        FROM notes
        WHERE id = $1
    `, noteID)

	return scanNote(row)
}

func (r *NoteRepo) Delete(ctx context.Context, noteID string) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM notes WHERE id = $1`, noteID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return domain.ErrNoteNotFound
	}

	return nil
}
