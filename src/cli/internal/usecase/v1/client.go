package v1

import (
	"cli/internal/dto"
	"cli/internal/service"
	"cli/internal/usecase"
	"context"
	"fmt"
)

type ClientUseCase struct {
	svc service.AggregatorService
}

func NewClientUseCase(svc service.AggregatorService) usecase.Client {
	return &ClientUseCase{
		svc: svc,
	}
}

func (uc *ClientUseCase) Register(ctx context.Context, username, email, password string) (*dto.Tokens, error) {
	tokens, err := uc.svc.Register(ctx, username, email, password)
	return tokens, err
}

func (uc *ClientUseCase) Login(ctx context.Context, email, password string) (*dto.Tokens, error) {
	tokens, err := uc.svc.Login(ctx, email, password)
	return tokens, err
}

func (uc *ClientUseCase) Refresh(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error) {
	resp, err := uc.svc.Refresh(ctx, refreshToken)
	return resp, err
}

func (uc *ClientUseCase) Validate(ctx context.Context, token string) (*dto.ValidateTokenResponse, error) {
	resp, err := uc.svc.Validate(ctx, token)
	return resp, err
}

func (uc *ClientUseCase) Logout(ctx context.Context, refreshToken string) error {
	err := uc.svc.Logout(ctx, refreshToken)
	return err
}

func (uc *ClientUseCase) ShowBoards(ctx context.Context) {
	tokens, ok := ctx.Value("tokens").(dto.Tokens)
	if !ok {
		fmt.Println("failed to get tokens from context")
	}

	resp, err := uc.svc.Validate(tokens.AccessToken)
	if err != nil {
		fmt.Println("access token expired")
	}

	userID := resp.UserID
	role := resp.Role

	boards, err := uc.svc.ShowBoards(ctx)
}

func (uc *ClientUseCase) ShowBoard(ctx context.Context, boardID string) {
	columns, err := uc.svc.ShowBoard(ctx, boardID)
}

func (uc *ClientUseCase) ShowColumn(ctx context.Context, columnID string) {
	cards, err := uc.svc.ShowColumn(ctx, columnID)
}

func (uc *ClientUseCase) ShowCard(ctx context.Context, cardID string) {
	card, err := uc.svc.ShowCard(ctx, cardID)
}

func (uc *ClientUseCase) CreateBoard(ctx context.Context, title string) {
	err := uc.svc.CreateBoard(ctx, board)
}

func (uc *ClientUseCase) CreateColumn(ctx context.Context, boardID, title string) {
	err := uc.svc.CreateColumn(ctx, column)
}

func (uc *ClientUseCase) CreateCard(ctx context.Context, columnID, title, description string) {
	err := uc.svc.CreateCard(ctx, card)
}

func (uc *ClientUseCase) UpdateBoard(ctx context.Context, boardID, title string) {
	err := uc.svc.UpdateBoard(ctx, board)
}

func (uc *ClientUseCase) UpdateColumn(ctx context.Context, boardID, columnID, title string) {
	err := uc.svc.UpdateColumn(ctx, column)
}

func (uc *ClientUseCase) UpdateCardTitle(ctx context.Context, boardID, columnID, cardID, title string) {
	err := uc.svc.UpdateCard(ctx, card)
}

func (uc *ClientUseCase) UpdateCardDescription(ctx context.Context, boardID, columnID, cardID, description string) {
	err := uc.svc.UpdateCard(ctx, card)
}

func (uc *ClientUseCase) DeleteBoard(ctx context.Context, id string) {
	err := uc.svc.DeleteBoard(ctx, id)
}

func (uc *ClientUseCase) DeleteColumn(ctx context.Context, id string) {
	err := uc.svc.DeleteColumn(ctx, id)
}

func (uc *ClientUseCase) DeleteCard(ctx context.Context, id string) {
	err := uc.svc.DeleteCard(ctx, id)
}

func (uc *ClientUseCase) Stats(ctx context.Context, from, to string) {
	stats, err := uc.svc.Stats(ctx, from, to)
}
