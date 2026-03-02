package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/shadowpr1est/knigapoisk-auth-service/internal/domain/entity"
	"github.com/shadowpr1est/knigapoisk-auth-service/internal/domain/repository"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) repository.UserRepository {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users(email, username, password_hash, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowxContext(ctx, query,
		user.Email,
		user.Username,
		user.PasswordHash,
		user.Role,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepo) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	var user entity.User
	query := `
		SELECT id, email, username, password_hash, role, created_at, updated_at
		FROM users WHERE id = $1 
	`
	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	query := `
		SELECT id, email, username, password_hash, role, created_at, updated_at
		FROM users WHERE email = $1
	`
	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		return nil, err
	}
	return &user, err
}

func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	query := `
		SELECT id, email, username, password_hash, role, created_at, updated_at
		FROM users WHERE username = $1
	`
	err := r.db.GetContext(ctx, &user, query, username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) Update(ctx context.Context, user *entity.User) error {
	query := `
		UPDATE users
		SET email=$1, username=$2, password_hash=$3, updated_at=now()
		WHERE id =$4
	`
	_, err := r.db.ExecContext(ctx, query,
		user.Email,
		user.Username,
		user.PasswordHash,
		user.ID)

	return err
}

