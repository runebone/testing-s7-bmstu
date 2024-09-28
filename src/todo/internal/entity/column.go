package entity

import (
	"time"

	"github.com/google/uuid"
)

type Column struct {
	ID        uuid.UUID `db:"id"`
	UserID    uuid.UUID `db:"user_id"`
	BoardID   uuid.UUID `db:"board_id"`
	Title     string    `db:"title"`
	Position  float64   `db:"position"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
