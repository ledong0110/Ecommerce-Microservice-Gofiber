package db

import (
	"log"
	"os"
	"auth_service/app/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDB() {
	var err error
	dsn := os.Getenv("DB_USER") + ":" + os.Getenv("DB_PWD") + "@tcp(" + os.Getenv("DB_HOST") + ")/" + os.Getenv("DB_NAME") + "?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	err = DB.AutoMigrate(&models.User{}, &models.OTP{})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database connected !")
}
