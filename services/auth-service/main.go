package main

import (
	_ "fmt"
	"log"
	_ "net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
    // jwtware "github.com/gofiber/contrib/jwt"
	db "auth_service/config/db"
	routes "auth_service/resources/routes"
	_ "auth_service/resources/utility"
    
)

func main() {  
    // Load env variables
    godotenv.Load()
    // Connect database
    db.ConnectDB()
    // // Initialize server
    app := fiber.New(fiber.Config{})
    // app.Use(jwtware.New(jwtware.Config{
	// 	SigningKey: jwtware.SigningKey{Key: []byte(os.Getenv("ACCESS_TOKEN_SECRET"))},
	// }))
    app.Use(cors.New(cors.Config{
        AllowCredentials: true,
        AllowOrigins: "http://locahost:3000",
        AllowHeaders:  "Origin, Content-Type, Accept",  
    }))
    // app.Static("/", "./public/build")
    
    routes.Route(app)
    
    log.Fatal(app.Listen(":"+os.Getenv("PORT")))
    
}