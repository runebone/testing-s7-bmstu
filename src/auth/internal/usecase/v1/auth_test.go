package v1_test

import (
	log "auth/internal/adapter/logger"
	"auth/internal/dto"
	"auth/mocks"
	"context"
	"errors"
	"testing"

	v1 "auth/internal/usecase/v1"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestRegister(t *testing.T) {
	runner.Run(t, "TestRegister", func(pt provider.T) {
		tests := []struct {
			name         string
			username     string
			email        string
			password     string
			accessToken  string
			refreshToken string
			mockSetup    func(mockTokenRepo *mocks.TokenRepository, mockUserSvc *mocks.UserService, mockTokenSvc *mocks.TokenService, username, email, password, accessToken, refreshToken string)
			wantErr      bool
			err          error
		}{
			{
				name:         "positive",
				username:     "PositiveUsername",
				email:        "positive@email.com",
				password:     "P0s1t1v3P@ssw0rD",
				accessToken:  "PositiveAccessToken",
				refreshToken: "PositiveRefreshToken",
				mockSetup: func(mockTokenRepo *mocks.TokenRepository, mockUserSvc *mocks.UserService, mockTokenSvc *mocks.TokenService, username, email, password, accessToken, refreshToken string) {
					mockUserSvc.On("CreateUser", context.Background(), username, email, password).Return(nil)

					userID := uuid.New()
					role := "user"
					pwdHash, _ := hashPassword(password)
					userDTO := &dto.User{
						ID:           userID,
						Username:     username,
						Email:        email,
						Role:         role,
						PasswordHash: pwdHash,
					}

					mockUserSvc.On("GetUserByEmail", context.Background(), email).Return(userDTO, nil)
					mockTokenSvc.On("GenerateAccessToken", context.Background(), userID.String(), role).Return(accessToken, nil)
					mockTokenSvc.On("GenerateRefreshToken", context.Background(), userID.String(), role).Return(refreshToken, nil)
					mockTokenRepo.On("Save", context.Background(), mock.Anything).Return(nil)
				},
				wantErr: false,
			},
			{
				name:         "negative",
				username:     "NegativeUsername",
				email:        "negative@email.com",
				password:     "N3g@t1v3P@ssw0rD",
				accessToken:  "NegativeAccessToken",
				refreshToken: "NegativeRefreshToken",
				mockSetup: func(mockTokenRepo *mocks.TokenRepository, mockUserSvc *mocks.UserService, mockTokenSvc *mocks.TokenService, username, email, password, accessToken, refreshToken string) {
					mockUserSvc.On("CreateUser", context.Background(), username, email, password).Return(errors.New(""))
				},
				wantErr: true,
				err:     v1.ErrCreateUser,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				runner.Run(t, tt.name, func(pt provider.T) {
					mockTokenRepo := new(mocks.TokenRepository)
					mockUserSvc := new(mocks.UserService)
					mockTokenSvc := new(mocks.TokenService)
					logger := log.NewEmptyLogger()

					uc := v1.NewAuthUseCase(mockTokenRepo, mockUserSvc, mockTokenSvc, logger)

					tt.mockSetup(mockTokenRepo, mockUserSvc, mockTokenSvc, tt.username, tt.email, tt.password, tt.accessToken, tt.refreshToken)

					pt.WithNewStep("Call Register", func(sCtx provider.StepCtx) {
						_, err := uc.Register(context.Background(), tt.username, tt.email, tt.password)

						if tt.wantErr {
							sCtx.Assert().Error(err, "Expected error")
							sCtx.Assert().ErrorIs(err, tt.err)
						} else {
							sCtx.Assert().NoError(err, "Expected no error")
						}

						mockUserSvc.AssertExpectations(t)
					})
				})
			})
		}
	})
}

