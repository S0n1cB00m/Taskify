package board

import "errors"

var (
	ErrBoardNotFound = errors.New("board not found")
	ErrTitleRequired = errors.New("board title is required")
	ErrTitleTooLong  = errors.New("board title is too long")
	ErrEmptyOwner    = errors.New("owner is empty")
)
