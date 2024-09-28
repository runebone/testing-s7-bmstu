package user

import (
	"auth/internal/dto"
	"context"
)

type UserService interface {
	GetUserByEmail(ctx context.Context, email string) (*dto.User, error)
}
