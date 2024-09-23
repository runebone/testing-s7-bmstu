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
	ErrGetBoardByID           = errors.New("failed to get board by id")
	ErrGetBoardsByUser        = errors.New("failed to get boards by user")
	ErrUpdateBoard            = errors.New("failed to update board")
	ErrDeleteBoard            = errors.New("failed to delete board")
	ErrGetColumnByID          = errors.New("failed to get column by id")
	ErrGetColumnsByBoard      = errors.New("failed to get columns by board")
	ErrUpdateColumn           = errors.New("failed to update column")
	ErrDeleteColumn           = errors.New("failed to delete column")
	ErrGetCardByID            = errors.New("failed to get card by id")
	ErrGetCardsByColumn       = errors.New("failed to get cards by column")
	ErrUpdateCard             = errors.New("failed to update card")
	ErrDeleteCard             = errors.New("failed to delete card")
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
	board, err := uc.boardRepo.GetBoardByID(ctx, id)

	if err != nil {
		return nil, ErrGetBoardByID
	}

	return board, nil
}

func (uc *todoUseCase) GetBoardsByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]entity.Board, error) {
	err := validateLimitAndOffset(limit, offset)

	if err != nil {
		return nil, err
	}

	columns, err := uc.boardRepo.GetBoardsByUser(ctx, userID, limit, offset)

	if err != nil {
		return nil, ErrGetBoardsByUser
	}

	return columns, nil
}

func (uc *todoUseCase) UpdateBoard(ctx context.Context, board *entity.Board) error {
	err := validateBoard(board)

	if err != nil {
		return err
	}

	err = uc.boardRepo.UpdateBoard(ctx, board)

	if err != nil {
		return ErrUpdateBoard
	}

	return nil
}

func (uc *todoUseCase) DeleteBoard(ctx context.Context, id uuid.UUID) error {
	err := uc.boardRepo.DeleteBoard(ctx, id)

	if err != nil {
		return ErrDeleteBoard
	}

	return nil
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
	column, err := uc.columnRepo.GetColumnByID(ctx, id)

	if err != nil {
		return nil, ErrGetColumnByID
	}

	return column, nil
}

func (uc *todoUseCase) GetColumnsByBoard(ctx context.Context, boardID uuid.UUID, limit, offset int) ([]entity.Column, error) {
	err := validateLimitAndOffset(limit, offset)

	if err != nil {
		return nil, err
	}

	columns, err := uc.columnRepo.GetColumnsByBoard(ctx, boardID, limit, offset)

	if err != nil {
		return nil, ErrGetColumnsByBoard
	}

	return columns, nil
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

	err = uc.columnRepo.UpdateColumn(ctx, column)

	if err != nil {
		return ErrUpdateColumn
	}

	return nil
}

func (uc *todoUseCase) DeleteColumn(ctx context.Context, id uuid.UUID) error {
	err := uc.columnRepo.DeleteColumn(ctx, id)

	if err != nil {
		return ErrDeleteColumn
	}

	return nil
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
	card, err := uc.cardRepo.GetCardByID(ctx, id)

	if err != nil {
		return nil, ErrGetCardByID
	}

	return card, nil
}

func (uc *todoUseCase) GetCardsByColumn(ctx context.Context, columnID uuid.UUID, limit, offset int) ([]entity.Card, error) {
	err := validateLimitAndOffset(limit, offset)

	if err != nil {
		return nil, err
	}

	cards, err := uc.cardRepo.GetCardsByColumn(ctx, columnID, limit, offset)

	if err != nil {
		return nil, ErrGetCardsByColumn
	}

	return cards, nil
}

func (uc *todoUseCase) UpdateCard(ctx context.Context, card *entity.Card) error {
	err := validateCard(card)

	if err != nil {
		return err
	}

	err = uc.cardRepo.UpdateCard(ctx, card)

	if err != nil {
		return ErrUpdateCard
	}

	return nil
}

func (uc *todoUseCase) DeleteCard(ctx context.Context, id uuid.UUID) error {
	err := uc.cardRepo.DeleteCard(ctx, id)

	if err != nil {
		return ErrDeleteCard
	}

	return nil
}
