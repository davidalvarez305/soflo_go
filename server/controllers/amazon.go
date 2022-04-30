package controllers

import (
	"github.com/davidalvarez305/soflo_go/server/handlers"
	"github.com/gofiber/fiber/v2"
)

func Amazon(router fiber.Router) {

	amazon := router.Group("amazon")

	amazon.Post("/crawl", handlers.CrawlAmazon)
	amazon.Post("/paapi5", handlers.SearchPAAPI5)
}
