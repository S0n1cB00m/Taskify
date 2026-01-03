package persistence

import (
	"database/sql"
	"time"

	"Taskify/services/board-service/internal/domain/board"
)

type BoardModel struct {
	ID          int64          `json:"id" db:"id" example:"1"`
	Title       string         `json:"title" db:"title" example:"Daily routine"`
	Description sql.NullString `json:"description" db:"description" example:"Board description"`
	Owner       int64          `json:"owner" db:"user_id" example:"1"`
	CreatedAt   time.Time      `json:"createdAt" db:"created_at" example:"2019-09-07 17:40:58"`
	UpdatedAt   time.Time      `json:"updatedAt" db:"updated_at" example:"2019-09-07 17:40:58"`
}

func (m *BoardModel) toDomain() *board.Board {
	return &board.Board{
		ID:          m.ID,
		Title:       m.Title,
		Description: m.Description.String,
		Owner:       m.Owner,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func fromDomain(b *board.Board) *BoardModel {
	return &BoardModel{
		ID:    b.ID,
		Title: b.Title,
		Description: sql.NullString{
			String: b.Description,
			Valid:  b.Description != "",
		},
		Owner:     b.Owner,
		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}
}
