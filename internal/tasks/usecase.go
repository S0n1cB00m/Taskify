package tasks

import "context"

type UseCase interface {
	Create(ctx context.Context, task *Task) error
	GetByID(ctx context.Context, id int64) (*Task, error)
	Update(ctx context.Context, task *Task) error
	Delete(ctx context.Context, id int64) error
}

type useCase struct {
	repo Repository
}

func NewUseCase(repo Repository) UseCase {
	return &useCase{repo: repo}
}

func (uc *useCase) Create(ctx context.Context, task *Task) error {
	return nil
}

func (uc *useCase) GetByID(ctx context.Context, id int64) (*Task, error) {
	return nil, nil
}

func (uc *useCase) Update(ctx context.Context, task *Task) error {
	return nil
}

func (uc *useCase) Delete(ctx context.Context, id int64) error {
	return nil
}
