package main

import (
	_ "fmt"
	"log"
	_ "net/http"
	"os"
    "github.com/gofiber/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
    // jwtware "github.com/gofiber/contrib/jwt"
	db "auth_service/config/db"
	routes "auth_service/resources/routes"
	_ "auth_service/resources/utility"
    _ "auth_service/docs"
)

// @title Ecommerce Authentication Service API
// @version 1.0
// @description This is a sample swagger for Fiber
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email fiber@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host 127.0.0.1:8000
// @BasePath /
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
        // AllowOrigins: "http://locahost:3000",
        AllowHeaders:  "Origin, Content-Type, Accept",  
    }))
    // app.Static("/", "./public/build")
    
    app.Get("/swagger/*", swagger.HandlerDefault)

	// app.Get("/swagger/*", swagger.New(swagger.Config{ // custom
	// 	URL: "http://example.com/doc.json",
	// 	DeepLinking: false,
	// 	// Expand ("list") or Collapse ("none") tag groups by default
	// 	DocExpansion: "none",
	// 	// Prefill OAuth ClientId on Authorize popup
	// 	OAuth: &swagger.OAuthConfig{
	// 		AppName:  "OAuth Provider",
	// 		ClientId: "21bb4edc-05a7-4afc-86f1-2e151e4ba6e2",
	// 	},
	// 	// Ability to change OAuth2 redirect uri location
	// 	OAuth2RedirectUrl: "http://localhost:8080/swagger/oauth2-redirect.html",
	// }))

    routes.Route(app)
    
    log.Fatal(app.Listen(":"+os.Getenv("PORT")))
    
}