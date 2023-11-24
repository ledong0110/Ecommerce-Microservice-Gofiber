package main

import (
	db "cart_service/db"
	_ "fmt"
	"log"
	_ "net/http"
	"os"

	routes "cart_service/routes"

	// _ "product_service/docs"

	"github.com/gofiber/fiber/v2"
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
	// app.Get("/swagger/*", swagger.HandlerDefault)
	// Add route
	routes.Route(app)

	log.Fatal(app.Listen(":" + os.Getenv("PORT")))

}
