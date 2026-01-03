package board

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"

	"Taskify/services/board-service/internal/domain/board"
	domain "Taskify/services/board-service/internal/domain/board"
)

type UpdateBoardUseCase struct {
	repo board.Repository
}

func NewUpdateBoardUseCase(repo board.Repository) *UpdateBoardUseCase {
	return &UpdateBoardUseCase{repo: repo}
}

func (uc *UpdateBoardUseCase) Handle(ctx context.Context, cmd UpdateBoardCommand) (*board.Board, error) {
	// 1. Сначала получаем текущую доску, чтобы убедиться, что она существует
	currentBoard, err := uc.repo.GetByID(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}

	log.Debug().Msgf("current board: %v", *currentBoard)

	// 2. Применяем изменения к доменной сущности (в памяти)
	if cmd.Title != nil {
		if *cmd.Title == "" {
			return nil, domain.ErrTitleRequired
		}
		currentBoard.Title = *cmd.Title
	}

	if cmd.Description != nil {
		currentBoard.Description = *cmd.Description
	}

	// Обновляем время
	currentBoard.UpdatedAt = time.Now()

	log.Debug().Msgf("board data to update: %v", *currentBoard)

	// 3. Сохраняем обновленную сущность
	updatedBoard, err := uc.repo.Update(ctx, currentBoard)
	if err != nil {
		return nil, err
	}

	log.Debug().Msgf("updated board: %v", *updatedBoard)

	return updatedBoard, nil
}
