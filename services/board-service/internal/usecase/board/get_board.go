package board

import (
	"context"

	"Taskify/services/board-service/internal/domain/board"
)

type GetBoardUseCase struct {
	repo board.Repository
}

func NewGetBoardUseCase(repo board.Repository) *GetBoardUseCase {
	return &GetBoardUseCase{repo: repo}
}

func (uc *GetBoardUseCase) Handle(ctx context.Context, id int64) (*board.Board, error) {
	receivedBoard, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return receivedBoard, nil
}
