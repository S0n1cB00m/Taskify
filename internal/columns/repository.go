package columns

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type Repository interface {
	Create(ctx context.Context, column *Column) error
	Delete(ctx context.Context, column *Column) error
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, column *Column) error {
	log := zerolog.Ctx(ctx)

	log.Trace().Msg("attempting to create column")

	var columnId int64

	// noinspection SqlNoDataSourceInspection
	query := "INSERT INTO columns (user_id, index, name) VALUES ($1, COALESCE((SELECT MAX(index) + 1 FROM columns WHERE board_id = $1), 1), $2) RETURNING index;"

	err := r.db.QueryRow(ctx, query, column.BoardId).Scan(&column.Index)
	if err != nil {
		log.Error().Err(err).Msg("failed to create board")
		return fmt.Errorf("columns.repository.Create: %w", err)
	}

	log.Debug().Int64("board_id", columnId).Msg("board created successfully")

	return nil
}

func (r *repository) Delete(ctx context.Context, column *Column) error {
	log := zerolog.Ctx(ctx)

	log.Trace().Int64("column_id", column.Index).Msg("attempting to delete column")

	// noinspection SqlNoDataSourceInspection
	query := "DELETE FROM columns WHERE board_id = $1 AND index = $2"

	commandTag, err := r.db.Exec(ctx, query, column.BoardId, column.Index)

	if err != nil {
		log.Error().Err(err).Int64("column_id", column.Index).Msg("failed to execute delete query")
		return fmt.Errorf("columns.repository.Delete: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		log.Debug().Int64("column_id", column.Index).Msg("column not found for deletion")
		return ErrColumnNotFound
	}

	log.Debug().Int64("column_id", column.Index).Msg("column deleted successfully")

	return nil
}
