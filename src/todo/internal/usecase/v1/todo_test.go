package v1_test

import (
	"context"
	"errors"
	"testing"
	log "todo/internal/adapter/logger"
	"todo/internal/entity"
	"todo/mocks"

	v1 "todo/internal/usecase/v1"

	"github.com/google/uuid"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
)

func TestCreateBoard(t *testing.T) {
	runner.Run(t, "TestCreateBoard", func(pt provider.T) {
		tests := []struct {
			name      string
			board     entity.Board
			mockSetup func(mockBoardRepo *mocks.BoardRepository, board *entity.Board)
			wantErr   bool
			err       error
		}{
			{
				name: "positive",
				board: entity.Board{
					ID:     uuid.New(),
					UserID: uuid.New(),
					Title:  "PositiveBoard",
				},
				mockSetup: func(mockBoardRepo *mocks.BoardRepository, board *entity.Board) {
					mockBoardRepo.On("CreateBoard", context.Background(), board).Return(nil)
				},
				wantErr: false,
			},
			{
				name: "negative",
				board: entity.Board{
					ID:     uuid.New(),
					UserID: uuid.New(),
					Title:  "NegativeBoard",
				},
				mockSetup: func(mockBoardRepo *mocks.BoardRepository, board *entity.Board) {
					mockBoardRepo.On("CreateBoard", context.Background(), board).Return(errors.New(""))
				},
				wantErr: true,
				err:     v1.ErrCreateBoard,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				runner.Run(t, tt.name, func(pt provider.T) {
					mockBoardRepo := new(mocks.BoardRepository)
					mockColumnRepo := new(mocks.ColumnRepository)
					mockCardRepo := new(mocks.CardRepository)
					logger := log.NewEmptyLogger()

					uc := v1.NewTodoUseCase(mockBoardRepo, mockColumnRepo, mockCardRepo, logger)

					tt.mockSetup(mockBoardRepo, &tt.board)

					pt.WithNewStep("Call CreateBoard", func(sCtx provider.StepCtx) {
						err := uc.CreateBoard(context.Background(), &tt.board)

						if tt.wantErr {
							sCtx.Assert().Error(err, "Expected error")
							sCtx.Assert().ErrorIs(err, tt.err)
						} else {
							sCtx.Assert().NoError(err, "Expected no error")
						}

						mockCardRepo.AssertExpectations(t)
					})
				})
			})
		}
	})
}
