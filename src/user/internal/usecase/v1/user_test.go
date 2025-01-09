package v1_test

import (
	"context"
	"errors"
	"testing"
	"time"
	log "user/internal/adapter/logger"
	"user/internal/entity"
	"user/internal/repository"
	"user/internal/testdata"
	v1 "user/internal/usecase/v1"
	"user/mocks"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	"github.com/stretchr/testify/mock"
)

func TestCreateUser(t *testing.T) {
	runner.Run(t, "Test CreateUser", func(pt provider.T) {
		objectMother := &testdata.UserObjectMother{}

		tests := []struct {
			name      string
			user      entity.User
			mockSetup func(mockRepo *mocks.UserRepository)
			wantErr   bool
		}{
			{
				name: "positive",
				user: objectMother.ValidUser(),
				mockSetup: func(mockRepo *mocks.UserRepository) {
					mockRepo.On("CreateUser", context.Background(), mock.Anything).Return(nil)
				},
				wantErr: false,
			},
			{
				name:      "negative",
				user:      objectMother.InvalidEmailUser(),
				mockSetup: func(mockRepo *mocks.UserRepository) {},
				wantErr:   true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				runner.Run(t, tt.name, func(pt provider.T) {
					mockRepo := new(mocks.UserRepository)
					logger := log.NewEmptyLogger()
					userUC := v1.NewUserUseCase(mockRepo, logger)

					tt.mockSetup(mockRepo)

					pt.WithNewStep("Call CreateUser", func(sCtx provider.StepCtx) {
						err := userUC.CreateUser(context.Background(), tt.user)

						if tt.wantErr {
							sCtx.Assert().Error(err, "Expected error")
						} else {
							sCtx.Assert().NoError(err, "Expected no error")
						}

						mockRepo.AssertExpectations(t)
					})
				})
			})
		}
	})
}

func TestGetUserByID(t *testing.T) {
	runner.Run(t, "Test GetUserByID", func(pt provider.T) {
		tests := []struct {
			name      string
			id        uuid.UUID
			mockSetup func(mockRepo *mocks.UserRepository, id uuid.UUID)
			wantErr   bool
		}{
			{
				name: "positive",
				id:   uuid.New(),
				mockSetup: func(mockRepo *mocks.UserRepository, id uuid.UUID) {
					user := testdata.NewUserBuilder().
						WithUsername("PositiveUser").
						WithID(id).
						Build()
					mockRepo.On("GetUserByID", context.Background(), id).Return(&user, nil)
				},
				wantErr: false,
			},
			{
				name: "negative",
				id:   uuid.New(),
				mockSetup: func(mockRepo *mocks.UserRepository, id uuid.UUID) {
					mockRepo.On("GetUserByID", context.Background(), id).Return(nil, errors.New(""))
				},
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				runner.Run(t, tt.name, func(pt provider.T) {
					mockRepo := new(mocks.UserRepository)
					logger := log.NewEmptyLogger()
					userUC := v1.NewUserUseCase(mockRepo, logger)

					tt.mockSetup(mockRepo, tt.id)

					pt.WithNewStep("Call GetUserByID", func(sCtx provider.StepCtx) {
						_, err := userUC.GetUserByID(context.Background(), tt.id)

						if tt.wantErr {
							sCtx.Assert().Error(err, "Expected error")
							sCtx.Assert().ErrorIs(err, v1.ErrGetUserByID, "Expected ErrGetUserByID")
						} else {
							sCtx.Assert().NoError(err, "Expected no error")
						}

						mockRepo.AssertExpectations(t)
					})
				})
			})
		}
	})
}

func TestGetUsers(t *testing.T) {
	runner.Run(t, "Test GetUsers", func(pt provider.T) {
		objectMother := &testdata.UserObjectMother{}

		tests := []struct {
			name      string
			filter    repository.UserFilter
			mockSetup func(mockRepo *mocks.UserRepository, filter repository.UserFilter)
			wantErr   bool
		}{
			{
				name:   "positive",
				filter: objectMother.PositiveUsernameFilter(),
				mockSetup: func(mockRepo *mocks.UserRepository, filter repository.UserFilter) {
					user := testdata.NewUserBuilder().
						WithUsername(*filter.Username).
						Build()

					users := []entity.User{user}

					mockRepo.On("GetUsers", context.Background(), filter).Return(users, nil)
				},
				wantErr: false,
			},
			{
				name:   "negative",
				filter: objectMother.NegativeUsernameFilter(),
				mockSetup: func(mockRepo *mocks.UserRepository, filter repository.UserFilter) {
					mockRepo.On("GetUsers", context.Background(), filter).Return(nil, errors.New(""))
				},
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				runner.Run(t, tt.name, func(pt provider.T) {
					mockRepo := new(mocks.UserRepository)
					logger := log.NewEmptyLogger()
					userUC := v1.NewUserUseCase(mockRepo, logger)

					tt.mockSetup(mockRepo, tt.filter)

					pt.WithNewStep("Call GetUsers", func(sCtx provider.StepCtx) {
						_, err := userUC.GetUsers(context.Background(), tt.filter)

						if tt.wantErr {
							sCtx.Assert().Error(err, "Expected error")
							sCtx.Assert().ErrorIs(err, v1.ErrGetUsers, "Expected ErrGetUsers")
						} else {
							sCtx.Assert().NoError(err, "Expected no error")
						}

						mockRepo.AssertExpectations(t)
					})
				})
			})
		}
	})
}

