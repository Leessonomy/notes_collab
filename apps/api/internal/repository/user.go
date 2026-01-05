package repository

import (
	"context"
	"notes-collab-api/internal/domain"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, u domain.User) error {
	_, err := r.db.Exec(ctx, `
        INSERT INTO users (id, name, email, password, created_at)
        VALUES ($1, $2, $3, $4, $5)
    `,
		u.ID,
		u.Name,
		u.Email,
		u.Password,
		u.CreatedAt,
	)
	return err
}

func (r *UserRepo) scanOne(row pgx.Row) (*domain.User, error) {
	var u domain.User

	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.CreatedAt)

	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	return &u, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	return r.scanOne(r.db.QueryRow(ctx,
		`SELECT id, name, email, password FROM users WHERE email = $1`,
		email,
	))
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	return r.scanOne(r.db.QueryRow(ctx,
		`SELECT id, name, email, password FROM users WHERE id = $1`,
		id,
	))
}
