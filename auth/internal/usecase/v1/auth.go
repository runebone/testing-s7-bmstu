package v1

import (
	"auth/internal/dto"
	"auth/internal/entity"
	"auth/internal/repository"
	"auth/internal/service/tokengen"
	"auth/internal/service/user"
	"auth/internal/usecase"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type authUseCase struct {
	tokenRepo    repository.TokenRepository
	userService  user.UserService
	tokenService tokengen.TokenService
}

func NewAuthUseCase(tokenRepo repository.TokenRepository, userService user.UserService, tokenService tokengen.TokenService) usecase.AuthUsecase {
	return &authUseCase{
		tokenRepo:    tokenRepo,
		userService:  userService,
		tokenService: tokenService,
	}
}

func (uc *authUseCase) Login(ctx context.Context, email, password string) (*dto.LoginResponse, error) {
	user, err := uc.userService.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if !validatePassword(password, user.PasswordHash) {
		return nil, errors.New("incorrect password")
	}

	userID := user.ID.String()

	accessToken, err := uc.tokenService.GenerateAccessToken(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("couldn't generate access token: %w", err)
	}

	refreshToken, err := uc.tokenService.GenerateRefreshToken(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("couldn't generate refresh token: %w", err)
	}

	token := &entity.Token{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     refreshToken,
		CreatedAt: time.Now(),
	}
	err = uc.tokenRepo.Save(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("couldn't save refresh token: %w", err)
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func validatePassword(password, passwordHash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)) == nil
}

func (uc *authUseCase) Refresh(ctx context.Context, refreshToken string) (*dto.RefreshTokenResponse, error) {
	userID, err := uc.tokenService.ValidateToken(ctx, refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	newAccessToken, err := uc.tokenService.GenerateAccessToken(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("couldn't generate access token: %w", err)
	}

	return &dto.RefreshTokenResponse{
		AccessToken: newAccessToken,
	}, nil
}

func (uc *authUseCase) Logout(ctx context.Context, refreshToken string) error {
	token, err := uc.tokenRepo.FindByToken(ctx, refreshToken)
	if err != nil {
		return errors.New("couldn't find refresh token")
	}

	tokenID := token.ID.String()

	err = uc.tokenRepo.Delete(ctx, tokenID)
	if err != nil {
		return errors.New("couldn't delete refresh token")
	}

	return nil
}
