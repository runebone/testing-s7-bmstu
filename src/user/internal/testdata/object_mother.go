package testdata

import (
	"time"
	"user/internal/entity"

	"github.com/google/uuid"
)

type UserObjectMother struct{}

func (u *UserObjectMother) ValidUser() entity.User {
	return entity.User{
		ID:           uuid.New(),
		Username:     "ValidUser",
		Email:        "valid@example.com",
		PasswordHash: "Password@123",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

func (u *UserObjectMother) InvalidEmailUser() entity.User {
	return entity.User{
		ID:           uuid.New(),
		Username:     "InValidUser",
		Email:        "invalid-email",
		PasswordHash: "Password@123",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}
