package v1_test

import (
	log "aggregator/internal/adapter/logger"
	"aggregator/internal/dto"
	"aggregator/mocks"
	"context"
	"errors"
	"testing"
	"time"

	v1 "aggregator/internal/usecase/v1"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
)

func TestGetStats(t *testing.T) {
	runner.Run(t, "TestGetStats", func(pt provider.T) {
		fromTime := time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC)
		toTime := time.Date(2023, 9, 30, 23, 59, 59, 0, time.UTC)

		tests := []struct {
			name      string
			from      time.Time
			to        time.Time
			mockSetup func(mockUserSvc *mocks.UserService, mockTodoSvc *mocks.TodoService, from, to time.Time)
			wantErr   bool
			err       error
		}{
			{
				name: "positive",
				from: fromTime,
				to:   toTime,
				mockSetup: func(mockUserSvc *mocks.UserService, mockTodoSvc *mocks.TodoService, from, to time.Time) {
					userDTOs := make([]dto.User, 3)

					userDTOs[0] = dto.User{
						ID:       uuid.New(),
						Email:    "user@zero.com",
						Username: "UserZero",
					}
					userDTOs[1] = dto.User{
						ID:       uuid.New(),
						Email:    "user@one.com",
						Username: "UserOne",
					}
					userDTOs[2] = dto.User{
						ID:       uuid.New(),
						Email:    "user@two.com",
						Username: "UserTwo",
					}

					mockUserSvc.On("GetNewUsers", context.Background(), from, to).Return(userDTOs, nil)

					cardDTOs := make([]dto.Card, 3)

					cardDTOs[0] = dto.Card{
						ID:       uuid.New(),
						UserID:   uuid.New(),
						ColumnID: uuid.New(),
						Title:    "TitleZero",
					}
					cardDTOs[1] = dto.Card{
						ID:       uuid.New(),
						UserID:   userDTOs[1].ID,
						ColumnID: uuid.New(),
						Title:    "TitleOne",
					}
					cardDTOs[2] = dto.Card{
						ID:       uuid.New(),
						UserID:   userDTOs[2].ID,
						ColumnID: uuid.New(),
						Title:    "TitleTwo",
					}

					mockTodoSvc.On("GetNewCards", context.Background(), from, to).Return(cardDTOs, nil)
				},
				wantErr: false,
			},
			{
				name:      "negative",
				from:      toTime,
				to:        fromTime,
				mockSetup: func(mockUserSvc *mocks.UserService, mockTodoSvc *mocks.TodoService, from, to time.Time) {},
				wantErr:   true,
				err:       v1.ErrInvalidTimeRange,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				runner.Run(t, tt.name, func(pt provider.T) {
					mockUserSvc := new(mocks.UserService)
					mockAuthSvc := new(mocks.AuthService)
					mockTodoSvc := new(mocks.TodoService)
					logger := log.NewEmptyLogger()

					uc := v1.NewAggregatorUseCase(mockUserSvc, mockAuthSvc, mockTodoSvc, logger)

					tt.mockSetup(mockUserSvc, mockTodoSvc, tt.from, tt.to)

					pt.WithNewStep("Call GetStats", func(sCtx provider.StepCtx) {
						_, err := uc.GetStats(context.Background(), tt.from, tt.to)

						if tt.wantErr {
							sCtx.Assert().Error(err, "Expected error")
							sCtx.Assert().ErrorIs(err, tt.err)
						} else {
							sCtx.Assert().NoError(err, "Expected no error")
						}

						mockUserSvc.AssertExpectations(t)
						mockTodoSvc.AssertExpectations(t)
					})
				})
			})
		}
	})
}

func TestRegister(t *testing.T) {
	runner.Run(t, "TestRegister", func(pt provider.T) {
		tests := []struct {
			name      string
			username  string
			email     string
			password  string
			mockSetup func(mockAuthSvc *mocks.AuthService, username, email, password string)
			wantErr   bool
			err       error
		}{
			{
				name:     "positive",
				username: "PositiveUsername",
				email:    "positive@email.com",
				password: "P0s1t1v3P@ssw0rD",
				mockSetup: func(mockAuthSvc *mocks.AuthService, username, email, password string) {
					tokens := dto.Tokens{
						AccessToken:  "AccessToken",
						RefreshToken: "RefreshToken",
					}

					mockAuthSvc.On("Register", context.Background(), username, email, password).Return(&tokens, nil)
				},
				wantErr: false,
			},
			{
				name:     "negative",
				username: "NegativeUsername",
				email:    "negative@email.com",
				password: "N3g@t1v3P@ssw0rD",
				mockSetup: func(mockAuthSvc *mocks.AuthService, username, email, password string) {
					mockAuthSvc.On("Register", context.Background(), username, email, password).Return(nil, errors.New(""))
				},
				wantErr: true,
				err:     v1.ErrRegister,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				runner.Run(t, tt.name, func(pt provider.T) {
					mockUserSvc := new(mocks.UserService)
					mockAuthSvc := new(mocks.AuthService)
					mockTodoSvc := new(mocks.TodoService)
					logger := log.NewEmptyLogger()

					uc := v1.NewAggregatorUseCase(mockUserSvc, mockAuthSvc, mockTodoSvc, logger)

					tt.mockSetup(mockAuthSvc, tt.username, tt.email, tt.password)

					pt.WithNewStep("Call Register", func(sCtx provider.StepCtx) {
						_, err := uc.Register(context.Background(), tt.username, tt.email, tt.password)

						if tt.wantErr {
							sCtx.Assert().Error(err, "Expected error")
							sCtx.Assert().ErrorIs(err, tt.err)
						} else {
							sCtx.Assert().NoError(err, "Expected no error")
						}

						mockAuthSvc.AssertExpectations(t)
					})
				})
			})
		}
	})
}

