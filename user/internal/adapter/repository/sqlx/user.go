package repository

import (
	"context"
	"fmt"
	"user/internal/entity"
	"user/internal/repository"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SQLXUserRepository struct {
	db *sqlx.DB
}

func NewSQLXUserRepository(db *sqlx.DB) *SQLXUserRepository {
	return &SQLXUserRepository{db: db}
}

func (r *SQLXUserRepository) CreateUser(ctx context.Context, user *entity.User) error {
	query := `
    INSERT INTO users (id, username, email, password_hash, created_at, updated_at)
    VALUES (:id, :username, :email, :password_hash, :created_at, :updated_at)
    `
	_, err := r.db.NamedExecContext(ctx, query, user)
	return err
}

func (r *SQLXUserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := r.db.GetContext(ctx, &user, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *SQLXUserRepository) GetUsers(ctx context.Context, filter repository.UserFilter) ([]entity.User, error) {
	query := "SELECT * FROM users WHERE 1=1"
	args := []interface{}{}
	i := 1

	if filter.ID != nil {
		str := fmt.Sprintf(" AND id = $%d", i)
		query += str
		args = append(args, *filter.ID)
		i += 1
	}

	if filter.Email != nil {
		str := fmt.Sprintf(" AND email = $%d", i)
		query += str
		args = append(args, *filter.Email)
		i += 1
	}

	if filter.Username != nil {
		str := fmt.Sprintf(" AND username = $%d", i)
		query += str
		args = append(args, *filter.Username)
		i += 1
	}

	var users []entity.User
	stmt, err := r.db.PreparexContext(ctx, query)
	if err != nil {
		return nil, err
	}

	err = stmt.SelectContext(ctx, &users, args...)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *SQLXUserRepository) GetUsersBatch(ctx context.Context, limit, offset int) ([]entity.User, error) {
	var users []entity.User
	err := r.db.SelectContext(ctx, &users, "SELECT * FROM users ORDER BY created_at ASC LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *SQLXUserRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	query := `
    UPDATE users SET username = :username, email = :email, updated_at = :updated_at
    WHERE id = :id
    `
	_, err := r.db.NamedExecContext(ctx, query, user)
	return err
}

func (r *SQLXUserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
	return err
}
