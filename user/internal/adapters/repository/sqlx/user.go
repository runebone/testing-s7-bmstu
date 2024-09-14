package repositories

import (
	"context"
	"user/internal/entities"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SQLXUserRepository struct {
	db *sqlx.DB
}

func NewSQLXUserRepository(db *sqlx.DB) *SQLXUserRepository {
	return &SQLXUserRepository{db: db}
}

func (r *SQLXUserRepository) CreateUser(ctx context.Context, user *entities.User) error {
	query := `
    INSERT INTO users (id, username, email, password_hash, created_at, updated_at)
    VALUES (:id, :username, :email, :password_hash, :created_at, :updated_at)
    `
	_, err := r.db.NamedExecContext(ctx, query, user)
	return err
}

func (r *SQLXUserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	var user entities.User
	err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *SQLXUserRepository) UpdateUser(ctx context.Context, user *entities.User) error {
	query := `
    UPDATE users SET username = :username, email = :email, password_hash = :password_hash, updated_at = :updated_at
    WHERE id = :id
    `
	_, err := r.db.NamedExecContext(ctx, query, user)
	return err
}

func (r *SQLXUserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
	return err
}
