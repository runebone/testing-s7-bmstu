package entity

import (
	"time"

	"github.com/google/uuid"
)

type Card struct {
	ID          uuid.UUID `db:"id"`
	UserID      uuid.UUID `db:"user_id"`
	ColumnID    uuid.UUID `db:"column_id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	Position    float64   `db:"position"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
