package v1

import (
	"context"
	"errors"
	"regexp"
	"time"
	"user/internal/entity"
	"user/internal/repository"
	"user/internal/usecase"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidFromTo = errors.New("<<to>> date should be not less than <<from>>")
)

type userUseCase struct {
	repo repository.UserRepository
}

func NewUserUseCase(repo repository.UserRepository) usecase.UserUseCase {
	return &userUseCase{
		repo: repo,
	}
}

func (u *userUseCase) CreateUser(ctx context.Context, user entity.User) error {
	if !isValidEmail(user.Email) {
		return errors.New("invalid email format")
	}

	if !isValidUsername(user.Username) {
		return errors.New("invalid username format")
	}

	// TODO: Maybe move password validation to auth service,
	// and deal only with hashes in user service.
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

	return u.repo.CreateUser(ctx, &user)
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
	user, err := u.repo.GetUserByID(ctx, id)

	if err != nil {
		return nil, errors.New("failed to get user by id")
	}

	return user, nil
}

func (u *userUseCase) GetUsers(ctx context.Context, filter repository.UserFilter) ([]entity.User, error) {
	return u.repo.GetUsers(ctx, filter)
}

func (u *userUseCase) GetUsersBatch(ctx context.Context, limit, offset int) ([]entity.User, error) {
	if limit < 0 || offset < 0 {
		return nil, errors.New("limit and offset can't be negative")
	} else if limit == 0 {
		return nil, errors.New("limit should be greater than zero")
	}

	return u.repo.GetUsersBatch(ctx, limit, offset)
}

func (u *userUseCase) GetNewUsers(ctx context.Context, from time.Time, to time.Time) ([]entity.User, error) {
	if from.Unix() > to.Unix() {
		return nil, ErrInvalidFromTo
	}

	return u.repo.GetNewUsers(ctx, from, to)
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
