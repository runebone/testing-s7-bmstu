package dto

import "github.com/google/uuid"

type CreateCardRequest struct {
	UserID      uuid.UUID `json:"user_id"`
	ColumnID    uuid.UUID `json:"column_id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Position    float64   `json:"position"`
}

type Card struct {
	ID          uuid.UUID `json:"id"`
	ColumnID    uuid.UUID `json:"column_id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Position    float64   `json:"position"`
}

type UpdateCardRequest struct {
	ID          uuid.UUID `json:"id"`
	ColumnID    uuid.UUID `json:"column_id,omitempty"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	Position    float64   `json:"position,omitempty"`
}
