package utils

import (
	"auth_service/app/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateAccessToken(user models.User) (string, error) {

	exp := time.Now().Add(time.Second * 60).Unix()
	claims := jwt.MapClaims{
		"Issuer":   user.ID.String(),
		"exp":      exp,
		"Subject":  "accessToken",
		"IssuedAt": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(os.Getenv("ACCESS_TOKEN_SECRET")))
	if err != nil {
		return "", err
	}
	return t, nil
}

func CreateRefreshToken(user models.User) (string, error) {
	exp := time.Now().Add(time.Hour * 72).Unix()
	claims := jwt.MapClaims{
		"Issuer":   user.ID.String(),
		"exp":      exp,
		"Subject":  "refreshToken",
		"IssuedAt": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(os.Getenv("REFRESH_TOKEN_SECRET")))
	if err != nil {
		return "", err
	}

	return t, nil
}

func CreateOTPToken(otp string) (string, error) {
	exp := time.Now().Add(time.Minute * 5).Unix()
	claims := jwt.MapClaims{
		"Issuer":   otp,
		"exp":      exp,
		"Subject":  "otpToken",
		"IssuedAt": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(os.Getenv("OTP_CREDENTIAL")))
	if err != nil {
		return "", err
	}

	return t, nil
}
