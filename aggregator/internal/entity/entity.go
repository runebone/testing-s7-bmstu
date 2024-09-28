package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID
	Username string
	Email    string
}

type Card struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Title       string
	Description string
}

type NewUsersAndCardsStats struct {
	Date               time.Time
	Users              []User
	Cards              []Card
	NumCardsByNewUsers int
}
