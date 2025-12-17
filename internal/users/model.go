package users

import "golang.org/x/crypto/bcrypt"

type User struct {
	Id       int64  `json:"id" db:"id" example:"101"`
	Email    string `json:"email" db:"email" example:"taskify@example.com"`
	Username string `json:"username" db:"username" example:"johnlennon"`
	Password string `json:"-" db:"password" example:"123456"`
}

type CreateUserDTO struct {
	Email    string `json:"email" example:"taskify@example.com"`
	Username string `json:"username" example:"johnlennon"`
	Password string `json:"password" example:"123456"`
}

// HashPassword хеширует пароль и записывает его в структуру
func (u *User) HashPassword(plainPassword string) error {
	// Cost (сложность) - bcrypt.DefaultCost (обычно 10)
	bytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}

// CheckPassword сравнивает присланный пароль с сохраненным хешем
func (u *User) CheckPassword(plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainPassword))
	return err == nil
}
