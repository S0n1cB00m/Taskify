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
	Create(ctx context.Context, user *User) (*User, error)
	GetByID(ctx context.Context, id int64) (*User, error)
	Update(ctx context.Context, user *User) (*User, error)
	Delete(ctx context.Context, id int64) error
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, user *User) (*User, error) {
	log := zerolog.Ctx(ctx)

	log.Trace().Str("email", user.Email).Msg("attempting to create user")

	query := "INSERT INTO users (email, username, password) VALUES ($1, $2, $3) RETURNING id"

	err := r.db.QueryRow(ctx, query, user.Email, user.Username, user.Password).Scan(&user.ID)
	if err != nil {
		log.Error().Err(err).Str("email", user.Email).Msg("failed to create user")
		return nil, fmt.Errorf("users.repository.Create: %w", err)
	}

	log.Debug().Int64("user_id", user.ID).Msg("user created successfully")

	return user, nil
}

func (r *repository) GetByID(ctx context.Context, id int64) (*User, error) {
	log := zerolog.Ctx(ctx)

	log.Trace().Int64("user_id", id).Msg("attempting to get user")

	var user User

	user.ID = id

	query := "SELECT email, username FROM users WHERE id = $1"

	row := r.db.QueryRow(ctx, query, id)

	err := row.Scan(&user.Email, &user.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Debug().Int64("user_id", id).Msg("user not found")
			return nil, ErrUserNotFound
		}

		log.Error().Err(err).Int64("user_id", id).Msg("failed to get user by id")
		return nil, fmt.Errorf("users.repository.GetByID: %w", err)
	}

	log.Debug().Int64("user_id", id).Msg("user retrieved successfully")
	return &user, nil
}

func (r *repository) Update(ctx context.Context, user *User) (*User, error) {
	log := zerolog.Ctx(ctx)

	log.Trace().Int64("user_id", user.ID).Str("email", user.Email).Msg("attempting to update user")

	query := "UPDATE users SET email = $1, username = $2, password = $3 WHERE id = $4"

	_, err := r.db.Exec(ctx, query, user.Email, user.Username, user.Password, user.ID)

	if err != nil {
		log.Error().Err(err).Str("email", user.Email).Msg("failed to update user")
		return nil, fmt.Errorf("users.repository.Update: %w", err)
	}

	log.Debug().Int64("user_id", user.ID).Msg("user updated successfully")

	return user, nil
}

func (r *repository) Delete(ctx context.Context, id int64) error {
	log := zerolog.Ctx(ctx)

	log.Trace().Int64("user_id", id).Msg("attempting to delete user")

	query := "DELETE FROM users WHERE id = $1"

	commandTag, err := r.db.Exec(ctx, query, id)

	if err != nil {
		log.Error().Err(err).Int64("user_id", id).Msg("failed to execute delete query")
		return fmt.Errorf("users.repository.Delete: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		log.Debug().Int64("user_id", id).Msg("user not found for deletion")
		return ErrUserNotFound
	}

	log.Debug().Int64("user_id", id).Msg("user deleted successfully")

	return nil
}
