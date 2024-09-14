package usecase

import (
	"context"
	"user/internal/entities"
	"user/internal/repositories"

	"github.com/google/uuid"
)

type UserUseCase interface {
	CreateUser(ctx context.Context, user *entities.User) error
	GetUserByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	UpdateUser(ctx context.Context, user *entities.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type userUseCase struct {
	repo repositories.UserRepository
}

func NewUserUseCase(repo repositories.UserRepository) UserUseCase {
	return &userUseCase{
		repo: repo,
	}
}

func (u *userUseCase) CreateUser(ctx context.Context, user *entities.User) error {
	return u.repo.CreateUser(ctx, user)
}

func (u *userUseCase) GetUserByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	return u.repo.GetUserByID(ctx, id)
}

func (u *userUseCase) UpdateUser(ctx context.Context, user *entities.User) error {
	return u.repo.UpdateUser(ctx, user)
}

func (u *userUseCase) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return u.repo.DeleteUser(ctx, id)
}
