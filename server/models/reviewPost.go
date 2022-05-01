package models

type ReviewPost struct {
	ID                            uint
	Title                         string
	Slug                          string
	Content                       string
	CategoryID                    int
	Headline                      string
	Intro                         string
	Description                   string
	ProductLabel                  string
	ProductName                   string
	ProductDescription            string
	ProductAffiliateUrl           string
	Faq_Answer_1                  string
	Faq_Answer_2                  string
	Faq_Answer_3                  string
	Faq_Question_1                string
	Faq_Question_2                string
	Faq_Question_3                string
	HorizontalCardProductImageUrl string
	HorizontalCardProductImageAlt string
}
