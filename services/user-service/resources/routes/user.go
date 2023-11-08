package routes

import (
	"github.com/gofiber/fiber/v2"

	controllers "user_service/app/controllers"
)

var userController controllers.UserController = controllers.InitializeUserController()

func UserRouter(user fiber.Router) {
	user.Post("/register", userController.Register)
	user.Get("/refresh", userController.RefreshToken)

	user.Get("/:username", userController.GetUserByUsername)
	user.Get("/", userController.GetAllUsers)
	user.Put("/:username", userController.Update)
	user.Post("/login", userController.Login)
	user.Delete("/:username", userController.DeleteUser)
	// user.Get("*", userController.EmptyPage)
}
