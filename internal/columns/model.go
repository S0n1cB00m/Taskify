package columns

type Column struct {
	ID      int64 `json:"id" db:"id" example:"101"`
	BoardID int64 `json:"board_id" db:"board_id" example:"6"`
}
