package dto

import (
	"todo/internal/entity"

	"github.com/google/uuid"
)

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

func ToColumnDTOs(columns []entity.Column) []Column {
	columnDTOs := make([]Column, len(columns))
	for i, column := range columns {
		columnDTOs[i] = Column{
			ID:       column.ID,
			BoardID:  column.BoardID,
			Title:    column.Title,
			Position: column.Position,
		}
	}
	return columnDTOs
}
