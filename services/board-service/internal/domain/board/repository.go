package board

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, board *Board) error

	GetByID(ctx context.Context, int int64) (*Board, error)

	GetList(ctx context.Context) ([]*Board, error)

	Update(ctx context.Context, board *Board) error

	Delete(ctx context.Context, int int64) error
}
