package dto

import (
	"aggregator/internal/entity"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func UserToEntity(userDTO User) *entity.User {
	return &entity.User{
		ID:       userDTO.ID,
		Username: userDTO.Username,
		Email:    userDTO.Email,
	}
}

type Card struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

func CardToEntity(cardDTO Card) *entity.Card {
	return &entity.Card{
		ID:          cardDTO.ID,
		UserID:      cardDTO.UserID,
		Title:       cardDTO.Title,
		Description: cardDTO.Description,
	}
}
