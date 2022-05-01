package models

type ReviewPost struct {
	ID                            uint
	Title                         string
	Slug                          string `gorm:"unique"`
	Content                       string
	CategoryID                    int `gorm:"column:categoryId"`
	Headline                      string
	Intro                         string
	Description                   string
	ProductLabel                  string `gorm:"column:productlabel"`
	ProductName                   string `gorm:"column:productname"`
	ProductDescription            string `gorm:"column:productdescription"`
	ProductAffiliateUrl           string `gorm:"column:productaffiliateurl"`
	Faq_Answer_1                  string
	Faq_Answer_2                  string
	Faq_Answer_3                  string
	Faq_Question_1                string
	Faq_Question_2                string
	Faq_Question_3                string
	HorizontalCardProductImageUrl string `gorm:"column:horizontalcardproductimageurl"`
	HorizontalCardProductImageAlt string `gorm:"column:horizontalcardproductimagealt"`
}
