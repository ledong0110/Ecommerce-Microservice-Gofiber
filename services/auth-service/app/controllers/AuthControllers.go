package controllers

import (
	// "context"
	"os"
	"time"
	// "encoding/json"
	"log"

	// "github.com/jinzhu/copier"

	"gorm.io/gorm"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	guuid "github.com/google/uuid"
	"auth_service/app/models"
	database "auth_service/config/db"
	utils "auth_service/resources/utility"

	
)


type AuthController struct {
}


// Register godoc
//
//	@Summary		Create new user
//	@Description	Create new user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		models.RegisterForm	true	"UserModel"
//	@Success		200		{string}  string  "OK"
//	@Failure		400		{string}  error  "Bad Request"
//	@Router			/auth/register [post]
func (ctrl AuthController) Register(c*fiber.Ctx) error {
	payload := models.RegisterForm{}
	if err := c.BodyParser(&payload); err != nil {
		return c.SendStatus(403)
	}
	// log.Println(payload)
	payload.Password, _ = utils.CreatePassword(payload.Password)
	db := database.DB
	queryUser := models.User{Username: payload.Username}
	foundUser := models.User{}
	err := db.First(&foundUser, &queryUser).Error

	if err != gorm.ErrRecordNotFound {
		return c.Status(400).JSON(fiber.Map{
			"code":    400,
			"message": "User already exits",
		})
	}
	log.Println("Start saving")
	newUser := models.User {
		ID: guuid.New(),
		Username: payload.Username,
		Password: payload.Password,
		Email: payload.Email,
		Role: payload.Role,
		Picture: payload.Picture,
		Fullname: payload.Fullname,
	}
	db.Create(&newUser)
	return c.SendStatus(200)
}


// Login godoc
//
//	@Summary		Login to the system
//	@Description	Login to the system
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		models.LoginForm	true	"LoginForm"
//	@Success		200		{string}  string  "OK"
//	@Failure		400		{string}  error  "Bad Request"
//	@Router			/auth/login [post]
func (ctrl AuthController) Login(c*fiber.Ctx) error {
	payload := models.LoginForm{}
	
	if err := c.BodyParser(&payload); err != nil {
		return c.SendStatus(403)
	}
	// log.Println(payload)
	db := database.DB
	// log.Println(string(acc))
	// opts := options.FindOne().SetProjection(bson.M{
		
	// 	"online": 0,
	// })
	userDetail := models.User{}
	queryUser := models.User{Username: payload.Username}
	err := db.First(&userDetail, &queryUser).Error

	if err == gorm.ErrRecordNotFound {
		return c.Status(404).JSON(fiber.Map{
			"code":    404,
			"message": "User not found",
		})
	}
	if err != nil {
		
		return c.SendStatus(401)
	}
	if utils.ComparePasswords(payload.Password, userDetail.Password){
		log.Println("Successfully authentication")
	} else {
		log.Println("Wrong Password")
		c.ClearCookie("refresh_jwt")
		return c.SendStatus(401)
	}
	RefreshToken, err := utils.CreateRefreshToken(userDetail)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	AccessToken, err := utils.CreateAccessToken(userDetail)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	
	
	cookie := fiber.Cookie{
		Name: "refresh_jwt", 
		Value: RefreshToken,
		HTTPOnly: true,
		Secure: true,
		SameSite: "None",
		Path: "/auth/refresh",
		MaxAge: 24*60*60*1000,
	}
	
	c.Cookie(&cookie)
	return c.JSON(fiber.Map{"user": fiber.Map{"id": userDetail.ID ,"role": userDetail.Role, "picture": userDetail.Picture, "name": userDetail.Fullname}, "accessToken": AccessToken})
}

// Logout godoc
//
//	@Summary		Logout to the system
//	@Description	Logout to the system
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Success		200		{string}  string  "OK"
//	@Failure		400		{string}  error  "Bad Request"
//	@Router			/auth/logout [get]
func (ctrl AuthController) Logout(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_jwt", "none")
	if refreshToken == "none" {
		return c.SendStatus(200)
	}
	c.ClearCookie("refresh_jwt")
	return c.SendStatus(200)	
}

// Refresh godoc
//
//	@Summary		Refresh new access token
//	@Description	Refresh new access token
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Success		200		{string}  string  "OK"
//	@Failure		400		{string}  error  "Bad Request"
//	@Router			/auth/refresh [get]
func (ctrl AuthController) RefreshToken(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_jwt", "none")
	if refreshToken == "none" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}		
	refreshClaims := jwt.StandardClaims{}
	token, _ := jwt.ParseWithClaims(refreshToken, &refreshClaims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("REFRESH_TOKEN_SECRET")), nil
		},
	)
	
	if token.Valid {
		if refreshClaims.ExpiresAt < time.Now().Unix() {
			c.ClearCookie("refresh_jwt")
			return c.SendStatus(fiber.StatusUnauthorized)
		}
	} else {
		c.ClearCookie(("refresh_jwt"))
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	db := database.DB
	userID, _ := guuid.Parse(refreshClaims.Issuer)
	userDetail := models.User{}
	queryUser := models.User{ID: userID}
	err := db.First(&userDetail, &queryUser).Error

	if err == gorm.ErrRecordNotFound {
		return c.SendStatus(fiber.StatusNotFound)
	}
	if err != nil {
		
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	AccessToken, _ := utils.CreateAccessToken(userDetail)
	
	return c.JSON(fiber.Map{"user": fiber.Map{"id": userDetail.ID ,"role": userDetail.Role, "picture": userDetail.Picture, "name": userDetail.Fullname}, "accessToken": AccessToken})
}

func InitializeAuthController() AuthController {
	return AuthController{}
}
