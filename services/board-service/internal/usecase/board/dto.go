package board

type CreateBoardCommand struct {
	Title       string
	Description string
	OwnerID     int64
}

type UpdateBoardCommand struct {
	ID          int64
	Title       *string
	Description *string
}

type MoveBoardCommand struct {
	Owner int64
}
