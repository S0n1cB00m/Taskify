package columns

import "context"

type UseCase interface {
	Create(ctx context.Context, column *Column) error
	Delete(ctx context.Context, column *Column) error
}

type useCase struct {
	repo Repository
}

func NewUseCase(repo Repository) UseCase {
	return &useCase{repo: repo}
}

func (uc *useCase) Create(ctx context.Context, column *Column) error {
	err := uc.repo.Create(ctx, column)
	return err
}

func (uc *useCase) Delete(ctx context.Context, column *Column) error {
	err := uc.repo.Delete(ctx, column)
	return err
}
