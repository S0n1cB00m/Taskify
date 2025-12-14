package columns

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, column *Column) error
	GetByID(ctx context.Context, id int64) (*Column, error)
	Update(ctx context.Context, column *Column) error
	Delete(ctx context.Context, id int64) error
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, column *Column) error {
	return nil
}

func (r *repository) GetByID(ctx context.Context, id int64) (*Column, error) {
	return nil, nil
}

func (r *repository) Update(ctx context.Context, column *Column) error {
	return nil
}

func (r *repository) Delete(ctx context.Context, id int64) error {
	return nil
}
