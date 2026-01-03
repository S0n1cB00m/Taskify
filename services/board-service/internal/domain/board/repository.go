package board

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, board *Board) error

	GetByID(ctx context.Context, id int64) (*Board, error)

	GetList(ctx context.Context) ([]*Board, error)

	Update(ctx context.Context, board *Board) (*Board, error)

	Delete(ctx context.Context, id int64) error
}
