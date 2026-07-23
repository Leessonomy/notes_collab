package repository

import (
	"context"
	"errors"
	"notes-collab-api/internal/domain"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RefreshTokenRepo struct {
	db *pgxpool.Pool
}

func NewRefreshTokenRepo(db *pgxpool.Pool) *RefreshTokenRepo {
	return &RefreshTokenRepo{db: db}
}

func (r *RefreshTokenRepo) Save(ctx context.Context, t domain.RefreshToken) error {
	_, err := r.db.Exec(ctx, `
        INSERT INTO refresh_tokens (token, user_id, expires_at, created_at)
        VALUES ($1, $2, $3, $4)
    `,
		t.Token,
		t.UserID,
		t.ExpiresAt,
		t.CreatedAt,
	)
	return err
}

func (r *RefreshTokenRepo) Refresh(ctx context.Context, oldToken, newToken string, expiresAt time.Time) (string, error) {
	var userID string

	err := r.db.QueryRow(ctx, `
        UPDATE refresh_tokens
        SET token = $1, expires_at = $2, created_at = now()
        WHERE token = $3 AND expires_at > now()
        RETURNING user_id
    `, newToken, expiresAt, oldToken).Scan(&userID)

	if errors.Is(err, pgx.ErrNoRows) {
		return "", domain.ErrTokenNotFound
	}
	if err != nil {
		return "", err
	}

	return userID, nil
}

func (r *RefreshTokenRepo) Delete(ctx context.Context, token string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM refresh_tokens WHERE token = $1`, token)
	return err
}

func (r *RefreshTokenRepo) DeleteByUser(ctx context.Context, userId string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM refresh_tokens WHERE user_id = $1`, userId)
	return err
}
