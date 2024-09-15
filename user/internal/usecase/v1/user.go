package v1

import (
	"context"
	"errors"
	"regexp"
	"user/internal/entity"
	"user/internal/repository"
	"user/internal/usecase"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userUseCase struct {
	repo repository.UserRepository
}

func NewUserUseCase(repo repository.UserRepository) usecase.UserUseCase {
	return &userUseCase{
		repo: repo,
	}
}

func (u *userUseCase) CreateUser(ctx context.Context, user *entity.User) error {
	if !isValidEmail(user.Email) {
		return errors.New("invalid email format")
	}

	if !isValidUsername(user.Username) {
		return errors.New("invalid username format")
	}

	err := validatePassword(user.PasswordHash) // NOTE: Plain password, unencrypted initially
	if err != nil {
		return err
	}

	hashedPassword, err := hashPassword(user.PasswordHash)
	if err != nil {
		return err
	}
	user.PasswordHash = hashedPassword

	user.ID = uuid.New()

	return u.repo.CreateUser(ctx, user)
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`)
	return re.MatchString(email)
}

func isValidUsername(username string) bool {
	re := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]{4,}$`)
	return re.MatchString(username)
}

func validatePassword(password string) error {
	if len(password) < 8 || len(password) > 32 {
		return errors.New("password must be between 8 and 32 characters long")
	}

	hasUpperCase := regexp.MustCompile(`[A-Z]`).MatchString(password)
	if !hasUpperCase {
		return errors.New("password must contain at least one uppercase letter")
	}

	hasLowerCase := regexp.MustCompile(`[a-z]`).MatchString(password)
	if !hasLowerCase {
		return errors.New("password must contain at least one lowercase letter")
	}

	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasDigit {
		return errors.New("password must contain at least one digit")
	}

	hasSpecialChar := regexp.MustCompile(`[!@#~$%^&*(),.?":{}|<>]`).MatchString(password)
	if !hasSpecialChar {
		return errors.New("password must contain at least one special character")
	}

	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (u *userUseCase) GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	return u.repo.GetUserByID(ctx, id)
}

func (u *userUseCase) UpdateUser(ctx context.Context, user *entity.User) error {
	existing_user, _ := u.repo.GetUserByID(ctx, user.ID)

	if existing_user == nil {
		return errors.New("user does not exist")
	}

	return u.repo.UpdateUser(ctx, user)
}

func (u *userUseCase) DeleteUser(ctx context.Context, id uuid.UUID) error {
	existing_user, _ := u.repo.GetUserByID(ctx, id)

	if existing_user == nil {
		return errors.New("user does not exist")
	}

	return u.repo.DeleteUser(ctx, id)
}
