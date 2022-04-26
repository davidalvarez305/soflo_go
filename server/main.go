package main

import (
	"github.com/gofiber/fiber/v2"
)

func Router(app *fiber.App) {
	app.Get("/api/", func(c *fiber.Ctx) error {
		msg := "Hello world!"
		return c.JSON(fiber.Map{
			"data": msg,
		})
	})
}

func main() {
	app := fiber.New()

	app.Listen(":4007")
}
