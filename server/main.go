package main

import (
	"github.com/davidalvarez305/soflo_go/server/database"
	"github.com/davidalvarez305/soflo_go/server/routes"
	"github.com/davidalvarez305/soflo_go/server/sessions"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:4007",
	}))
	database.Connect()
	sessions.Init()
	routes.Router(app)

	app.Listen(":4007")
}