func TestLogin(t *testing.T) {
	runner.Run(t, "TestLogin", func(pt provider.T) {
		tests := []struct {
			name         string
			email        string
			password     string
			accessToken  string
			refreshToken string
			mockSetup    func(mockTokenRepo *mocks.TokenRepository, mockUserSvc *mocks.UserService, mockTokenSvc *mocks.TokenService, email, password, accessToken, refreshToken string)
			wantErr      bool
			err          error
		}{
			{
				name:         "positive",
				email:        "positive@email.com",
				password:     "P0s1t1v3P@ssw0rD",
				accessToken:  "PositiveAccessToken",
				refreshToken: "PositiveRefreshToken",
				mockSetup: func(mockTokenRepo *mocks.TokenRepository, mockUserSvc *mocks.UserService, mockTokenSvc *mocks.TokenService, email, password, accessToken, refreshToken string) {
					userID := uuid.New()
					role := "user"
					pwdHash, _ := hashPassword(password)
					userDTO := &dto.User{
						ID:           userID,
						Email:        email,
						Role:         role,
						PasswordHash: pwdHash,
					}

					mockUserSvc.On("GetUserByEmail", context.Background(), email).Return(userDTO, nil)
					mockTokenSvc.On("GenerateAccessToken", context.Background(), userID.String(), role).Return(accessToken, nil)
					mockTokenSvc.On("GenerateRefreshToken", context.Background(), userID.String(), role).Return(refreshToken, nil)
					mockTokenRepo.On("Save", context.Background(), mock.Anything).Return(nil)
				},
				wantErr: false,
			},
			{
				name:         "negative",
				email:        "negative@email.com",
				password:     "N3g@t1v3P@ssw0rD",
				accessToken:  "NegativeAccessToken",
				refreshToken: "NegativeRefreshToken",
				mockSetup: func(mockTokenRepo *mocks.TokenRepository, mockUserSvc *mocks.UserService, mockTokenSvc *mocks.TokenService, email, password, accessToken, refreshToken string) {
					mockUserSvc.On("GetUserByEmail", context.Background(), email).Return(nil, errors.New(""))
				},
				wantErr: true,
				err:     v1.ErrGetUserByEmail,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				runner.Run(t, tt.name, func(pt provider.T) {
					mockTokenRepo := new(mocks.TokenRepository)
					mockUserSvc := new(mocks.UserService)
					mockTokenSvc := new(mocks.TokenService)
					logger := log.NewEmptyLogger()

					uc := v1.NewAuthUseCase(mockTokenRepo, mockUserSvc, mockTokenSvc, logger)

					tt.mockSetup(mockTokenRepo, mockUserSvc, mockTokenSvc, tt.email, tt.password, tt.accessToken, tt.refreshToken)

					pt.WithNewStep("Call Login", func(sCtx provider.StepCtx) {
						_, err := uc.Login(context.Background(), tt.email, tt.password)

						if tt.wantErr {
							sCtx.Assert().Error(err, "Expected error")
							sCtx.Assert().ErrorIs(err, tt.err)
						} else {
							sCtx.Assert().NoError(err, "Expected no error")
						}

						mockUserSvc.AssertExpectations(t)
					})
				})
			})
		}
	})
}

func TestRefresh(t *testing.T) {
	runner.Run(t, "TestRefresh", func(pt provider.T) {
		tests := []struct {
			name         string
			refreshToken string
			mockSetup    func(mockTokenRepo *mocks.TokenRepository, mockUserSvc *mocks.UserService, mockTokenSvc *mocks.TokenService, refreshToken string)
			wantErr      bool
			err          error
		}{
			{
				name:         "positive",
				refreshToken: "PositiveRefreshToken",
				mockSetup: func(mockTokenRepo *mocks.TokenRepository, mockUserSvc *mocks.UserService, mockTokenSvc *mocks.TokenService, refreshToken string) {
					userID := uuid.New()
					role := "user"
					accessToken := "PositiveAccessToken"

					mockTokenSvc.On("ValidateToken", context.Background(), refreshToken).Return(userID.String(), role, nil)
					mockTokenSvc.On("GenerateAccessToken", context.Background(), userID.String(), role).Return(accessToken, nil)
				},
				wantErr: false,
			},
			{
				name:         "negative",
				refreshToken: "NegativeRefreshToken",
				mockSetup: func(mockTokenRepo *mocks.TokenRepository, mockUserSvc *mocks.UserService, mockTokenSvc *mocks.TokenService, refreshToken string) {
					mockTokenSvc.On("ValidateToken", context.Background(), refreshToken).Return("", "", errors.New(""))
				},
				wantErr: true,
				err:     v1.ErrValidateToken,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				runner.Run(t, tt.name, func(pt provider.T) {
					mockTokenRepo := new(mocks.TokenRepository)
					mockUserSvc := new(mocks.UserService)
					mockTokenSvc := new(mocks.TokenService)
					logger := log.NewEmptyLogger()

					uc := v1.NewAuthUseCase(mockTokenRepo, mockUserSvc, mockTokenSvc, logger)

					tt.mockSetup(mockTokenRepo, mockUserSvc, mockTokenSvc, tt.refreshToken)

					pt.WithNewStep("Call Refresh", func(sCtx provider.StepCtx) {
						_, err := uc.Refresh(context.Background(), tt.refreshToken)

						if tt.wantErr {
							sCtx.Assert().Error(err, "Expected error")
							sCtx.Assert().ErrorIs(err, tt.err)
						} else {
							sCtx.Assert().NoError(err, "Expected no error")
						}

						mockUserSvc.AssertExpectations(t)
					})
				})
			})
		}
	})
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
