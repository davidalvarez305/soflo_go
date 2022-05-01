package models

type Category struct {
	ID             int
	Title          string
	Slug           string `gorm:"unique"`
	ReviewProducts []ReviewPost
}
