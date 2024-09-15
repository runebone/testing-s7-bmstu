package repository

import (
	"context"
	"time"
	"user/internal/common/logger"
	"user/internal/entity"
	r "user/internal/repository"

	"github.com/google/uuid"
)

type LoggingUserRepository struct {
	repo   r.UserRepository
	logger logger.Logger
}

func NewLoggingUserRepository(repo r.UserRepository, logger logger.Logger) r.UserRepository {
	return &LoggingUserRepository{
		repo:   repo,
		logger: logger,
	}
}

func (l *LoggingUserRepository) CreateUser(ctx context.Context, user *entity.User) error {
	start := time.Now()

	l.logger.WithFields(map[string]interface{}{
		"action":  "CreateUser",
		"user_id": user.ID,
	}).Info(ctx, "Starting to create user")

	err := l.repo.CreateUser(ctx, user)
	if err != nil {
		l.logger.WithFields(map[string]interface{}{
			"action":  "CreateUser",
			"user_id": user.ID,
			"error":   err.Error(),
		}).Error(ctx, "Failed to create user")
		return err
	}

	l.logger.WithFields(map[string]interface{}{
		"action":   "CreateUser",
		"user_id":  user.ID,
		"duration": time.Since(start),
	}).Info(ctx, "User created successfully")
	return nil
}

func (l *LoggingUserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	start := time.Now()

	l.logger.WithFields(map[string]interface{}{
		"action":  "GetUserByID",
		"user_id": id,
	}).Info(ctx, "Fetching user by ID")

	user, err := l.repo.GetUserByID(ctx, id)
	if err != nil {
		l.logger.WithFields(map[string]interface{}{
			"action":  "GetUserByID",
			"user_id": id,
			"error":   err.Error(),
		}).Error(ctx, "Failed to fetch user by ID")
		return nil, err
	}

	l.logger.WithFields(map[string]interface{}{
		"action":   "GetUserByID",
		"user_id":  id,
		"duration": time.Since(start),
	}).Info(ctx, "User fetched successfully")
	return user, nil
}

func (l *LoggingUserRepository) GetUsersBatch(ctx context.Context, limit, offset int) ([]entity.User, error) {
	start := time.Now()

	l.logger.WithFields(map[string]interface{}{
		"action": "GetUsersBatch",
		"limit":  limit,
		"offset": offset,
	}).Info(ctx, "Fetching users batch")

	users, err := l.repo.GetUsersBatch(ctx, limit, offset)
	if err != nil {
		l.logger.WithFields(map[string]interface{}{
			"action": "GetUsersBatch",
			"limit":  limit,
			"offset": offset,
			"error":  err.Error(),
		}).Error(ctx, "Failed to fetch users batch")
		return nil, err
	}

	l.logger.WithFields(map[string]interface{}{
		"action":   "GetUserBatch",
		"limit":    limit,
		"offset":   offset,
		"duration": time.Since(start),
	}).Info(ctx, "Users batch fetched successfully")
	return users, nil
}

func (l *LoggingUserRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	start := time.Now()

	l.logger.WithFields(map[string]interface{}{
		"action":  "UpdateUser",
		"user_id": user.ID,
	}).Info(ctx, "Updating user")

	err := l.repo.UpdateUser(ctx, user)
	if err != nil {
		l.logger.WithFields(map[string]interface{}{
			"action":  "UpdateUser",
			"user_id": user.ID,
			"error":   err.Error(),
		}).Error(ctx, "Failed to update user")
		return err
	}

	l.logger.WithFields(map[string]interface{}{
		"action":   "UpdateUser",
		"user_id":  user.ID,
		"duration": time.Since(start),
	}).Info(ctx, "User updated successfully")
	return nil
}

func (l *LoggingUserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	start := time.Now()

	l.logger.WithFields(map[string]interface{}{
		"action":  "DeleteUser",
		"user_id": id,
	}).Info(ctx, "Deleting user")

	err := l.repo.DeleteUser(ctx, id)
	if err != nil {
		l.logger.WithFields(map[string]interface{}{
			"action":  "DeleteUser",
			"user_id": id,
			"error":   err.Error(),
		}).Error(ctx, "Failed to delete user")
		return err
	}

	l.logger.WithFields(map[string]interface{}{
		"action":   "DeleteUser",
		"user_id":  id,
		"duration": time.Since(start),
	}).Info(ctx, "User deleted successfully")
	return nil
}
