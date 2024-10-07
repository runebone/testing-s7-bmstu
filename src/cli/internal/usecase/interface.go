package client

import (
	"cli/internal/dto"
	"context"
)

type Client interface {
	Register(ctx context.Context, username, email, password string)
	Login(ctx context.Context, email, password string)
	Refresh(ctx context.Context, refreshToken string)
	Validate(ctx context.Context, token string)
	Logout(ctx context.Context, refreshToken string)

	ShowBoards(ctx context.Context)
	ShowBoard(ctx context.Context, boardID string)
	ShowColumn(ctx context.Context, columnID string)
	ShowCard(ctx context.Context, cardID string)

	CreateBoard(ctx context.Context, board dto.Board)
	CreateColumn(ctx context.Context, column dto.Column)
	CreateCard(ctx context.Context, card dto.Card)

	UpdateBoard(ctx context.Context, board *dto.Board)
	UpdateColumn(ctx context.Context, column *dto.Column)
	UpdateCard(ctx context.Context, card *dto.Card)

	DeleteBoard(ctx context.Context, id string)
	DeleteColumn(ctx context.Context, id string)
	DeleteCard(ctx context.Context, id string)

	Stats(ctx context.Context, from, to string)
}
