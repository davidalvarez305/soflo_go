package actions

import (
	"fmt"
	"strings"

	"github.com/davidalvarez305/soflo_go/server/database"
	"github.com/davidalvarez305/soflo_go/server/models"
	"github.com/davidalvarez305/soflo_go/server/types"
	"github.com/davidalvarez305/soflo_go/server/utils"
	"gorm.io/gorm/clause"
)

func GetProducts(keyword string) []types.AmazonSearchResultsPage {
	var products []types.AmazonSearchResultsPage

	q := types.GoogleQuery{
		Pagesize: 50,
		KeywordSeed: types.KeywordSeed{
			Keywords: [1]string{keyword},
		},
	}

	KW_ARR := QueryGoogle(q)

	seedKeywords := GetSeedKeywords(KW_ARR)
	commercialKeywords := GetCommercialKeywords(seedKeywords)

	for i := 0; i < len(commercialKeywords); i++ {
		data := ScrapeSearchResultsPage(commercialKeywords[i])
		if len(data) == 0 {
			fmt.Println("Keyword: " + commercialKeywords[i] + "0")
		}
		if len(data) > 0 {
			products = append(products, data...)
		}
		total := fmt.Sprintf("Keyword #%v of %v - %s - Total Products = %v", i+1, len(commercialKeywords), commercialKeywords[i], len(data))
		fmt.Println(total)
	}

	productsTotal := fmt.Sprintf("Total Products = %v", len(products))
	fmt.Println(productsTotal)

	return products
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

	ins := database.DB.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(p, len(p))

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
			Title: strings.Title(a),
			Slug:  utils.CreateCategorySlug(a),
		}
		categories = append(categories, cat)
	}

	db := database.DB.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(categories, len(categories))

	if db.Error != nil {
		fmt.Println("Error while trying to insert categories.")
		return nil, db.Error
	}

	sel := database.DB.Raw("SELECT * FROM category;").Scan(&categories)

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
		for _, a := range c {
			if strings.Title(products[i].Category) == a.Title {
				categoryId = a.ID
			}
		}
		p := utils.CreateReviewPostFields(products[i], dictionary, sentences, categoryId)
		posts = append(posts, p)
	}

	db := database.DB.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(posts, len(posts))

	if db.Error != nil {
		return db.Error
	}

	return nil
}
