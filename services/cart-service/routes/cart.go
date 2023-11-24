package routes

import (
	controllers "cart_service/app/controllers"

	"github.com/gofiber/fiber/v2"
)

var cartController controllers.CartController = controllers.InitializeCartController()

func CartRouter(cart fiber.Router) {
	cart.Post("/Add", cartController.Add)
	cart.Get("/GetByUser/:id", cartController.GetByUser)
	cart.Delete("/DeleteById/:id", cartController.DeleteById)
}
