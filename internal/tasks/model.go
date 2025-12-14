package tasks

type Task struct {
	ID          int64  `json:"id" db:"id" example:"101"`
	ColumnID    int64  `json:"column_id" db:"column_id" example:"34"`
	Title       string `json:"title" db:"title" example:"/POST deposit implementation"`
	Description string `json:"description" db:"description" example:"/POST deposit implementation by dev1 and dev2. The deadline is until next week."`
	AssigneeID  int64  `json:"assignee_id" db:"assignee_id" example:"7"`
	ColumnOrder int64  `json:"column_order" db:"column_order" example:"1"`
}