func TestLogin(t *testing.T) {
	runner.Run(t, "TestLogin", func(pt provider.T) {
		tests := []struct {
			name      string
			email     string
			password  string
			mockSetup func(mockAuthSvc *mocks.AuthService, email, password string)
			wantErr   bool
			err       error
		}{
			{
				name:     "positive",
				email:    "positive@email.com",
				password: "P0s1t1v3P@ssw0rD",
				mockSetup: func(mockAuthSvc *mocks.AuthService, email, password string) {
					tokens := dto.Tokens{
						AccessToken:  "AccessToken",
						RefreshToken: "RefreshToken",
					}

					mockAuthSvc.On("Login", context.Background(), email, password).Return(&tokens, nil)
				},
				wantErr: false,
			},
			{
				name:     "negative",
				email:    "negative@email.com",
				password: "N3g@t1v3P@ssw0rD",
				mockSetup: func(mockAuthSvc *mocks.AuthService, email, password string) {
					mockAuthSvc.On("Login", context.Background(), email, password).Return(nil, errors.New(""))
				},
				wantErr: true,
				err:     v1.ErrLogin,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				runner.Run(t, tt.name, func(pt provider.T) {
					mockUserSvc := new(mocks.UserService)
					mockAuthSvc := new(mocks.AuthService)
					mockTodoSvc := new(mocks.TodoService)
					logger := log.NewEmptyLogger()

					uc := v1.NewAggregatorUseCase(mockUserSvc, mockAuthSvc, mockTodoSvc, logger)

					tt.mockSetup(mockAuthSvc, tt.email, tt.password)

					pt.WithNewStep("Call Login", func(sCtx provider.StepCtx) {
						_, err := uc.Login(context.Background(), tt.email, tt.password)

						if tt.wantErr {
							sCtx.Assert().Error(err, "Expected error")
							sCtx.Assert().ErrorIs(err, tt.err)
						} else {
							sCtx.Assert().NoError(err, "Expected no error")
						}

						mockAuthSvc.AssertExpectations(t)
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
			mockSetup    func(mockAuthSvc *mocks.AuthService, refreshToken string)
			wantErr      bool
			err          error
		}{
			{
				name:         "positive",
				refreshToken: "PositiveToken",
				mockSetup: func(mockAuthSvc *mocks.AuthService, refreshToken string) {
					resp := dto.RefreshResponse{
						AccessToken: "AccessToken",
					}

					mockAuthSvc.On("Refresh", context.Background(), refreshToken).Return(&resp, nil)
				},
				wantErr: false,
			},
			{
				name:         "negative",
				refreshToken: "NegativeToken",
				mockSetup: func(mockAuthSvc *mocks.AuthService, refreshToken string) {
					mockAuthSvc.On("Refresh", context.Background(), refreshToken).Return(nil, errors.New(""))
				},
				wantErr: true,
				err:     v1.ErrRefresh,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				runner.Run(t, tt.name, func(pt provider.T) {
					mockUserSvc := new(mocks.UserService)
					mockAuthSvc := new(mocks.AuthService)
					mockTodoSvc := new(mocks.TodoService)
					logger := log.NewEmptyLogger()

					uc := v1.NewAggregatorUseCase(mockUserSvc, mockAuthSvc, mockTodoSvc, logger)

					tt.mockSetup(mockAuthSvc, tt.refreshToken)

					pt.WithNewStep("Call Refresh", func(sCtx provider.StepCtx) {
						_, err := uc.Refresh(context.Background(), tt.refreshToken)

						if tt.wantErr {
							sCtx.Assert().Error(err, "Expected error")
							sCtx.Assert().ErrorIs(err, tt.err)
						} else {
							sCtx.Assert().NoError(err, "Expected no error")
						}

						mockAuthSvc.AssertExpectations(t)
					})
				})
			})
		}
	})
}

func TestValidate(t *testing.T) {
	runner.Run(t, "TestValidate", func(pt provider.T) {
		tests := []struct {
			name      string
			token     string
			mockSetup func(mockAuthSvc *mocks.AuthService, token string)
			wantErr   bool
			err       error
		}{
			{
				name:  "positive",
				token: "PositiveToken",
				mockSetup: func(mockAuthSvc *mocks.AuthService, token string) {
					resp := dto.ValidateTokenResponse{
						UserID: uuid.New().String(),
						Role:   "User",
					}

					mockAuthSvc.On("ValidateToken", context.Background(), token).Return(&resp, nil)
				},
				wantErr: false,
			},
			{
				name:  "negative",
				token: "NegativeToken",
				mockSetup: func(mockAuthSvc *mocks.AuthService, token string) {
					mockAuthSvc.On("ValidateToken", context.Background(), token).Return(nil, errors.New(""))
				},
				wantErr: true,
				err:     v1.ErrValidate,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				runner.Run(t, tt.name, func(pt provider.T) {
					mockUserSvc := new(mocks.UserService)
					mockAuthSvc := new(mocks.AuthService)
					mockTodoSvc := new(mocks.TodoService)
					logger := log.NewEmptyLogger()

					uc := v1.NewAggregatorUseCase(mockUserSvc, mockAuthSvc, mockTodoSvc, logger)

					tt.mockSetup(mockAuthSvc, tt.token)

					pt.WithNewStep("Call Validate", func(sCtx provider.StepCtx) {
						_, err := uc.Validate(context.Background(), tt.token)

						if tt.wantErr {
							sCtx.Assert().Error(err, "Expected error")
							sCtx.Assert().ErrorIs(err, tt.err)
						} else {
							sCtx.Assert().NoError(err, "Expected no error")
						}

						mockAuthSvc.AssertExpectations(t)
					})
				})
			})
		}
	})
}

func TestLogout(t *testing.T) {
	runner.Run(t, "TestLogout", func(pt provider.T) {
	})
}

func TestGetBoards(t *testing.T) {
	runner.Run(t, "TestGetBoards", func(pt provider.T) {
	})
}

func TestGetColumns(t *testing.T) {
	runner.Run(t, "TestGetColumns", func(pt provider.T) {
	})
}

func TestGetCards(t *testing.T) {
	runner.Run(t, "TestGetCards", func(pt provider.T) {
	})
}

func TestGetCard(t *testing.T) {
	runner.Run(t, "TestGetCard", func(pt provider.T) {
	})
}

func TestCreateBoard(t *testing.T) {
	runner.Run(t, "TestCreateBoard", func(pt provider.T) {
	})
}

func TestCreateColumn(t *testing.T) {
	runner.Run(t, "TestCreateColumn", func(pt provider.T) {
	})
}

func TestCreateCard(t *testing.T) {
	runner.Run(t, "TestCreateCard", func(pt provider.T) {
	})
}

func TestUpdateBoard(t *testing.T) {
	runner.Run(t, "TestUpdateBoard", func(pt provider.T) {
	})
}

func TestUpdateColumn(t *testing.T) {
	runner.Run(t, "TestUpdateColumn", func(pt provider.T) {
	})
}

func TestUpdateCard(t *testing.T) {
	runner.Run(t, "TestUpdateCard", func(pt provider.T) {
	})
}

func TestDeleteBoard(t *testing.T) {
	runner.Run(t, "TestDeleteBoard", func(pt provider.T) {
	})
}

func TestDeleteColumn(t *testing.T) {
	runner.Run(t, "TestDeleteColumn", func(pt provider.T) {
	})
}

func TestDeleteCard(t *testing.T) {
	runner.Run(t, "TestDeleteCard", func(pt provider.T) {
		tests := []struct {
			name      string
			mockSetup func()
			wantErr   bool
			err       error
		}{
			{
				name:      "positive",
				mockSetup: func() {},
				wantErr:   false,
			},
			{
				name:      "negative",
				mockSetup: func() {},
				wantErr:   true,
				// err:       ,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				runner.Run(t, tt.name, func(pt provider.T) {
					// mockUserSvc := new(mocks.UserService)
					// mockAuthSvc := new(mocks.AuthService)
					// mockTodoSvc := new(mocks.TodoService)
					// logger := log.NewEmptyLogger()

					// uc := v1.NewAggregatorUseCase(mockUserSvc, mockAuthSvc, mockTodoSvc, logger)

					// tt.mockSetup()

					pt.WithNewStep("Call TODO", func(sCtx provider.StepCtx) {
						// err := userUC.CreateUser(context.Background(), tt.user)

						if tt.wantErr {
							// sCtx.Assert().Error(err, "Expected error")
							// sCtx.Assert().ErrorIs(err, tt.err)
						} else {
							// sCtx.Assert().NoError(err, "Expected no error")
						}

						// mockRepo.AssertExpectations(t)
					})
				})
			})
		}
	})
}
