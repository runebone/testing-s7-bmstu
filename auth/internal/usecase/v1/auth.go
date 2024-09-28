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
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrIncorrectPassword    error = errors.New("incorrect password")
	ErrGetAccessToken       error = errors.New("couldn't get access token")
	ErrGenerateRefreshToken error = errors.New("couldn't generate refresh token")
	ErrSaveRefreshToken     error = errors.New("couldn't save refresh token")
	ErrInvalidRefreshToken  error = errors.New("invalid refresh token")
	ErrGenerateAccessToken  error = errors.New("couldn't generate access token")
	ErrFindRefreshToken     error = errors.New("couldn't find refresh token")
	ErrDeleteRefreshToken   error = errors.New("couldn't delete refresh token")
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
		return nil, ErrIncorrectPassword
	}

	userID := user.ID.String()

	accessToken, err := uc.tokenService.GenerateAccessToken(ctx, userID)
	if err != nil {
		return nil, ErrGetAccessToken
	}

	refreshToken, err := uc.tokenService.GenerateRefreshToken(ctx, userID)
	if err != nil {
		return nil, ErrGenerateRefreshToken
	}

	token := &entity.Token{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     refreshToken,
		CreatedAt: time.Now(),
	}
	err = uc.tokenRepo.Save(ctx, token)
	if err != nil {
		return nil, ErrSaveRefreshToken
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
		return nil, ErrInvalidRefreshToken
	}

	newAccessToken, err := uc.tokenService.GenerateAccessToken(ctx, userID)
	if err != nil {
		return nil, ErrGenerateAccessToken
	}

	return &dto.RefreshTokenResponse{
		AccessToken: newAccessToken,
	}, nil
}

func (uc *authUseCase) Logout(ctx context.Context, refreshToken string) error {
	token, err := uc.tokenRepo.FindByToken(ctx, refreshToken)
	if err != nil {
		return ErrFindRefreshToken
	}

	tokenID := token.ID.String()

	err = uc.tokenRepo.Delete(ctx, tokenID)
	if err != nil {
		return ErrDeleteRefreshToken
	}

	return nil
}
