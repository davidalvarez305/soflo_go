package handlers

import (
	"github.com/davidalvarez305/soflo_go/server/actions"
	"github.com/davidalvarez305/soflo_go/server/utils"
	"github.com/gofiber/fiber/v2"
)

func GetContent(c *fiber.Ctx) error {
	data := actions.PullDynamicContent()

	if len(data) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"data": "No results found",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": data,
	})
}

func GetDictionary(c *fiber.Ctx) error {
	data := actions.PullContentDictionary()
	sentences := actions.PullDynamicContent()

	s := utils.GenerateContentUtil(data, sentences)
	if len(data) == 0 {
		return c.Status(404).JSON(fiber.Map{
			"data": "No results found",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": s,
	})
}
