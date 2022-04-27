package controllers

import (
	"github.com/davidalvarez305/soflo_go/server/actions"
	"github.com/gofiber/fiber/v2"
)

func Google(c *fiber.Ctx) error {
	type reqBody struct {
		Searches string `json:"searches"`
	}
	keywordList := [1]string{"hello"}

	var body reqBody
	c.BodyParser(&body)
	keywordList[0] = body.Searches

	q := actions.GoogleQuery{
		Pagesize: 1000,
		KeywordSeed: actions.KeywordSeed{
			Keywords: keywordList,
		},
	}

	keywords := actions.QueryGoogle(q)

	if len(keywords) == 0 {
		return c.Status(400).JSON(fiber.Map{
			"data": "Bad Request.",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": keywords,
	})
}
