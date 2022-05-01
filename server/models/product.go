package models

type Product struct {
	ID             uint
	AffiliateUrl   string `gorm:"column:affiliateUrl"`
	ProductPrice   string `gorm:"column:productPrice"`
	ProductReviews string `gorm:"column:productReviews"`
	ProductRatings string `gorm:"column:productRatings"`
	ProductImage   string `gorm:"column:productImage"`
}
