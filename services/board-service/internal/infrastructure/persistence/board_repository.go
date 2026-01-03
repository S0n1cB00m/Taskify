package persistence

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"Taskify/services/board-service/internal/domain/board"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ board.Repository = (*BoardRepository)(nil)

type BoardRepository struct {
	db *pgxpool.Pool
}

func NewBoardRepository(db *pgxpool.Pool) *BoardRepository {
	return &BoardRepository{db: db}
}

func (r *BoardRepository) Create(ctx context.Context, b *board.Board) error {
	query := "INSERT INTO boards(title, description, user_id) VALUES ($1, $2, $3) RETURNING id"

	model := fromDomain(b)

	err := r.db.QueryRow(ctx, query, model.Title, model.Description, model.Owner).Scan(&model.ID)
	if err != nil {
		// Здесь можно залогировать или обернуть ошибку
		return fmt.Errorf("failed to create board: %w", err)
	}

	b.ID = model.ID

	return nil
}

func (r *BoardRepository) GetByID(ctx context.Context, id int64) (*board.Board, error) {
	query := "SELECT id, title, description, user_id, created_at, updated_at FROM boards WHERE id = $1"

	var model BoardModel

	err := r.db.QueryRow(ctx, query, id).Scan(&model.ID, &model.Title, &model.Description, &model.Owner, &model.CreatedAt, &model.UpdatedAt)
	if err != nil {
		// 3. Обрабатываем случай, когда запись не найдена
		if errors.Is(err, pgx.ErrNoRows) {
			// Возвращаем доменную ошибку (она должна быть объявлена в domain/board/errors.go)
			return nil, board.ErrBoardNotFound
		}
		return nil, fmt.Errorf("failed to get board: %w", err)
	}

	return model.toDomain(), nil
}

func (r *BoardRepository) GetList(ctx context.Context) ([]*board.Board, error) {
	query := "SELECT id, title, description, user_id, created_at, updated_at FROM boards"

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query boards: %w", err)
	}
	// Гарантируем закрытие rows (освобождение соединения)
	defer rows.Close()

	// Инициализируем слайс, чтобы он был [] (empty json), а не null, если записей нет
	boardsList := make([]*board.Board, 0)

	for rows.Next() {
		var model BoardModel
		if err := rows.Scan(
			&model.ID,
			&model.Title,
			&model.Description,
			&model.Owner,
			&model.CreatedAt,
			&model.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan board: %w", err)
		}

		// Конвертируем и добавляем в результат
		boardsList = append(boardsList, model.toDomain())
	}

	// Проверяем ошибки, которые могли возникнуть во время итерации (сеть и т.д.)
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return boardsList, nil
}

func (r *BoardRepository) Update(ctx context.Context, b *board.Board) (*board.Board, error) {
	query := "UPDATE boards SET title = $1, description = $2, updated_at = $3 WHERE id = $4 RETURNING id, title, description, user_id, created_at, updated_at"

	// Конвертируем string -> sql.NullString для description
	desc := sql.NullString{String: b.Description, Valid: b.Description != ""}

	var model BoardModel

	err := r.db.QueryRow(ctx, query, b.Title, desc, b.UpdatedAt, b.ID).Scan(
		&model.ID,
		&model.Title,       // <--- добавил &
		&model.Description, // <--- добавил &
		&model.Owner,       // <--- добавил &
		&model.CreatedAt,
		&model.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update board: %w", err)
	}

	return model.toDomain(), nil
}

func (r *BoardRepository) Delete(ctx context.Context, id int64) error {
	query := "DELETE FROM boards WHERE id = $1"

	// Выполняем запрос и получаем tag команды (информацию о результате)
	commandTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete board: %w", err)
	}

	// Проверяем, сколько строк было удалено
	if commandTag.RowsAffected() == 0 {
		// Если 0, значит доски с таким ID не было
		return board.ErrBoardNotFound
	}

	return nil
}
