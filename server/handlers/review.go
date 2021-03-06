package handlers

import (
	"github.com/davidalvarez305/soflo_go/server/actions"
	"github.com/gofiber/fiber/v2"
)

func CreatePosts(c *fiber.Ctx) error {
	type reqBody struct {
		Keyword string `json:"keyword"`
	}

	var body reqBody
	c.BodyParser(&body)

	products, err := actions.CreateReviewPosts(body.Keyword)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"data": err,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": products,
	})
}
