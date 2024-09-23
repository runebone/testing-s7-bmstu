package dto

import (
	"todo/internal/entity"

	"github.com/google/uuid"
)

type CreateBoardRequest struct {
	UserID uuid.UUID `json:"user_id"`
	Title  string    `json:"title"`
}

type Board struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
}

type UpdateBoardRequest struct {
	Board
}

func ToBoardDTOs(boards []entity.Board) []Board {
	boardDTOs := make([]Board, len(boards))
	for i, board := range boards {
		boardDTOs[i] = Board{
			ID:    board.ID,
			Title: board.Title,
		}
	}
	return boardDTOs
}
