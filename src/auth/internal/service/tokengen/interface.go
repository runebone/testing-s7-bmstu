package tokengen

import "context"

type TokenService interface {
	GenerateAccessToken(ctx context.Context, userID string) (string, error)
	GenerateRefreshToken(ctx context.Context, userID string) (string, error)
	ValidateToken(ctx context.Context, token string) (string, error) // Returns userID
}
