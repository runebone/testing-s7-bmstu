package usecase

import (
	"context"
	"time"
	"user/internal/common/logger"
	"user/internal/entity"
	"user/internal/repository"
	"user/internal/usecase"

	"github.com/google/uuid"
)

// TODO: Add logging WithFields for better logs parsing

type LoggingUserUseCase struct {
	useCase usecase.UserUseCase
	logger  logger.Logger
}

func NewLoggingUserUseCase(u usecase.UserUseCase, l logger.Logger) usecase.UserUseCase {
	return &LoggingUserUseCase{
		useCase: u,
		logger:  l,
	}
}

func (l *LoggingUserUseCase) CreateUser(ctx context.Context, user *entity.User) error {
	start := time.Now()
	l.logger.Info(ctx, "CreateUser called", "user_id", user.ID)

	err := l.useCase.CreateUser(ctx, user)
	if err != nil {
		l.logger.Error(ctx, "CreateUser failed", "error", err)
		return err
	}

	l.logger.Info(ctx, "CreateUser succeeded", "duration", time.Since(start))
	return nil
}

func (l *LoggingUserUseCase) GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	start := time.Now()
	l.logger.Info(ctx, "GetUserByID called", "user_id", id)

	user, err := l.useCase.GetUserByID(ctx, id)
	if err != nil {
		l.logger.Error(ctx, "GetUserByID failed", "error", err)
		return nil, err
	}

	l.logger.Info(ctx, "GetUserByID succeeded", "duration", time.Since(start))
	return user, nil
}

func (l *LoggingUserUseCase) GetUsers(ctx context.Context, filter repository.UserFilter) ([]entity.User, error) {
	start := time.Now()
	l.logger.Info(ctx, "GetUsers called", "filter", filter)

	users, err := l.useCase.GetUsers(ctx, filter)
	if err != nil {
		l.logger.Error(ctx, "GetUsers failed", "error", err)
		return nil, err
	}

	l.logger.Info(ctx, "GetUsers succeeded", "duration", time.Since(start))
	return users, nil
}

func (l *LoggingUserUseCase) GetUsersBatch(ctx context.Context, limit, offset int) ([]entity.User, error) {
	start := time.Now()
	l.logger.Info(ctx, "GetUsersBatch called", "limit", limit, "offset", offset)

	users, err := l.useCase.GetUsersBatch(ctx, limit, offset)
	if err != nil {
		l.logger.Error(ctx, "GetUsersBatch failed", "error", err)
		return nil, err
	}

	l.logger.Info(ctx, "GetUsersBatch succeeded", "duration", time.Since(start))
	return users, nil
}

func (l *LoggingUserUseCase) GetNewUsers(ctx context.Context, from time.Time, to time.Time) ([]entity.User, error) {
	start := time.Now()
	l.logger.Info(ctx, "GetUsersNew called", "from", from, "to", to)

	users, err := l.useCase.GetNewUsers(ctx, from, to)
	if err != nil {
		l.logger.Error(ctx, "GetUsersNew failed", "error", err)
		return nil, err
	}

	l.logger.Info(ctx, "GetUsersNew succeeded", "duration", time.Since(start))
	return users, nil
}

func (l *LoggingUserUseCase) UpdateUser(ctx context.Context, user *entity.User) error {
	start := time.Now()
	l.logger.Info(ctx, "UpdateUser called", "user_id", user.ID)

	err := l.useCase.UpdateUser(ctx, user)
	if err != nil {
		l.logger.Error(ctx, "UpdateUser failed", "error", err)
		return err
	}

	l.logger.Info(ctx, "UpdateUser succeeded", "duration", time.Since(start))
	return nil
}

func (l *LoggingUserUseCase) DeleteUser(ctx context.Context, id uuid.UUID) error {
	start := time.Now()
	l.logger.Info(ctx, "DeleteUser called", "user_id", id)

	err := l.useCase.DeleteUser(ctx, id)
	if err != nil {
		l.logger.Error(ctx, "DeleteUser failed", "error", err)
		return err
	}

	l.logger.Info(ctx, "DeleteUser succeeded", "duration", time.Since(start))
	return nil
}
