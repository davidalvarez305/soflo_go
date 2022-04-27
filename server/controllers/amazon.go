package controllers

import (
	"github.com/davidalvarez305/soflo_go/server/actions"
	"github.com/gofiber/fiber/v2"
)

func Amazon(router fiber.Router) {

	amazon := router.Group("amazon")

	amazon.Post("/crawl", actions.AmazonCrawler)
}
