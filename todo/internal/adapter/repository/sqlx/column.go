package repository

import (
	"context"
	"todo/internal/entity"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SQLXColumnRepository struct {
	db *sqlx.DB
}

func NewSQLXColumnRepository(db *sqlx.DB) *SQLXColumnRepository {
	return &SQLXColumnRepository{db: db}
}

func (r *SQLXColumnRepository) CreateColumn(ctx context.Context, column *entity.Column) error {
	query := `
	INSERT INTO columns (id, board_id, user_id, title, position, created_at, updated_at)
	VALUES (:id, :board_id, :user_id, :title, :position, :created_at, :updated_at)
	`

	_, err := r.db.NamedExecContext(ctx, query, column)

	return err
}

func (r *SQLXColumnRepository) GetColumnByID(ctx context.Context, id uuid.UUID) (*entity.Column, error) {
	query := `
	SELECT FROM columns WHERE id = $1
	`

	var column *entity.Column
	err := r.db.GetContext(ctx, column, query, id)

	if err != nil {
		return nil, err
	}

	return column, nil
}

func (r *SQLXColumnRepository) GetColumnsByBoard(ctx context.Context, boardID uuid.UUID, limit, offset int) ([]entity.Column, error) {
	query := `
	SELECT FROM columns WHERE board_id = $1
	ORDER BY created_at ASC
	LIMIT $2
	OFFSET $3
	`

	var columns []entity.Column
	err := r.db.SelectContext(ctx, columns, query, boardID, limit, offset)

	if err != nil {
		return nil, err
	}

	return columns, nil
}

func (r *SQLXColumnRepository) UpdateColumn(ctx context.Context, column *entity.Column) error {
	query := `
    UPDATE columns SET
	title = :title
	position = :position
	updated_at = :updated_at
    WHERE id = :id
    `

	_, err := r.db.NamedExecContext(ctx, query, column)

	return err
}

func (r *SQLXColumnRepository) DeleteColumn(ctx context.Context, id uuid.UUID) error {
	query := `
	DELETE FROM columns WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id)

	return err
}
