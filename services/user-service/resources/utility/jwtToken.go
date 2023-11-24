package utils

import (
	"user_service/app/models"
	"time"
	"os"
	"github.com/golang-jwt/jwt/v4"
)
func CreateAccessToken(user models.User) (string, error) {
	exp := time.Now().Add(time.Second*60).Unix()
	claims := jwt.StandardClaims{
		Issuer: user.ID.String(),
		ExpiresAt: exp,
		Subject: "accessToken",
		IssuedAt: time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(os.Getenv("ACCESS_TOKEN_SECRET")))
	if err != nil {
		return "", err
	}
	return t, nil
}

func CreateRefreshToken(user models.User) (string, error) {
	exp := time.Now().Add(time.Hour*24).Unix()
	claims := jwt.StandardClaims{
		Issuer: user.ID.String(),
		ExpiresAt: exp,
		Subject: "refreshToken",
		IssuedAt: time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(os.Getenv("REFRESH_TOKEN_SECRET")))
	if err != nil {
		return "", err
	}

	return t, nil
}