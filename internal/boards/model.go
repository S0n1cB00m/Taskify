package boards

type Board struct {
	ID          int64  `json:"id" db:"id" example:"101"`
	Name        string `json:"daily" db:"name" example:"Deposit"`
	Description string `json:"description" db:"description" example:"Deposit feature development"`
	OwnerID     int64  `json:"owner_id" db:"owner_id" example:"2"`
}
