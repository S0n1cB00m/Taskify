package boards

import (
	"context"
)

type UseCase interface {
	Create(ctx context.Context, board *Board) (*Board, error)
	GetByID(ctx context.Context, board *Board) (*Board, error)
	Update(ctx context.Context, board *Board) (*Board, error)
	Delete(ctx context.Context, board *Board) error
}

type useCase struct {
	repo Repository
}

func NewUseCase(repo Repository) UseCase {
	return &useCase{repo: repo}
}

func (uc *useCase) Create(ctx context.Context, board *Board) (*Board, error) {
	createdBoard, err := uc.repo.Create(ctx, board)
	if err != nil {
		return nil, err
	}

	return createdBoard, err
}

func (uc *useCase) GetByID(ctx context.Context, board *Board) (*Board, error) {
	receivedBoard, err := uc.repo.GetByID(ctx, board)
	return receivedBoard, err
}

func (uc *useCase) Update(ctx context.Context, board *Board) (*Board, error) {
	updatedBoard, err := uc.repo.Update(ctx, board)
	if err != nil {
		return nil, err
	}

	return updatedBoard, err
}

func (uc *useCase) Delete(ctx context.Context, board *Board) error {
	err := uc.repo.Delete(ctx, board)
	return err
}
