package usecase

import (
	"auth/internal/dto"
	"context"
)

type AuthUsecase interface {
	Login(ctx context.Context, username, password string) (*dto.LoginResponse, error)
	Refresh(ctx context.Context, refreshToken string) (*dto.RefreshTokenResponse, error)
	Logout(ctx context.Context, refreshToken string) error
}