func TestGetUsersBatch(t *testing.T) {
	runner.Run(t, "Test GetUsersBatch", func(pt provider.T) {
		tests := []struct {
			name      string
			limit     int
			offset    int
			mockSetup func(mockRepo *mocks.UserRepository, limit, offset int)
			wantErr   bool
		}{
			{
				name:   "positive",
				limit:  10,
				offset: 0,
				mockSetup: func(mockRepo *mocks.UserRepository, limit, offset int) {
					user1 := testdata.NewUserBuilder().
						WithUsername("User1").
						Build()

					user2 := testdata.NewUserBuilder().
						WithUsername("User2").
						Build()

					users := []entity.User{user1, user2}

					mockRepo.On("GetUsersBatch", context.Background(), limit, offset).Return(users, nil)
				},
				wantErr: false,
			},
			{
				name:   "negative",
				limit:  10,
				offset: 0,
				mockSetup: func(mockRepo *mocks.UserRepository, limit, offset int) {
					mockRepo.On("GetUsersBatch", context.Background(), limit, offset).Return(nil, errors.New(""))
				},
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				runner.Run(t, tt.name, func(pt provider.T) {
					mockRepo := new(mocks.UserRepository)
					logger := log.NewEmptyLogger()
					userUC := v1.NewUserUseCase(mockRepo, logger)

					tt.mockSetup(mockRepo, tt.limit, tt.offset)

					pt.WithNewStep("Call GetUsersBatch", func(sCtx provider.StepCtx) {
						_, err := userUC.GetUsersBatch(context.Background(), tt.limit, tt.offset)

						if tt.wantErr {
							sCtx.Assert().Error(err, "Expected error")
							sCtx.Assert().ErrorIs(err, v1.ErrGetUsersBatch, "Expected ErrGetUsers")
						} else {
							sCtx.Assert().NoError(err, "Expected no error")
						}

						mockRepo.AssertExpectations(t)
					})
				})
			})
		}
	})
}

func TestGetNewUsers(t *testing.T) {
	runner.Run(t, "Test GetUsersBatch", func(pt provider.T) {
		fromTime := time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC)
		toTime := time.Date(2023, 9, 30, 23, 59, 59, 0, time.UTC)

		tests := []struct {
			name      string
			from      time.Time
			to        time.Time
			mockSetup func(mockRepo *mocks.UserRepository, from, to time.Time)
			wantErr   bool
		}{
			{
				name: "positive",
				from: fromTime,
				to:   toTime,
				mockSetup: func(mockRepo *mocks.UserRepository, from, to time.Time) {
					user1 := testdata.NewUserBuilder().
						WithUsername("User1").
						Build()

					user2 := testdata.NewUserBuilder().
						WithUsername("User2").
						Build()

					users := []entity.User{user1, user2}

					mockRepo.On("GetNewUsers", context.Background(), from, to).Return(users, nil)
				},
				wantErr: false,
			},
			{
				name: "negative",
				from: fromTime,
				to:   toTime,
				mockSetup: func(mockRepo *mocks.UserRepository, from, to time.Time) {
					mockRepo.On("GetNewUsers", context.Background(), from, to).Return(nil, errors.New(""))
				},
				wantErr: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				runner.Run(t, tt.name, func(pt provider.T) {
					mockRepo := new(mocks.UserRepository)
					logger := log.NewEmptyLogger()
					userUC := v1.NewUserUseCase(mockRepo, logger)

					tt.mockSetup(mockRepo, tt.from, tt.to)

					pt.WithNewStep("Call GetNewUsers", func(sCtx provider.StepCtx) {
						_, err := userUC.GetNewUsers(context.Background(), tt.from, tt.to)

						if tt.wantErr {
							sCtx.Assert().Error(err, "Expected error")
							sCtx.Assert().ErrorIs(err, v1.ErrGetNewUsers, "Expected ErrGetNewUsers")
						} else {
							sCtx.Assert().NoError(err, "Expected no error")
						}

						mockRepo.AssertExpectations(t)
					})
				})
			})
		}
	})
}

func TestUpdateUser(t *testing.T) {
	runner.Run(t, "Test UpdateUser", func(pt provider.T) {
		objectMother := &testdata.UserObjectMother{}

		tests := []struct {
			name      string
			user      entity.User
			mockSetup func(mockRepo *mocks.UserRepository, user *entity.User)
			wantErr   bool
			err       error
		}{
			{
				name: "positive",
				user: objectMother.ValidUser(),
				mockSetup: func(mockRepo *mocks.UserRepository, user *entity.User) {
					mockRepo.On("GetUserByID", context.Background(), user.ID).Return(user, nil)
					mockRepo.On("UpdateUser", context.Background(), user).Return(nil)
				},
				wantErr: false,
			},
			{
				name: "negative",
				user: objectMother.ValidUser(),
				mockSetup: func(mockRepo *mocks.UserRepository, user *entity.User) {
					mockRepo.On("GetUserByID", context.Background(), user.ID).Return(user, nil)
					mockRepo.On("UpdateUser", context.Background(), user).Return(errors.New(""))
				},
				wantErr: true,
				err:     v1.ErrUpdateUser,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				runner.Run(t, tt.name, func(pt provider.T) {
					mockRepo := new(mocks.UserRepository)
					logger := log.NewEmptyLogger()
					userUC := v1.NewUserUseCase(mockRepo, logger)

					tt.mockSetup(mockRepo, &tt.user)

					pt.WithNewStep("Call UpdateUser", func(sCtx provider.StepCtx) {
						err := userUC.UpdateUser(context.Background(), &tt.user)

						if tt.wantErr {
							sCtx.Assert().Error(err, "Expected error")
							sCtx.Assert().ErrorIs(err, tt.err)
						} else {
							sCtx.Assert().NoError(err, "Expected no error")
						}

						mockRepo.AssertExpectations(t)
					})
				})
			})
		}
	})
}

func TestDeleteUser(t *testing.T) {
}
