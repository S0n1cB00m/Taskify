package users

import "context"

type UseCase interface {
	Create(ctx context.Context, user *User) (*User, error)
	GetByID(ctx context.Context, id int64) (*User, error)
	Update(ctx context.Context, user *User) (*User, error)
	Delete(ctx context.Context, id int64) error
}

type useCase struct {
	repo Repository
}

func NewUseCase(repo Repository) UseCase {
	return &useCase{repo: repo}
}

func (uc *useCase) Create(ctx context.Context, user *User) (*User, error) {
	userDB, err := uc.repo.Create(ctx, user)
	return userDB, err
}

func (uc *useCase) GetByID(ctx context.Context, id int64) (*User, error) {
	user, err := uc.repo.GetByID(ctx, id)
	return user, err
}

func (uc *useCase) Update(ctx context.Context, user *User) (*User, error) {
	userDB, err := uc.repo.Update(ctx, user)
	return userDB, err
}

func (uc *useCase) Delete(ctx context.Context, id int64) error {
	err := uc.repo.Delete(ctx, id)
	return err
}
