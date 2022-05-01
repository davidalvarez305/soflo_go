package controllers

import (
	"github.com/davidalvarez305/soflo_go/server/handlers"
	"github.com/gofiber/fiber/v2"
)

func DynamicContent(router fiber.Router) {

	content := router.Group("content")
	content.Get("/", handlers.GetContent)
	content.Get("/dictionary", handlers.GetDictionary)
	content.Post("/dynamic", handlers.GetDynamicContent)
}
