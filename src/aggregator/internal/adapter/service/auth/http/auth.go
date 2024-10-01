package http

import (
	"aggregator/internal/dto"
	"aggregator/internal/service/auth"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type AuthService struct {
	baseURL    string
	httpClient *http.Client
}

func NewAuthService(baseURL string, timeout time.Duration) auth.AuthService {
	return &AuthService{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (s *AuthService) Register(ctx context.Context, username, email, password string) (*dto.Tokens, error) {
	// TODO:
	url := fmt.Sprintf("%s/register", s.baseURL)

	data := dto.RegisterRequest{
		Username: username,
		Email:    email,
		Password: password,
	}

	resp, err := s.makeRequest(ctx, http.MethodPost, url, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return nil, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*dto.Tokens, error) {
	// TODO:
	url := fmt.Sprintf("%s/login", s.baseURL)

	data := dto.LoginRequest{
		Email:    email,
		Password: password,
	}

	return nil, nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error) {
	// TODO:
	url := fmt.Sprintf("%s/refresh", s.baseURL)

	data := dto.RefreshRequest{
		RefreshToken: refreshToken,
	}

	return nil, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, token string) (string, string, error) {
	// TODO:
	url := fmt.Sprintf("%s/validate", s.baseURL)

	data := dto.ValidateTokenRequest{
		Token: token,
	}

	return "", "", nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	// TODO:
	url := fmt.Sprintf("%s/logout", s.baseURL)

	data := dto.LogoutRequest{
		RefreshToken: refreshToken,
	}

	return nil
}

func (s *AuthService) makeRequest(ctx context.Context, method, url string, data any) (*http.Response, error) {
	jsonBody, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error marshaling user data: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}

	return resp, nil
}
