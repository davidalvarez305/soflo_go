package handlers

import (
	"fmt"
	"strings"

	"github.com/davidalvarez305/soflo_go/server/actions"
	"github.com/davidalvarez305/soflo_go/server/types"
	"github.com/gofiber/fiber/v2"
)

func GetCommercialKeywords(c *fiber.Ctx) error {
	type reqBody struct {
		Keyword string `json:"keyword"`
	}

	var body reqBody
	c.BodyParser(&body)

	s := strings.Split(body.Keyword, "\n")

	if len(s) > 1 {
		return c.Status(400).JSON(fiber.Map{
			"data": "Only one seed keyword allowed per query.",
		})
	}

	q := types.GoogleQuery{
		Pagesize: 1000,
		KeywordSeed: types.KeywordSeed{
			Keywords: [1]string{body.Keyword},
		},
	}

	results := actions.QueryGoogle(q)

	if len(results.Results) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"data": "Bad Request.",
		})
	}

	seedKeywords := actions.GetSeedKeywords(results)
	fmt.Printf("Length of Seed Keywords: %v", seedKeywords)

	if len(seedKeywords) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"data": "No Seed Keywords Found.",
		})
	}

	keywords := actions.GetCommercialKeywords(seedKeywords)
	fmt.Printf("Length of Commercial Keywords: %v", keywords)

	if len(keywords) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"data": "No Commercial Keywords Found.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": keywords,
	})
}

func GetSeedKeywords(c *fiber.Ctx) error {
	type reqBody struct {
		Keyword string `json:"keyword"`
	}

	var body reqBody
	c.BodyParser(&body)

	s := strings.Split(body.Keyword, "\n")

	if len(s) > 1 {
		return c.Status(400).JSON(fiber.Map{
			"data": "Only one seed keyword allowed per query.",
		})
	}

	q := types.GoogleQuery{
		Pagesize: 1000,
		KeywordSeed: types.KeywordSeed{
			Keywords: [1]string{body.Keyword},
		},
	}

	results := actions.QueryGoogle(q)

	if len(results.Results) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"data": "Bad Request.",
		})
	}

	seedKeywords := actions.GetSeedKeywords(results)

	if len(seedKeywords) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"data": "No Seed Keywords Found.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": seedKeywords,
	})
}

func GetPeopleAlsoAsk(c *fiber.Ctx) error {

	keywords := actions.CrawlGoogleSERP("sell my junk car for 500")

	return c.Status(200).JSON(fiber.Map{
		"data": keywords,
	})
}
