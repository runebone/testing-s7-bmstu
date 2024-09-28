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

func UserToEntity(userDTO *User) entity.User {
	return entity.User{
		ID:       userDTO.ID,
		Username: userDTO.Username,
		Email:    userDTO.Email,
	}
}

func ToUserEntities(userDTOs []User) []entity.User {
	users := make([]entity.User, len(userDTOs))
	for i, userDTO := range userDTOs {
		users[i] = UserToEntity(&userDTO)
	}
	return users
}

type Card struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

func CardToEntity(cardDTO *Card) entity.Card {
	return entity.Card{
		ID:          cardDTO.ID,
		UserID:      cardDTO.UserID,
		Title:       cardDTO.Title,
		Description: cardDTO.Description,
	}
}

func ToCardEntities(cardDTOs []Card) []entity.Card {
	cards := make([]entity.Card, len(cardDTOs))
	for i, cardDTO := range cardDTOs {
		cards[i] = CardToEntity(&cardDTO)
	}
	return cards
}
