package board

import (
	"context"

	"Taskify/services/board-service/internal/domain/board"
)

type ListBoardsUseCase struct {
	repo board.Repository
}

func NewListBoardsUseCase(repo board.Repository) *ListBoardsUseCase {
	return &ListBoardsUseCase{repo: repo}
}

func (uc *ListBoardsUseCase) Handle(ctx context.Context) ([]*board.Board, error) {
	receivedBoards, err := uc.repo.GetList(ctx)
	if err != nil {
		return nil, err
	}

	return receivedBoards, nil
}
