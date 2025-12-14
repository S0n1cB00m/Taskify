package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id int64) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int64) error
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, user *User) error {
	return nil
}

func (r *repository) GetByID(ctx context.Context, id int64) (*User, error) {
	log := zerolog.Ctx(ctx)

	ctx.log.Trace().Int64("user_id", id).Msg("attempting to get user")

	var user User

	user.ID = id

	row := r.db.QueryRow(ctx, "SELECT email, username, password FROM users WHERE id = $1", id)

	err := row.Scan(&user.Email, &user.Username, &user.Password)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.log.Debug().Int64("user_id", id).Msg("user not found")
			return nil, ErrUserNotFound
		}

		r.log.Error().Err(err).Int64("user_id", id).Msg("failed to get user by id")
		return nil, fmt.Errorf("users.repository.GetByID: %w", err)
	}

	r.log.Debug().Int64("user_id", id).Msg("user retrieved successfully")
	return &user, nil
}

func (r *repository) Update(ctx context.Context, user *User) error {
	return nil
}

func (r *repository) Delete(ctx context.Context, id int64) error {
	return nil
}
