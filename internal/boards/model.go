package boards

type Board struct {
	Index       int64  `json:"id" db:"index" example:"101"`
	Name        string `json:"name" db:"name" example:"Deposit"`
	Description string `json:"description" db:"description" example:"Deposit feature development"`
	UserId      int64  `json:"-" db:"user_id" example:"2"`
}

type CreateBoardDTO struct {
	Name        string `json:"name" example:"Deposit"`
	Description string `json:"description" example:"Deposit feature development"`
}
