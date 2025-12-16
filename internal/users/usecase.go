package users

import (
	"context"
	"fmt"
)

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
	if err := user.HashPassword(user.Password); err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	createdUser, err := uc.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	createdUser.Password = ""

	return createdUser, err
}

func (uc *useCase) GetByID(ctx context.Context, id int64) (*User, error) {
	user, err := uc.repo.GetByID(ctx, id)
	return user, err
}

func (uc *useCase) Update(ctx context.Context, user *User) (*User, error) {
	if err := user.HashPassword(user.Password); err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	updatedUser, err := uc.repo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	updatedUser.Password = ""

	return updatedUser, err
}

func (uc *useCase) Delete(ctx context.Context, id int64) error {
	err := uc.repo.Delete(ctx, id)
	return err
}
