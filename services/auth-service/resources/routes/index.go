package routes

import (
	// middleware "auth_service/app/middlewares"
	// "os"

	"github.com/gofiber/fiber/v2"
	// jwtware "github.com/gofiber/jwt/v3"
)

func Route(app *fiber.App) {
	// task := app.Group("/task")
	// task.Use(jwtware.New(jwtware.Config{
	// 	SigningKey: []byte(os.Getenv("ACCESS_TOKEN_SECRET")),
	// }))
	// TaskRouter(task)

	auth := app.Group("/auth")
	AuthRouter(auth)
}
