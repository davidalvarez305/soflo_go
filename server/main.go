package main

import (
	"github.com/davidalvarez305/soflo_go/server/actions"
	"github.com/davidalvarez305/soflo_go/server/database"
	"github.com/davidalvarez305/soflo_go/server/sessions"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func Router(app *fiber.App) {
	app.Get("/api/", func(c *fiber.Ctx) error {
		data, err := actions.RefreshAuthToken()
		if err != nil {
			c.Status(400).JSON(fiber.Map{
				"data": "Error while trying to get Google Auth Token.",
			})
		}

		return c.Status(200).JSON(fiber.Map{
			"data": data,
		})
	})
}

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:4007",
	}))
	database.Connect()
	sessions.Init()
	Router(app)

	app.Listen(":4007")
}
