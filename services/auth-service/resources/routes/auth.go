package routes

import (
	"github.com/gofiber/fiber/v2"

	controllers "auth_service/app/controllers"
	
)
var authController controllers.AuthController= controllers.InitializeAuthController()

func AuthRouter(auth fiber.Router) {
	auth.Post("/register", authController.Register)
	auth.Get("/refresh", authController.RefreshToken)
	auth.Get("/logout",  authController.Logout)
	auth.Post("/login", authController.Login)
	// auth.Get("*", authController.EmptyPage)
}