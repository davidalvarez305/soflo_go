package handlers

import (
	"fmt"

	"github.com/davidalvarez305/soflo_go/server/actions"
	"github.com/gofiber/fiber/v2"
)

func CreatePosts(c *fiber.Ctx) error {
	var products []actions.AmazonSearchResultsPage
	type reqBody struct {
		Keyword string `json:"keyword"`
	}

	var body reqBody
	c.BodyParser(&body)

	q := actions.GoogleQuery{
		Pagesize: 50,
		KeywordSeed: actions.KeywordSeed{
			Keywords: [1]string{body.Keyword},
		},
	}

	KW_ARR := actions.QueryGoogle(q)

	seedKeywords := actions.GetSeedKeywords(KW_ARR)
	commercialKeywords := actions.GetCommercialKeywords(seedKeywords)

	for i := 0; i < len(commercialKeywords); i++ {
		data := actions.ScrapeSearchResultsPage(commercialKeywords[i])
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
	return c.Status(200).JSON(fiber.Map{
		"data": products,
	})
}
