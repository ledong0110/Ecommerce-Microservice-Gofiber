package routes

import (
	"github.com/gofiber/fiber/v2"
)

func Route(app *fiber.App) {
	prod := app.Group("/prod")
	ProdRouter(prod)
}
