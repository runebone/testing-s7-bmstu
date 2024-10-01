package http

import (
	"aggregator/internal/dto"
	"aggregator/internal/service/user"
	"context"
	"net/http"
	"time"
)

type UserService struct {
	baseURL    string
	httpClient *http.Client
}

func NewUserService(baseURL string, timeout time.Duration) user.UserService {
	return &UserService{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (s *UserService) GetNewUsers(ctx context.Context, from, to time.Time) ([]dto.User, error) {
	// TODO:
	return nil, nil
}
