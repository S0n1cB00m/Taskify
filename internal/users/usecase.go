package users

import "context"

type UseCase interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id int64) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int64) error
}

type useCase struct {
	repo Repository
}

func NewUseCase(repo Repository) UseCase {
	return &useCase{repo: repo}
}

func (uc *useCase) Create(ctx context.Context, user *User) error {
	return nil
}

func (uc *useCase) GetByID(ctx context.Context, id int64) (*User, error) {
	user, err := uc.repo.GetByID(ctx, id)
	return user, err
}

func (uc *useCase) Update(ctx context.Context, user *User) error {
	return nil
}

func (uc *useCase) Delete(ctx context.Context, id int64) error {
	return nil
}
