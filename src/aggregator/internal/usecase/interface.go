package usecase

import (
	"aggregator/internal/entity"
	"context"
	"time"
)

type AggregatorUseCase interface {
	GetStats(ctx context.Context, from, to time.Time) ([]entity.NewUsersAndCardsStats, error)
}
