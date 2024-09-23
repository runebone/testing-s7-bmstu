package dto

import "github.com/google/uuid"

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
