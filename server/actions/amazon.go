package actions

import "github.com/gofiber/fiber/v2"

func AmazonCrawler(c *fiber.Ctx) error {
	type reqBody struct {
		Searches string `json:"searches"`
	}

	var body reqBody
	c.BodyParser(&body)

	s := body.Searches

	return c.Status(200).JSON(fiber.Map{
		"data": s,
	})
}
