package board

import (
	"context"

	"Taskify/services/board-service/internal/domain/board"
)

type CreateBoardUseCase struct {
	repo board.Repository
}

func NewCreateBoardUseCase(repo board.Repository) *CreateBoardUseCase {
	return &CreateBoardUseCase{repo: repo}
}

func (uc *CreateBoardUseCase) Handle(ctx context.Context, cmd CreateBoardCommand) (*board.Board, error) {
	// 1. Создаем доменный агрегат (тут сработает валидация: пустой заголовок и т.д.)
	b, err := board.NewBoard(cmd.Title, cmd.Description, cmd.OwnerID)
	if err != nil {
		return nil, err
	}

	// 2. Сохраняем через репозиторий
	if err := uc.repo.Create(ctx, b); err != nil {
		return nil, err
	}

	// 3. Возвращаем созданную доску (у неё уже будет ID, проставленный репозиторием)
	return b, nil
}
