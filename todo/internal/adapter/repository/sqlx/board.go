package repository

import (
	"context"
	"todo/internal/entity"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SQLXBoardRepository struct {
	db *sqlx.DB
}

func NewSQLXBoardRepository(db *sqlx.DB) *SQLXBoardRepository {
	return &SQLXBoardRepository{db: db}
}

func (r *SQLXBoardRepository) CreateBoard(ctx context.Context, board *entity.Board) error {
	query := `
    INSERT INTO boards (id, user_id, title, created_at, updated_at)
	VALUES (:id, :user_id, :title, :created_at, :updated_at)
    `

	_, err := r.db.NamedExecContext(ctx, query, board)

	return err
}

func (r *SQLXBoardRepository) GetBoardByID(ctx context.Context, id uuid.UUID) (*entity.Board, error) {
	query := `
	SELECT * FROM boards WHERE id = $1
	`

	var board *entity.Board
	err := r.db.GetContext(ctx, board, query, id)

	if err != nil {
		return nil, err
	}

	return board, err
}

func (r *SQLXBoardRepository) GetBoardsByUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]entity.Board, error) {
	query := `
	SELECT FROM boards WHERE user_id = $1
	ORDER BY created_at ASC
	LIMIT $2
	OFFSET $3
	`

	var boards []entity.Board
	err := r.db.SelectContext(ctx, boards, query, userID, limit, offset)

	if err != nil {
		return nil, err
	}

	return boards, nil
}

func (r *SQLXBoardRepository) UpdateBoard(ctx context.Context, board *entity.Board) error {
	query := `
    UPDATE boards SET
	title = :title
	updated_at = :updated_at
    WHERE id = :id
    `

	_, err := r.db.NamedExecContext(ctx, query, board)

	return err
}

func (r *SQLXBoardRepository) DeleteBoard(ctx context.Context, id uuid.UUID) error {
	query := `
	DELETE FROM boards WHERE id = $1
	`

	_, err := r.db.ExecContext(ctx, query, id)

	return err
}
