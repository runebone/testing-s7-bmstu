package repository

import (
	"context"
	"time"
	"todo/internal/entity"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SQLXCardRepository struct {
	db *sqlx.DB
}

func NewSQLXCardRepository(db *sqlx.DB) *SQLXCardRepository {
	return &SQLXCardRepository{db: db}
}

func (r *SQLXCardRepository) CreateCard(ctx context.Context, card *entity.Card) error {
	query := `
	INSERT INTO cards (id, column_id, user_id, title, description, position, created_at, updated_at)
	VALUES (:id, :column_id, :user_id, :title, :description, :position, :created_at, :updated_at)
	`

	_, err := r.db.NamedExecContext(ctx, query, card)

	return err
}

func (r *SQLXCardRepository) GetCardByID(ctx context.Context, id uuid.UUID) (*entity.Card, error) {
	query := `
	SELECT FROM cards WHERE id = $1
	`

	var card entity.Card
	err := r.db.GetContext(ctx, &card, query, id)

	if err != nil {
		return nil, err
	}

	return &card, nil
}

func (r *SQLXCardRepository) GetCardsByColumn(ctx context.Context, columnID uuid.UUID, limit, offset int) ([]entity.Card, error) {
	query := `
	SELECT FROM cards WHERE column_id = $1
	ORDER BY created_at ASC
	LIMIT $2
	OFFSET $3
	`

	var cards []entity.Card
	err := r.db.SelectContext(ctx, &cards, query, columnID, limit, offset)

	if err != nil {
		return nil, err
	}

	return cards, nil
}

func (r *SQLXCardRepository) UpdateCard(ctx context.Context, card *entity.Card) error {
	query := `
    UPDATE cards SET
	column_id = :column_id,
	title = :title,
	description = :description,
	position = :position,
	updated_at = :updated_at
    WHERE id = :id
    `

	_, err := r.db.NamedExecContext(ctx, query, card)

	return err
}

func (r *SQLXCardRepository) DeleteCard(ctx context.Context, id uuid.UUID) error {
	query := `
	DELETE FROM cards WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id)

	return err
}

func (r *SQLXCardRepository) GetNewCards(ctx context.Context, from, to time.Time) ([]entity.Card, error) {
	query := `
	SELECT FROM cards
	WHERE $1 <= created_at AND created_at <= $2
	`

	var cards []entity.Card
	err := r.db.SelectContext(ctx, &cards, query, from, to)

	if err != nil {
		return nil, err
	}

	return cards, nil
}
