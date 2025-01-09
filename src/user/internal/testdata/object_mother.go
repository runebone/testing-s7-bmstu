package testdata

import (
	"time"
	"user/internal/entity"
	"user/internal/repository"

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

func (u *UserObjectMother) PositiveUsernameFilter() repository.UserFilter {
	username := "PositiveUser"
	return repository.UserFilter{
		Username: &username,
	}
}

func (u *UserObjectMother) NegativeUsernameFilter() repository.UserFilter {
	username := "NegativeUser"
	return repository.UserFilter{
		Username: &username,
	}
}
