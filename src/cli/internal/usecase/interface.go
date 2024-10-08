package usecase

import (
	"cli/internal/dto"
	"context"
)

type Client interface {
	Register(ctx context.Context, username, email, password string) (*dto.Tokens, error)
	Login(ctx context.Context, email, password string) (*dto.Tokens, error)
	Refresh(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error)
	Validate(ctx context.Context, token string) (*dto.ValidateTokenResponse, error)
	Logout(ctx context.Context, refreshToken string) error

	// context with value tokens
	ShowBoards(ctx context.Context)
	ShowBoard(ctx context.Context, boardID string)
	ShowColumn(ctx context.Context, columnID string)
	ShowCard(ctx context.Context, cardID string)

	CreateBoard(ctx context.Context, title string)                       // board dto.Board)
	CreateColumn(ctx context.Context, boardID, title string)             // column dto.Column)
	CreateCard(ctx context.Context, columnID, title, description string) // card dto.Card)

	UpdateBoard(ctx context.Context, boardID, title string)                                   // board *dto.Board)
	UpdateColumn(ctx context.Context, boardID, columnID, title string)                        // column *dto.Column)
	UpdateCardTitle(ctx context.Context, boardID, columnID, cardID, title string)             // card *dto.Card)
	UpdateCardDescription(ctx context.Context, boardID, columnID, cardID, description string) // card *dto.Card)

	DeleteBoard(ctx context.Context, id string)
	DeleteColumn(ctx context.Context, id string)
	DeleteCard(ctx context.Context, id string)

	Stats(ctx context.Context, from, to string)
}
