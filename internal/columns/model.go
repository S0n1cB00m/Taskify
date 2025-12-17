package columns

type Column struct {
	Index   int64 `json:"id" db:"id" example:"101"`
	BoardId int64 `json:"board_id" db:"board_id" example:"6"`
}
