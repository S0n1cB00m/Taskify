package board

import (
	"context"

	"Taskify/services/board-service/internal/domain/board"
)

type DeleteBoardUseCase struct {
	repo board.Repository
}

func NewDeleteBoardUseCase(repo board.Repository) *DeleteBoardUseCase {
	return &DeleteBoardUseCase{repo: repo}
}

func (uc *DeleteBoardUseCase) Handle(ctx context.Context, id int64) error {
	err := uc.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
