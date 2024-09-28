package usecase

import (
	"auth/internal/dto"
	"context"
)

type AuthUsecase interface {
	Login(ctx context.Context, email, password string) (*dto.LoginResponse, error)
	Refresh(ctx context.Context, refreshToken string) (*dto.RefreshTokenResponse, error)
	// TODO: ValidateToken(ctx context.Context, token string) (string, error) // userID
	// TODO: GetAccessToken(ctx context.Context, userID string) ()
	Logout(ctx context.Context, refreshToken string) error
}
