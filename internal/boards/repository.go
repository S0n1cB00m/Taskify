package boards

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, board *Board) error
	GetByID(ctx context.Context, id int64) (*Board, error)
	Update(ctx context.Context, board *Board) error
	Delete(ctx context.Context, id int64) error
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, user *Board) error {
	return nil
}

func (r *repository) GetByID(ctx context.Context, id int64) (*Board, error) {
	return nil, nil
}

func (r *repository) Update(ctx context.Context, user *Board) error {
	return nil
}

func (r *repository) Delete(ctx context.Context, id int64) error {
	return nil
}
