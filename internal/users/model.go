package users

type User struct {
	ID       int64  `json:"id" db:"id" example:"101"`
	Email    string `json:"email" db:"email" example:"taskify@example.com"`
	Username string `json:"username" db:"username" example:"johnlennon"`
	Password string `json:"password" db:"password" example:"123456"`
}

type CreateUserDTO struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}
