package v1_test

import (
	"context"
	"errors"
	"testing"
	log "user/internal/adapter/logger"
	"user/internal/entity"
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
}

func TestGetUsersBatch(t *testing.T) {
}

func TestGetNewUsers(t *testing.T) {
}

func TestUpdateUser(t *testing.T) {
}

func TestDeleteUser(t *testing.T) {
}
