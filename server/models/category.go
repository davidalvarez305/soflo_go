package models

type Category struct {
	ID             int
	Name           string
	Slug           string `gorm:"unique"`
	ReviewProducts []ReviewPost
}
