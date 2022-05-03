package utils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/davidalvarez305/soflo_go/server/database"
	"github.com/davidalvarez305/soflo_go/server/models"
	"github.com/davidalvarez305/soflo_go/server/types"
	"gorm.io/gorm/clause"
)

func CreateCategorySlug(str string) string {
	var final string
	r := regexp.MustCompile(`[a-z0-9]+`)
	res := r.FindAllString(str, -1)

	if len(res) > 0 {
		final = strings.Join(res, "-")
	}

	return final
}

func CreateReviewPostFields(input types.AmazonSearchResultsPage, dictionary []types.Dictionary, sentences []types.DynamicContent, categoryId int) models.ReviewPost {
	name := strings.Join(strings.Split(strings.ToLower(input.Name), " "), "-")
	slug := CreateCategorySlug(name)
	replacedImage := strings.Replace(input.Image, "UL320", "UL640", 1)

	data := GenerateContentUtil(input.Name, dictionary, sentences)

	fields := models.ReviewPost{
		Title:                         data.ReviewPostTitle,
		CategoryID:                    categoryId,
		Slug:                          slug,
		Content:                       data.ReviewPostContent,
		Headline:                      data.ReviewPostHeadline,
		Intro:                         data.ReviewPostIntro,
		Description:                   data.ReviewPostDescription,
		ProductLabel:                  data.ReviewPostProductLabel,
		ProductName:                   input.Name,
		ProductDescription:            data.ReviewPostProductDescription,
		ProductAffiliateUrl:           input.Link,
		Faq_Answer_1:                  data.ReviewPostFaq_Answer_1,
		Faq_Answer_2:                  data.ReviewPostFaq_Answer_2,
		Faq_Answer_3:                  data.ReviewPostFaq_Answer_3,
		Faq_Question_1:                data.ReviewPostFaq_Question_1,
		Faq_Question_2:                data.ReviewPostFaq_Question_2,
		Faq_Question_3:                data.ReviewPostFaq_Question_3,
		HorizontalCardProductImageUrl: replacedImage,
		HorizontalCardProductImageAlt: strings.ToLower(input.Name),
	}

	return fields
}

func insertProducts(products []types.AmazonSearchResultsPage) ([]models.Product, error) {
	var p []models.Product
	for i := 0; i < len(products); i++ {
		product := models.Product{
			AffiliateUrl:   products[i].Link,
			ProductPrice:   products[i].Price,
			ProductReviews: products[i].Reviews,
			ProductRatings: products[i].Rating,
			ProductImage:   products[i].Image,
		}
		p = append(p, product)
	}

	fmt.Printf("Inserting %v products!", len(p))
	ins := database.DB.Clauses(clause.OnConflict{DoNothing: true}).Save(p)

	if ins.Error != nil {
		return nil, ins.Error
	}

	return p, nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func insertCategories(products []types.AmazonSearchResultsPage) ([]models.Category, error) {
	var categories []models.Category
	var c []string
	for i := 0; i < len(products); i++ {
		if !contains(c, products[i].Category) {
			c = append(c, products[i].Category)
		}
	}

	for _, a := range c {
		cat := models.Category{
			Name: strings.Title(a),
			Slug: CreateCategorySlug(a),
		}
		categories = append(categories, cat)
	}

	fmt.Printf("Inserting %v categories!", len(c))
	db := database.DB.Clauses(clause.OnConflict{DoNothing: true}).Save(categories)

	if db.Error != nil {
		fmt.Println("Error while trying to insert categories.")
		return nil, db.Error
	}

	sel := database.DB.Raw("SELECT * FROM categories;").Scan(&categories)

	if sel.Error != nil {
		return nil, sel.Error
	}

	return categories, nil
}

func InsertReviewPosts(products []types.AmazonSearchResultsPage, dictionary []types.Dictionary, sentences []types.DynamicContent) error {
	var posts []models.ReviewPost

	_, err2 := insertProducts(products)

	if err2 != nil {
		return err2
	}

	c, err := insertCategories(products)

	if err != nil {
		return err
	}

	for i := 0; i < len(products); i++ {
		var categoryId int
		for _, cat := range c {
			if strings.Title(products[i].Category) == cat.Name {
				categoryId = cat.ID
			}
		}
		p := CreateReviewPostFields(products[i], dictionary, sentences, categoryId)
		posts = append(posts, p)
	}

	db := database.DB.Clauses(clause.OnConflict{DoNothing: true}).Save(posts)

	if db.Error != nil {
		return db.Error
	}

	return nil
}
