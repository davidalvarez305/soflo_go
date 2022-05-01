package models

type User struct {
	ID       uint
	Username string `gorm:"unique"`
	Email    string `gorm:"unique"`
	Password string
}
