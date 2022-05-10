package controllers

import (
	"github.com/davidalvarez305/soflo_go/server/handlers"
	"github.com/gofiber/fiber/v2"
)

func Google(router fiber.Router) {
	google := router.Group("google")

	google.Post("/keywords", handlers.GetCommercialKeywords)
	google.Post("/seed", handlers.GetSeedKeywords)
	google.Post("/crawl", handlers.GetPeopleAlsoAsk)
}
