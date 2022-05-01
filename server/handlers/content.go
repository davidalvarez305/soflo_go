package handlers

import (
	"github.com/davidalvarez305/soflo_go/server/actions"
	"github.com/davidalvarez305/soflo_go/server/utils"
	"github.com/gofiber/fiber/v2"
)

func GetContent(c *fiber.Ctx) error {
	sentences := actions.PullDynamicContent()

	if len(sentences) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"data": "No results found",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": sentences,
	})
}

func GetDictionary(c *fiber.Ctx) error {
	data := actions.PullContentDictionary()

	if len(data) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"data": "No results found",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": data,
	})
}

func GetDynamicContent(c *fiber.Ctx) error {
	type reqBody struct {
		ProductName string `json:"productName"`
	}
	var body reqBody
	c.BodyParser(&body)

	data := actions.PullContentDictionary()
	sentences := actions.PullDynamicContent()

	content := utils.GenerateContentUtil(body.ProductName, data, sentences)

	if len(sentences) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"data": "No results found",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": content,
	})
}
