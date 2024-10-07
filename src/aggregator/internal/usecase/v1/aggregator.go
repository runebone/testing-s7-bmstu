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
	ErrRegister         error  = errors.New("failed to register")
	ErrLogin            error  = errors.New("failed to login")
	ErrRefresh          error  = errors.New("failed to refresh")
	ErrValidate         error  = errors.New("failed to validate token")
	ErrLogout           error  = errors.New("failed to logout")
	ErrGetBoards        error  = errors.New("failed to get boards")
	ErrGetColumns       error  = errors.New("failed to get columns")
	ErrGetCards         error  = errors.New("failed to get cards")
	ErrGetCard          error  = errors.New("failed to get card")
	ErrCreateBoard      error  = errors.New("failed to create board")
	ErrCreateColumn     error  = errors.New("failed to create column")
	ErrCreateCard       error  = errors.New("failed to create card")
	ErrUpdateBoard      error  = errors.New("failed to update board")
	ErrUpdateColumn     error  = errors.New("failed to update column")
	ErrUpdateCard       error  = errors.New("failed to update card")
	ErrDeleteBoard      error  = errors.New("failed to delete board")
	ErrDeleteColumn     error  = errors.New("failed to delete column")
	ErrDeleteCard       error  = errors.New("failed to delete card")
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

func (uc *AggregatorUseCase) Register(ctx context.Context, username, email, password string) (*dto.Tokens, error) {
	tokens, err := uc.authSvc.Register(ctx, username, email, password)
	if err != nil {
		return nil, ErrRegister
	}
	return tokens, nil
}

func (uc *AggregatorUseCase) Login(ctx context.Context, email, password string) (*dto.Tokens, error) {
	tokens, err := uc.authSvc.Login(ctx, email, password)
	if err != nil {
		return nil, ErrLogin
	}
	return tokens, nil
}

func (uc *AggregatorUseCase) Refresh(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error) {
	resp, err := uc.authSvc.Refresh(ctx, refreshToken)
	if err != nil {
		return nil, ErrRefresh
	}
	return resp, nil
}

func (uc *AggregatorUseCase) Validate(ctx context.Context, token string) (*dto.ValidateTokenResponse, error) {
	resp, err := uc.authSvc.ValidateToken(ctx, token)
	if err != nil {
		return nil, ErrValidate
	}
	return resp, nil
}

func (uc *AggregatorUseCase) Logout(ctx context.Context, refreshToken string) error {
	err := uc.authSvc.Logout(ctx, refreshToken)
	if err != nil {
		return ErrLogout
	}
	return nil
}

func (uc *AggregatorUseCase) GetBoards(ctx context.Context, userID string) ([]dto.Board, error) {
	boards, err := uc.todoSvc.GetBoards(ctx, userID)
	if err != nil {
		return nil, ErrGetBoards
	}
	return boards, nil
}

func (uc *AggregatorUseCase) GetColumns(ctx context.Context, boardID string) ([]dto.Column, error) {
	columns, err := uc.todoSvc.GetColumns(ctx, boardID)
	if err != nil {
		return nil, ErrGetColumns
	}
	return columns, nil
}

func (uc *AggregatorUseCase) GetCards(ctx context.Context, columnID string) ([]dto.Card, error) {
	cards, err := uc.todoSvc.GetCards(ctx, columnID)
	if err != nil {
		return nil, ErrGetCards
	}
	return cards, nil
}

func (uc *AggregatorUseCase) GetCard(ctx context.Context, id string) (*dto.Card, error) {
	card, err := uc.todoSvc.GetCard(ctx, id)
	if err != nil {
		return nil, ErrGetCard
	}
	return card, nil
}

func (uc *AggregatorUseCase) CreateBoard(ctx context.Context, board dto.Board) error {
	err := uc.todoSvc.CreateBoard(ctx, board)
	if err != nil {
		return ErrCreateBoard
	}
	return nil
}

func (uc *AggregatorUseCase) CreateColumn(ctx context.Context, column dto.Column) error {
	err := uc.todoSvc.CreateColumn(ctx, column)
	if err != nil {
		return ErrCreateColumn
	}
	return nil
}

func (uc *AggregatorUseCase) CreateCard(ctx context.Context, card dto.Card) error {
	err := uc.todoSvc.CreateCard(ctx, card)
	if err != nil {
		return ErrCreateCard
	}
	return nil
}

func (uc *AggregatorUseCase) UpdateBoard(ctx context.Context, board *dto.Board) error {
	err := uc.todoSvc.UpdateBoard(ctx, board)
	if err != nil {
		return ErrUpdateBoard
	}
	return nil
}

func (uc *AggregatorUseCase) UpdateColumn(ctx context.Context, column *dto.Column) error {
	err := uc.todoSvc.UpdateColumn(ctx, column)
	if err != nil {
		return ErrUpdateColumn
	}
	return nil
}

func (uc *AggregatorUseCase) UpdateCard(ctx context.Context, card *dto.Card) error {
	err := uc.todoSvc.UpdateCard(ctx, card)
	if err != nil {
		return ErrUpdateCard
	}
	return nil
}

func (uc *AggregatorUseCase) DeleteBoard(ctx context.Context, id string) error {
	err := uc.todoSvc.DeleteBoard(ctx, id)
	if err != nil {
		return ErrDeleteBoard
	}
	return nil
}

func (uc *AggregatorUseCase) DeleteColumn(ctx context.Context, id string) error {
	err := uc.todoSvc.DeleteColumn(ctx, id)
	if err != nil {
		return ErrDeleteColumn
	}
	return nil
}

func (uc *AggregatorUseCase) DeleteCard(ctx context.Context, id string) error {
	err := uc.todoSvc.DeleteCard(ctx, id)
	if err != nil {
		return ErrDeleteCard
	}
	return nil
}
