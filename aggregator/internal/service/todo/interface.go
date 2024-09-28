package todo

import (
	"aggregator/internal/dto"
	"context"
	"time"
)

type TodoService interface {
	GetNewCards(ctx context.Context, from time.Time, to time.Time) ([]dto.Card, error)
}
