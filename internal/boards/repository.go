package boards

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type Repository interface {
	Create(ctx context.Context, board *Board) (*Board, error)
	GetByID(ctx context.Context, board *Board) (*Board, error)
	Update(ctx context.Context, board *Board) (*Board, error)
	Delete(ctx context.Context, board *Board) error
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

//CREATE TABLE boards (
//    id      BIGSERIAL PRIMARY KEY,
//    user_id BIGINT      NOT NULL,
//    index     INT         NOT NULL,
//    name    TEXT        NOT NULL,
//    -- ...
//    UNIQUE (user_id, index)
//);

func (r *repository) Create(ctx context.Context, board *Board) (*Board, error) {
	log := zerolog.Ctx(ctx)

	log.Trace().Str("board", board.Name).Msg("attempting to create board")

	// noinspection SqlNoDataSourceInspection
	query := "INSERT INTO boards (user_id, index, name) VALUES ($1, COALESCE((SELECT MAX(index) + 1 FROM boards WHERE user_id = $1), 1), $2) RETURNING index;"

	err := r.db.QueryRow(ctx, query, board.Name, board.Description, board.UserId).Scan(&board.Index)
	if err != nil {
		log.Error().Err(err).Str("board", board.Name).Msg("failed to create board")
		return nil, fmt.Errorf("boards.repository.Create: %w", err)
	}

	log.Debug().Int64("board_id", board.Index).Msg("board created successfully")

	return board, nil
}

func (r *repository) GetByID(ctx context.Context, board *Board) (*Board, error) {
	log := zerolog.Ctx(ctx)

	log.Trace().Int64("board_id", board.Index).Msg("attempting to get board")

	// noinspection SqlNoDataSourceInspection
	query := "SELECT name, description FROM boards WHERE user_id = $1 AND index = $2"

	row := r.db.QueryRow(ctx, query, board.UserId, board.Index)

	err := row.Scan(&board.Name, &board.Description)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Debug().Int64("board_id", board.Index).Msg("board not found")
			return nil, ErrBoardNotFound
		}

		log.Error().Err(err).Int64("board_id", board.Index).Msg("failed to get board by id")
		return nil, fmt.Errorf("boards.repository.GetByID: %w", err)
	}

	log.Debug().Int64("board_id", board.Index).Msg("board retrieved successfully")
	return board, nil
}

func (r *repository) Update(ctx context.Context, board *Board) (*Board, error) {
	log := zerolog.Ctx(ctx)

	log.Trace().Int64("board_id", board.Index).Str("name", board.Name).Msg("attempting to update board")

	// noinspection SqlNoDataSourceInspection
	query := "UPDATE boards SET name = $1, description = $2 WHERE user_id = $3 AND index = $4"

	_, err := r.db.Exec(ctx, query, board.Name, board.Description, board.UserId, board.Index)

	if err != nil {
		log.Error().Err(err).Str("name", board.Name).Msg("failed to update board")
		return nil, fmt.Errorf("boards.repository.Update: %w", err)
	}

	log.Debug().Int64("board_id", board.Index).Msg("board updated successfully")

	return board, nil
}

func (r *repository) Delete(ctx context.Context, board *Board) error {
	log := zerolog.Ctx(ctx)

	log.Trace().Int64("board_id", board.Index).Msg("attempting to delete board")

	// noinspection SqlNoDataSourceInspection
	query := "DELETE FROM boards WHERE owner_id = $1 AND index = $2"

	commandTag, err := r.db.Exec(ctx, query, board.UserId, board.Index)

	if err != nil {
		log.Error().Err(err).Int64("board_id", board.Index).Msg("failed to execute delete query")
		return fmt.Errorf("boards.repository.Delete: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		log.Debug().Int64("board_id", board.Index).Msg("board not found for deletion")
		return ErrBoardNotFound
	}

	log.Debug().Int64("board_id", board.Index).Msg("board deleted successfully")

	return nil
}
