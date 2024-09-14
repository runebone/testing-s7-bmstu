package dto

import "github.com/google/uuid"

type UserDTO struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}

type CreateUserDTO struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UpdateUserDTO struct {
	Username *string `json:"username,omitempty"`
	Email    *string `json:"email,omitempty"`
}
