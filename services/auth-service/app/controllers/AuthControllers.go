package controllers

import (
	// "context"
	// "os"
	//"time"
	// "encoding/json"
	"encoding/json"
	"log"
	"time"

	// "github.com/jinzhu/copier"

	"auth_service/app/models"
	database "auth_service/config/db"
	utils "auth_service/resources/utility"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	guuid "github.com/google/uuid"
	"gorm.io/gorm"

	amqp_config "auth_service/config/amqp"
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
func (ctrl AuthController) Register(c *fiber.Ctx) error {
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

	queryUser = models.User{Email: payload.Email}
	foundUser = models.User{}
	err = db.First(&foundUser, &queryUser).Error

	if err != gorm.ErrRecordNotFound {
		return c.Status(400).JSON(fiber.Map{
			"code":    400,
			"message": "Email already exits",
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
func (ctrl AuthController) Login(c *fiber.Ctx) error {
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

	// cookie := fiber.Cookie{
	// 	Name: "refresh_jwt",
	// 	Value: RefreshToken,
	// 	HTTPOnly: true,
	// 	Secure: true,
	// 	SameSite: "None",
	// 	Path: "/auth/refresh",
	// 	MaxAge: 24*60*60*1000,
	// }

	// c.Cookie(&cookie)
	return c.JSON(fiber.Map{"user": fiber.Map{"id": userDetail.ID, "role": userDetail.Role, "picture": userDetail.Picture, "name": userDetail.Fullname}, "accessToken": AccessToken, "refreshToken": RefreshToken})
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
	refreshToken := c.Locals("user").(*jwt.Token)

	refreshClaims := refreshToken.Claims.(jwt.MapClaims)
	db := database.DB
	userID, _ := guuid.Parse(refreshClaims["Issuer"].(string))
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

// Logout godoc
//
//	@Summary		Change password
//	@Description	Submit new password to change password
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		models.ChangePasswordForm	true	"Password"
//	@Success		200		{string}  string  "OK"
//	@Failure		400		{string}  error  "Bad Request"
//	@Router			/auth/submitOTP [post]
func (ctrl AuthController) ChangePassword(c *fiber.Ctx) error {
	payload := models.ChangePasswordForm{}
	if err := c.BodyParser(&payload); err != nil {
		return c.SendStatus(502)
	}

	db := database.DB
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
	if !utils.ComparePasswords(payload.OldPassword, userDetail.Password) {

		return c.SendStatus(401)
	}
	payload.NewPassword, _ = utils.CreatePassword(payload.NewPassword)
	db.Model(&models.User{}).Where("username = ?", payload.Username).Update("password", payload.NewPassword)

	return c.JSON(fiber.Map{"msg": "done"})
}

// Logout godoc
//
//	@Summary		Forgot password
//	@Description	User lost their password, they need to change another password
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		models.EmailForm	true	"Email"
//	@Success		200		{string}  string  "OK"
//	@Failure		400		{string}  error  "Bad Request"
//	@Router			/auth/forgot [post]
func (ctrl AuthController) ForgotPassword(c *fiber.Ctx) error {
	payload := models.EmailForm{}
	if err := c.BodyParser(&payload); err != nil {
		return c.SendStatus(502)
	}
	db := database.DB
	userDetail := models.User{}
	queryUser := models.User{Email: payload.Email}
	err := db.First(&userDetail, &queryUser).Error

	if err == gorm.ErrRecordNotFound {
		return c.SendStatus(fiber.StatusNotFound)
	}
	otp, _ := utils.GenerateOTP(6)
	newOTP := models.OTP{
		Email:  payload.Email,
		OTP:    otp,
		Status: 0,
	}
	db.Create(&newOTP)
	msgContent := map[string]string{"to": payload.Email, "otp": otp}
	// jsonMsgContent, _ := json.Marshal(msgContent)
	msgForm := models.MessageForm{
		From:    "Auth_service",
		To:      "Email_service",
		Task:    "OTP",
		Content: msgContent,
	}
	jsonPayload, _ := json.Marshal(msgForm)
	utils.Publish(amqp_config.EmailService.Name, string(jsonPayload))

	return c.JSON(fiber.Map{"msg": "done"})
}

// Logout godoc
//
//	@Summary		Submit OTP
//	@Description	Submit OTP to change password
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		models.OTPForm	true	"Email"
//	@Success		200		{string}  string  "OK"
//	@Failure		400		{string}  error  "Bad Request"
//	@Router			/auth/submitOTP [post]
func (ctrl AuthController) SubmitOTP(c *fiber.Ctx) error {
	payload := models.OTPForm{}
	if err := c.BodyParser(&payload); err != nil {
		return c.SendStatus(502)
	}
	db := database.DB
	otpDetail := models.OTP{}
	time_duration := time.Now().Add(-time.Minute * 3).Unix()
	err := db.Where("created_at > ? and status = ? and email = ?", time_duration, 0, payload.Email).First(&otpDetail).Error
	log.Println(otpDetail)
	log.Println(time_duration)
	if err == gorm.ErrRecordNotFound {
		return c.SendStatus(fiber.StatusNotFound)
	}
	if payload.OTP != otpDetail.OTP {
		return c.SendStatus(fiber.StatusNotFound)
	}
	db.Model(&models.OTP{}).Where("otp = ?", payload.OTP).Update("status", 1)
	AccessToken, _ := utils.CreateOTPToken(payload.OTP)
	return c.JSON(fiber.Map{"accessToken": AccessToken})
}

// Logout godoc
//
//	@Summary		Reset password
//	@Description	Submit new password to change password
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		models.NewPasswordForm	true	"Password"
//	@Success		200		{string}  string  "OK"
//	@Failure		400		{string}  error  "Bad Request"
//	@Router			/auth/submitOTP [post]
func (ctrl AuthController) ResetPassword(c *fiber.Ctx) error {
	payload := models.NewPasswordForm{}
	if err := c.BodyParser(&payload); err != nil {
		return c.SendStatus(502)
	}
	payload.Password, _ = utils.CreatePassword(payload.Password)
	db := database.DB

	db.Model(&models.User{}).Where("email = ?", payload.Email).Update("password", payload.Password)

	return c.JSON(fiber.Map{"msg": "done"})
}

func InitializeAuthController() AuthController {
	return AuthController{}
}
