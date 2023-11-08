package routes

import (
	controllers "product_service/app/controllers"

	"github.com/gofiber/fiber/v2"
)

var prodController controllers.ProdController = controllers.InitializeProdController()

func ProdRouter(prod fiber.Router) {
	prod.Post("/Add", prodController.Add)
	prod.Get("/GetById/:id", prodController.GetById)
	prod.Get("/GetByOwner/:id", prodController.GetByOwner)
	prod.Put("/EditById/:id", prodController.EditById)
	prod.Delete("/DeleteById/:id", prodController.DeleteById)
}
