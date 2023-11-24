package utils

import (
    "errors"

    "golang.org/x/crypto/bcrypt"
    "crypto/rand"
)

const otpChars = "1234567890"
// CreatePassword will create password using bcrypt
func CreatePassword(passwordString string) (string, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordString), 8)
    if err != nil {
        return "", errors.New("Error occurred while creating a Hash")
    }

    return string(hashedPassword), nil
}

// ComparePasswords will create password using bcrypt
func ComparePasswords(password string, hashedPassword string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    return err == nil 
}

func GenerateOTP(length int) (string, error) {
    buffer := make([]byte, length)
    _, err := rand.Read(buffer)
    if err != nil {
        return "", err
    }

    otpCharsLength := len(otpChars)
    for i := 0; i < length; i++ {
        buffer[i] = otpChars[int(buffer[i])%otpCharsLength]
    }

    return string(buffer), nil
}