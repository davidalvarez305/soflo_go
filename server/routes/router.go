package routes

import (
	"github.com/davidalvarez305/soflo_go/server/controllers"
	"github.com/gofiber/fiber/v2"
)

func Router(app *fiber.App) {
	api := app.Group("api")
	google := api.Group("google")

	google.Post("/", controllers.Google)
}
