package http

import (
	"aggregator/internal/dto"
	"aggregator/internal/service/todo"
	"context"
	"net/http"
	"time"
)

type TodoService struct {
	baseURL    string
	httpClient *http.Client
}

func NewTodoService(baseURL string, timeout time.Duration) todo.TodoService {
	return &TodoService{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (s *TodoService) GetNewCards(ctx context.Context, from, to time.Time) ([]dto.Card, error) {
	// TODO:
	return nil, nil
}

func (s *TodoService) GetBoards(ctx context.Context, userID string) ([]dto.Board, error) {
	// TODO:
	return nil, nil
}

func (s *TodoService) GetColumns(ctx context.Context, boardID string) ([]dto.Column, error) {
	// TODO:
	return nil, nil
}

func (s *TodoService) GetCards(ctx context.Context, columnID string) ([]dto.Card, error) {
	// TODO:
	return nil, nil
}

func (s *TodoService) GetCard(ctx context.Context, id string) (*dto.Card, error) {
	// TODO:
	return nil, nil
}

func (s *TodoService) CreateBoard(ctx context.Context, board dto.Board) error {
	// TODO:
	return nil
}

func (s *TodoService) CreateColumn(ctx context.Context, column dto.Column) error {
	// TODO:
	return nil
}

func (s *TodoService) CreateCard(ctx context.Context, card dto.Card) error {
	// TODO:
	return nil
}

func (s *TodoService) UpdateBoard(ctx context.Context, board *dto.Board) error {
	// TODO:
	return nil
}

func (s *TodoService) UpdateColumn(ctx context.Context, column *dto.Column) error {
	// TODO:
	return nil
}

func (s *TodoService) UpdateCard(ctx context.Context, card *dto.Card) error {
	// TODO:
	return nil
}

func (s *TodoService) DeleteBoard(ctx context.Context, id string) error {
	// TODO:
	return nil
}

func (s *TodoService) DeleteColumn(ctx context.Context, id string) error {
	// TODO:
	return nil
}

func (s *TodoService) DeleteCard(ctx context.Context, id string) error {
	// TODO:
	return nil
}
