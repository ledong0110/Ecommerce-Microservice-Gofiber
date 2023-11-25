package routes

import (
	controllers "auth_service/app/controllers"
	"os"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

var authController controllers.AuthController = controllers.InitializeAuthController()

func AuthRouter(auth fiber.Router) {
	auth.Post("/register", authController.Register)
	auth.Get("/logout", authController.Logout)
	auth.Post("/login", authController.Login)
	auth.Post("/changePassword", jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(os.Getenv("ACCESS_TOKEN_SECRET"))},
	}),
	authController.ChangePassword)
	auth.Post("/forgot", authController.ForgotPassword)
	auth.Post("/submitOTP", authController.SubmitOTP)
	
	auth.Get("/refresh", jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(os.Getenv("REFRESH_TOKEN_SECRET"))},
	}),
	authController.RefreshToken)
	
	auth.Post("/reset_password", jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(os.Getenv("OTP_CREDENTIAL"))},
	}),
	authController.ResetPassword)
	// auth.Get("*", authController.EmptyPage)
}
