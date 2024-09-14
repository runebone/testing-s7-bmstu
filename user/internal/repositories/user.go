package repositories

import (
	"context"
	"user/internal/entities"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entities.User) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	UpdateUser(ctx context.Context, user *entities.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}
