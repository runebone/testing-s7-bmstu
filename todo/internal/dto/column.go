package dto

import "github.com/google/uuid"

type CreateColumnRequest struct {
	UserID   uuid.UUID `json:"user_id"`
	BoardID  uuid.UUID `json:"board_id"`
	Title    string    `json:"title"`
	Position float64   `json:"position"`
}

type Column struct {
	ID       uuid.UUID `json:"id"`
	BoardID  uuid.UUID `json:"board_id"`
	Title    string    `json:"title"`
	Position float64   `json:"position"`
}

type UpdateColumnRequest struct {
	ID       uuid.UUID `json:"id"`
	Title    string    `json:"title,omitempty"`
	Position float64   `json:"position,omitempty"`
}
