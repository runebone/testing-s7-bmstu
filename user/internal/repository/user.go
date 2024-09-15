package repository

import (
	"context"
	"user/internal/entity"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetUsersBatch(ctx context.Context, limit, offset int) ([]entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}
