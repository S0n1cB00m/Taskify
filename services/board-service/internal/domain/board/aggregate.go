package board

import (
	"time"
)

type Board struct {
	ID          int64
	Title       string
	Description string
	Owner       int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewBoard(title, description string, owner int64) (*Board, error) {
	if title == "" {
		return nil, ErrTitleRequired
	}

	var runes = []rune(title)
	if len(runes) > 100 {
		return nil, ErrTitleTooLong
	}

	if owner == 0 {
		return nil, ErrEmptyOwner

	}

	return &Board{
		Title:       title,
		Description: description,
		Owner:       owner,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}
