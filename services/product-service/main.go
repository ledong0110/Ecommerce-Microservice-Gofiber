package main

import (
	_ "fmt"
	"log"
	_ "net/http"
	"os"
	db "product_service/db"
	routes "product_service/routes"

	_ "product_service/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
)

func main() {
	// Load env variables
	godotenv.Load()
	// Connect database
	db.ConnectDB()
	// Initialize server
	app := fiber.New(fiber.Config{})
	// Add Swagger
	app.Get("/prod/swagger/*", swagger.HandlerDefault)
	// Add route
	routes.Route(app)

	log.Fatal(app.Listen(":" + os.Getenv("PORT")))

}
