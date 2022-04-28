package handlers

import (
	"github.com/davidalvarez305/soflo_go/server/actions"
	"github.com/gofiber/fiber/v2"
)

func CrawlAmazon(c *fiber.Ctx) error {
	type reqBody struct {
		Keyword string `json:"keyword"`
	}

	var body reqBody
	c.BodyParser(&body)

	data := actions.ScrapeSearchResultsPage(body.Keyword)

	if len(data) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"data": "No results found",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": data,
	})
}
