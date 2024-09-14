package repositories

import (
	"context"
	"user/internal/entities"
	"user/internal/middleware/cache"
	r "user/internal/repositories"

	"github.com/google/uuid"
)

type CachedUserRepository struct {
	baseRepo r.UserRepository
	cache    *cache.CacheMiddleware
}

func NewCachedUserRepository(baseRepo r.UserRepository, cache *cache.CacheMiddleware) *CachedUserRepository {
	return &CachedUserRepository{
		baseRepo: baseRepo,
		cache:    cache,
	}
}

func (r *CachedUserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	cacheKey := id.String()

	cachedData, err := r.cache.GetOrSet(ctx, cacheKey, func() (interface{}, error) {
		return r.baseRepo.GetUserByID(ctx, id)
	})

	if err != nil {
		return nil, err
	}

	user := cachedData.(*entities.User)
	return user, nil
}

func (r *CachedUserRepository) CreateUser(ctx context.Context, user *entities.User) error {
	return r.baseRepo.CreateUser(ctx, user)
}

func (r *CachedUserRepository) UpdateUser(ctx context.Context, user *entities.User) error {
	return r.baseRepo.UpdateUser(ctx, user)
}

func (r *CachedUserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return r.baseRepo.DeleteUser(ctx, id)
}
