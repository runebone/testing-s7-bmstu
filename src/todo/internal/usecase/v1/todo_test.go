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

func TestCreateColumn(t *testing.T) {
	runner.Run(t, "TestCreateColumn", func(pt provider.T) {
		tests := []struct {
			name      string
			column    entity.Column
			mockSetup func(mockColumnRepo *mocks.ColumnRepository, column *entity.Column)
			wantErr   bool
			err       error
		}{
			{
				name: "positive",
				column: entity.Column{
					ID:      uuid.New(),
					UserID:  uuid.New(),
					BoardID: uuid.New(),
					Title:   "PositiveColumn",
				},
				mockSetup: func(mockColumnRepo *mocks.ColumnRepository, column *entity.Column) {
					mockColumnRepo.On("CreateColumn", context.Background(), column).Return(nil)
				},
				wantErr: false,
			},
			{
				name: "negative",
				column: entity.Column{
					ID:      uuid.New(),
					UserID:  uuid.New(),
					BoardID: uuid.New(),
					Title:   "NegativeColumn",
				},
				mockSetup: func(mockColumnRepo *mocks.ColumnRepository, column *entity.Column) {
					mockColumnRepo.On("CreateColumn", context.Background(), column).Return(errors.New(""))
				},
				wantErr: true,
				err:     v1.ErrCreateColumn,
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

					tt.mockSetup(mockColumnRepo, &tt.column)

					pt.WithNewStep("Call CreateColumn", func(sCtx provider.StepCtx) {
						err := uc.CreateColumn(context.Background(), &tt.column)

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

func TestCreateCard(t *testing.T) {
	runner.Run(t, "TestCreateCard", func(pt provider.T) {
		tests := []struct {
			name      string
			card      entity.Card
			mockSetup func(mockCardRepo *mocks.CardRepository, card *entity.Card)
			wantErr   bool
			err       error
		}{
			{
				name: "positive",
				card: entity.Card{
					ID:       uuid.New(),
					UserID:   uuid.New(),
					ColumnID: uuid.New(),
					Title:    "PositiveCard",
				},
				mockSetup: func(mockCardRepo *mocks.CardRepository, card *entity.Card) {
					mockCardRepo.On("CreateCard", context.Background(), card).Return(nil)
				},
				wantErr: false,
			},
			{
				name: "negative",
				card: entity.Card{
					ID:       uuid.New(),
					UserID:   uuid.New(),
					ColumnID: uuid.New(),
					Title:    "NegativeCard",
				},
				mockSetup: func(mockCardRepo *mocks.CardRepository, card *entity.Card) {
					mockCardRepo.On("CreateCard", context.Background(), card).Return(errors.New(""))
				},
				wantErr: true,
				err:     v1.ErrCreateCard,
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

					tt.mockSetup(mockCardRepo, &tt.card)

					pt.WithNewStep("Call CreateCard", func(sCtx provider.StepCtx) {
						err := uc.CreateCard(context.Background(), &tt.card)

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
