package sqlx

import (
	"auth/internal/entity"
	"context"

	"github.com/jmoiron/sqlx"
)

type SQLXTokenRepository struct {
	db *sqlx.DB
}

func NewSQLXTokenRepository(db *sqlx.DB) *SQLXTokenRepository {
	return &SQLXTokenRepository{
		db: db,
	}
}

func (r *SQLXTokenRepository) Save(ctx context.Context, token *entity.Token) error {
	query := `
        INSERT INTO tokens (id, user_id, token, created_at)
        VALUES ($1, $2, $3, $4)
    `
	_, err := r.db.ExecContext(ctx, query, token.ID, token.UserID, token.Token, token.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *SQLXTokenRepository) Delete(ctx context.Context, tokenID string) error {
	query := `DELETE FROM tokens WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, tokenID)
	if err != nil {
		return err
	}
	return nil
}

func (r *SQLXTokenRepository) FindByToken(ctx context.Context, tokenValue string) (*entity.Token, error) {
	query := `SELECT id, user_id, token, created_at FROM tokens WHERE token = $1`

	var token entity.Token
	err := r.db.GetContext(ctx, &token, query, tokenValue)
	if err != nil {
		return nil, err
	}

	return &token, nil
}
