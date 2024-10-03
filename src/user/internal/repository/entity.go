package repository

import (
	"time"
	"user/internal/entity"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `db:"id"`
	Username     string    `db:"username"`
	Email        string    `db:"email"`
	Role         string    `db:"role"`
	PasswordHash string    `db:"password_hash"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func RepoUser(e entity.User) User {
	return User{
		ID:           e.ID,
		Username:     e.Username,
		Email:        e.Email,
		Role:         e.Role,
		PasswordHash: e.PasswordHash,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
	}
}

func UserToEntity(r User) entity.User {
	return entity.User{
		ID:           r.ID,
		Username:     r.Username,
		Email:        r.Email,
		Role:         r.Role,
		PasswordHash: r.PasswordHash,
		CreatedAt:    r.CreatedAt,
		UpdatedAt:    r.UpdatedAt,
	}
}
