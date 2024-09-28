package v1

import (
	"aggregator/internal/dto"
	"aggregator/internal/entity"
	"aggregator/internal/service/auth"
	"aggregator/internal/service/todo"
	"aggregator/internal/service/user"
	"aggregator/internal/usecase"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	dateLayout          string = "2006-01-02"
	ErrInvalidTimeRange error  = errors.New("<<from>> should be not greater than <<to>>")
	ErrGetNewUsers      error  = errors.New("failed to get new users")
	ErrGetNewCards      error  = errors.New("failed to get new cards")
)

type AggregatorUseCase struct {
	userSvc user.UserService
	authSvc auth.AuthService
	todoSvc todo.TodoService
}

func NewAggregatorUseCase(userSvc user.UserService, authSvc auth.AuthService, todoSvc todo.TodoService) usecase.AggregatorUseCase {
	return &AggregatorUseCase{
		userSvc: userSvc,
		authSvc: authSvc,
		todoSvc: todoSvc,
	}
}

func (uc *AggregatorUseCase) GetStats(ctx context.Context, from, to time.Time) ([]entity.NewUsersAndCardsStats, error) {
	if from.After(to) {
		return nil, ErrInvalidTimeRange
	}

	users, err := uc.userSvc.GetNewUsers(ctx, from, to)
	if err != nil {
		return nil, ErrGetNewUsers
	}

	dateUsersMap, newUsersMap := getDateUsersAndNewUsersMap(users)

	cards, err := uc.todoSvc.GetNewCards(ctx, from, to)
	if err != nil {
		return nil, ErrGetNewCards
	}

	dateCardsMap, numCardsByNewUsersMap := getDateCardsAndNumCardsByNewUsersMap(cards, newUsersMap)

	stats := getStats(dateUsersMap, dateCardsMap, numCardsByNewUsersMap)

	return stats, nil
}

func getDateUsersAndNewUsersMap(users []dto.User) (map[string][]entity.User, map[uuid.UUID]dto.User) {
	dateUsersMap := map[string][]entity.User{}
	newUsersMap := map[uuid.UUID]dto.User{}

	for _, userDTO := range users {
		dateKey := userDTO.CreatedAt.Format(dateLayout)
		user := dto.UserToEntity(&userDTO)

		dateUsersMap[dateKey] = append(dateUsersMap[dateKey], user)

		newUsersMap[user.ID] = userDTO
	}

	return dateUsersMap, newUsersMap
}

func getDateCardsAndNumCardsByNewUsersMap(cards []dto.Card, newUsersMap map[uuid.UUID]dto.User) (map[string][]entity.Card, map[string]int) {
	dateCardsMap := map[string][]entity.Card{}
	numCardsByNewUsersMap := map[string]int{}

	for _, cardDTO := range cards {
		dateKey := cardDTO.CreatedAt.Format(dateLayout)
		card := dto.CardToEntity(&cardDTO)

		dateCardsMap[dateKey] = append(dateCardsMap[dateKey], card)

		if user, ok := newUsersMap[card.UserID]; ok && user.CreatedAt == cardDTO.CreatedAt {
			numCardsByNewUsersMap[dateKey]++
		}
	}

	return dateCardsMap, numCardsByNewUsersMap
}

func getStats(dateUsersMap map[string][]entity.User, dateCardsMap map[string][]entity.Card, numCardsByNewUsersMap map[string]int) []entity.NewUsersAndCardsStats {
	stats := []entity.NewUsersAndCardsStats{}

	dates := mergeKeys(dateUsersMap, dateCardsMap)

	for dateKey := range dates {
		date, _ := time.Parse(dateLayout, dateKey)

		stat := entity.NewUsersAndCardsStats{
			Date:               date,
			Users:              dateUsersMap[dateKey],
			Cards:              dateCardsMap[dateKey],
			NumCardsByNewUsers: numCardsByNewUsersMap[dateKey],
		}

		stats = append(stats, stat)
	}

	return stats
}

func mergeKeys(dateUsersMap map[string][]entity.User, dateCardsMap map[string][]entity.Card) map[string]bool {
	dates := map[string]bool{}

	for k := range dateUsersMap {
		dates[k] = true
	}

	for k := range dateCardsMap {
		dates[k] = true
	}

	return dates
}
