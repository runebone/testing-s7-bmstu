package v1

import (
	"context"
	"errors"
	"todo/internal/entity"
	"todo/internal/repository"
	"todo/internal/usecase"

	"github.com/google/uuid"
)

var (
	ErrBoardEmptyTitle        = errors.New("board should have a title")
	ErrBoardNoUserID          = errors.New("board should have a user id")
	ErrColumnEmptyTitle       = errors.New("column should have a title")
	ErrColumnNoUserID         = errors.New("column should have a user id")
	ErrColumnNoBoardID        = errors.New("column should have a board id")
	ErrColumnNegativePosition = errors.New("column cannot have a negative position")
	ErrNegativeLimitOrOffset  = errors.New("limit and offset cannot be negative")
	ErrZeroLimit              = errors.New("limit cannot be zero")
	ErrCardNoUserID           = errors.New("card should have a user id")
	ErrCardNoColumnID         = errors.New("card should have a column id")
	ErrCardNegativePosition   = errors.New("card cannot have a negative position")
	ErrCardEmptyTitle         = errors.New("card should have a title")
)

type todoUseCase struct {
	boardRepo  repository.BoardRepository
	columnRepo repository.ColumnRepository
	cardRepo   repository.CardRepository
}

func NewTodoUseCase(boardRepo repository.BoardRepository, columnRepo repository.ColumnRepository, cardRepo repository.CardRepository) usecase.TodoUseCase {
	return &todoUseCase{
		boardRepo:  boardRepo,
		columnRepo: columnRepo,
		cardRepo:   cardRepo,
	}
}

func (uc *todoUseCase) CreateBoard(ctx context.Context, board *entity.Board) error {
	err := validateBoard(board)
	if err != nil {
		return err
	}

	board.ID = uuid.New()

	return uc.boardRepo.CreateBoard(ctx, board)
}

func validateBoard(board *entity.Board) error {
	if board.Title == "" {
		return ErrBoardEmptyTitle
	}

	if board.UserID == uuid.Nil {
		return ErrBoardNoUserID
	}

	return nil
}

func (uc *todoUseCase) GetBoardByID(ctx context.Context, id uuid.UUID) (*entity.Board, error) {
	return uc.boardRepo.GetBoardByID(ctx, id)
}

func (uc *todoUseCase) UpdateBoard(ctx context.Context, board *entity.Board) error {
	err := validateBoard(board)
	if err != nil {
		return err
	}

	return uc.boardRepo.UpdateBoard(ctx, board)
}

func (uc *todoUseCase) DeleteBoard(ctx context.Context, id uuid.UUID) error {
	return uc.boardRepo.DeleteBoard(ctx, id)
}

func (uc *todoUseCase) CreateColumn(ctx context.Context, column *entity.Column) error {
	err := validateColumn(column)
	if err != nil {
		return err
	}

	column.ID = uuid.New()

	return uc.columnRepo.CreateColumn(ctx, column)
}

func validateColumn(column *entity.Column) error {
	if column.Title == "" {
		return ErrColumnEmptyTitle
	}

	if column.UserID == uuid.Nil {
		return ErrColumnNoUserID
	}

	if column.BoardID == uuid.Nil {
		return ErrColumnNoBoardID
	}

	if column.Position < 0 {
		return ErrColumnNegativePosition
	}

	return nil
}

func (uc *todoUseCase) GetColumnByID(ctx context.Context, id uuid.UUID) (*entity.Column, error) {
	return uc.columnRepo.GetColumnByID(ctx, id)
}

func (uc *todoUseCase) GetColumnsByBoard(ctx context.Context, boardID uuid.UUID, limit, offset int) ([]entity.Column, error) {
	err := validateLimitAndOffset(limit, offset)
	if err != nil {
		return nil, err
	}

	return uc.columnRepo.GetColumnsByBoard(ctx, boardID, limit, offset)
}

func validateLimitAndOffset(limit, offset int) error {
	if limit < 0 || offset < 0 {
		return ErrNegativeLimitOrOffset
	}

	if limit == 0 {
		return ErrZeroLimit
	}

	return nil
}

func (uc *todoUseCase) UpdateColumn(ctx context.Context, column *entity.Column) error {
	err := validateColumn(column)
	if err != nil {
		return err
	}

	return uc.columnRepo.UpdateColumn(ctx, column)
}

func (uc *todoUseCase) DeleteColumn(ctx context.Context, id uuid.UUID) error {
	return uc.columnRepo.DeleteColumn(ctx, id)
}

func (uc *todoUseCase) CreateCard(ctx context.Context, card *entity.Card) error {
	err := validateCard(card)
	if err != nil {
		return err
	}

	card.ID = uuid.New()

	return uc.cardRepo.CreateCard(ctx, card)
}

func validateCard(card *entity.Card) error {
	if card.UserID == uuid.Nil {
		return ErrCardNoUserID
	}

	if card.ColumnID == uuid.Nil {
		return ErrCardNoColumnID
	}

	if card.Position < 0 {
		return ErrCardNegativePosition
	}

	if card.Title == "" {
		return ErrCardEmptyTitle
	}

	return nil
}

func (uc *todoUseCase) GetCardByID(ctx context.Context, id uuid.UUID) (*entity.Card, error) {
	return uc.cardRepo.GetCardByID(ctx, id)
}

func (uc *todoUseCase) GetCardsByColumn(ctx context.Context, columnID uuid.UUID, limit, offset int) ([]entity.Card, error) {
	err := validateLimitAndOffset(limit, offset)
	if err != nil {
		return nil, err
	}

	return uc.cardRepo.GetCardsByColumn(ctx, columnID, limit, offset)
}

func (uc *todoUseCase) UpdateCard(ctx context.Context, card *entity.Card) error {
	err := validateCard(card)
	if err != nil {
		return err
	}

	return uc.cardRepo.UpdateCard(ctx, card)
}

func (uc *todoUseCase) DeleteCard(ctx context.Context, id uuid.UUID) error {
	return uc.cardRepo.DeleteCard(ctx, id)
}
