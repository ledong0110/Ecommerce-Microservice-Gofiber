package controllers

import (
	// "context"
	"os"
	"time"

	// "encoding/json"
	"log"

	// "github.com/jinzhu/copier"

	"user_service/app/models"
	database "user_service/config/db"
	utils "user_service/resources/utility"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	guuid "github.com/google/uuid"
	"gorm.io/gorm"
)

type UserController struct {
}

// Register godoc
//
//	@Summary		Create new user
//	@Description	Create new user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		models.RegisterForm	true	"UserModel"
//	@Success		200		{string}  string  "OK"
//	@Failure		400		{string}  error  "Bad Request"
//	@Router			/user/register [post]
func (ctrl UserController) Register(c *fiber.Ctx) error {
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
	newUser := models.User{
		ID:       guuid.New(),
		Username: payload.Username,
		Password: payload.Password,
		Email:    payload.Email,
		Role:     payload.Role,
		Picture:  payload.Picture,
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
func (ctrl UserController) Login(c *fiber.Ctx) error {
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
	if utils.ComparePasswords(payload.Password, userDetail.Password) {
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
		Name:     "refresh_jwt",
		Value:    RefreshToken,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "None",
		Path:     "/user/refresh",
		MaxAge:   24 * 60 * 60 * 1000,
	}

	c.Cookie(&cookie)
	return c.JSON(fiber.Map{"user": fiber.Map{"id": userDetail.ID, "role": userDetail.Role, "picture": userDetail.Picture, "name": userDetail.Fullname}, "accessToken": AccessToken})
}

// Logout godoc
//
//	@Summary		Logout to the system
//	@Description	Logout to the systems
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Success		200		{string}  string  "OK"
//	@Failure		400		{string}  error  "Bad Request"
//	@Router			/auth/logout [get]
func (ctrl UserController) Logout(c *fiber.Ctx) error {
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
func (ctrl UserController) RefreshToken(c *fiber.Ctx) error {
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

	return c.JSON(fiber.Map{"user": fiber.Map{"id": userDetail.ID, "role": userDetail.Role, "picture": userDetail.Picture, "name": userDetail.Fullname}, "accessToken": AccessToken})
}

func InitializeUserController() UserController {
	return UserController{}
}

// Get User godoc
//
//	@Summary		Get Infor User By Username
//	@Description	Get Infor User By Username
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			username	path		string  true	"username"
//	@Success		200		{string}  string  "OK"
//	@Failure		400		{string}  error  "Bad Request"
//	@Router			/user/{username} [get]
func (ctrl UserController) GetUserByUsername(c *fiber.Ctx) error {
	db := database.DB
	// // log.Println(string(acc))
	// // opts := options.FindOne().SetProjection(bson.M{

	// // 	"online": 0,
	// // })

	// userDetail := models.User{}
	queryUser := models.User{Username: c.Params("username")}
	foundUser := models.User{}
	err := db.First(&foundUser, &queryUser).Error

	if err == gorm.ErrRecordNotFound {
		return c.Status(404).JSON(fiber.Map{
			"code":    404,
			"message": "User not found",
		})
	}
	if err != nil {
		return c.Status(200).JSON(fiber.Map{"user": fiber.Map{"id": foundUser.ID, "role": foundUser.Role, "picture": foundUser.Picture, "name": foundUser.Fullname, "email": foundUser.Email}})
	}
	return c.Status(200).JSON(fiber.Map{"user": fiber.Map{"id": foundUser.ID, "role": foundUser.Role, "picture": foundUser.Picture, "name": foundUser.Fullname, "email": foundUser.Email}})
}

func CreateResponseUser(userModel models.User) models.User {
	return models.User{ID: userModel.ID, Fullname: userModel.Fullname, Username: userModel.Username, Email: userModel.Email, Role: userModel.Role, Picture: userModel.Picture}
}

// GetAll godoc
//
//	@Summary		Get All Users
//	@Description	Get All Users
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Success		200		{string}  string  "OK"
//	@Failure		400		{string}  error  "Bad Request"
//	@Router			/user/ [get]
func (ctrl UserController) GetAllUsers(c *fiber.Ctx) error {
	db := database.DB

	users := []models.User{}

	db.Find(&users)
	responseUsers := []models.User{}
	for _, user := range users {
		responseUser := CreateResponseUser(user)
		responseUsers = append(responseUsers, responseUser)
	}
	return c.Status(200).JSON(responseUsers)
}

// Update godoc
//
//	@Summary		Update User
//	@Description	Update User
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			username	path		string  true	"username"
//	@Param			payload	body		models.UpdateForm	true	"UpdateUserForm"
//	@Success		200		{string}  string  "OK"
//	@Failure		400		{string}  error  "Bad Request"
//	@Router			/user/{username} [put]
func (ctrl UserController) Update(c *fiber.Ctx) error {
	db := database.DB

	// userDetail := models.User{}
	queryUser := models.User{Username: c.Params("username")}
	foundUser := models.User{}
	err := db.First(&foundUser, &queryUser).Error

	if err == gorm.ErrRecordNotFound {
		return c.Status(404).JSON(fiber.Map{
			"code":    404,
			"message": "User not found",
		})
	}
	type UpdateUser struct {
		Fullname string `json:"fullname"`
		Email    string `json:"email"`
		Role     string `json:"role"`
		Picture  string `json:"picture"`
		// Password  string     `json:"-"`
	}
	var updateData UpdateUser
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(500).JSON(err.Error())
	}
	foundUser.Fullname = updateData.Fullname
	foundUser.Role = updateData.Role
	foundUser.Email = updateData.Email
	foundUser.Picture = updateData.Picture

	db.Save(&foundUser)

	responseUSer := CreateResponseUser(foundUser)
	return c.Status(200).JSON(responseUSer)
}

// DeleteUser godoc
//
//	@Summary		Delete a user
//	@Description	Delete a user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			username	path		string  true	"username"
//	@Success		200		{string}  string  "OK"
//	@Failure		400		{string}  error  "Bad Request"
//	@Router			/user/{username} [delete]
func (ctrl UserController) DeleteUser(c *fiber.Ctx) error {
	db := database.DB

	// userDetail := models.User{}
	queryUser := models.User{Username: c.Params("username")}
	foundUser := models.User{}
	err := db.First(&foundUser, &queryUser).Error

	if err == gorm.ErrRecordNotFound {
		return c.Status(404).JSON(fiber.Map{
			"code":    404,
			"message": "User not found",
		})
	}

	if err := db.Delete(&foundUser).Error; err != nil {
		return c.Status(404).JSON(err.Error())
	}

	return c.Status(200).SendString("Successfully Deleted User")
}
