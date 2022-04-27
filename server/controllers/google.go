package controllers

import (
	"github.com/davidalvarez305/soflo_go/server/actions"
	"github.com/gofiber/fiber/v2"
)

func Google(c *fiber.Ctx) error {
	type reqBody struct {
		Searches string `json:"searches"`
	}

	var body reqBody
	c.BodyParser(&body)

	actions.QueryGoogle()

	return c.Status(200).JSON(fiber.Map{
		"data": "Hello world!",
	})
}
