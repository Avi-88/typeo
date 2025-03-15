package config

import (
	"fmt"
	"log"
	"typeo/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	LoadEnv()
	dsn := GetEnv("DB_URL")
	var err error
	DB,err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	
	if err != nil {
		log.Fatal("Error connecting to the database")
	}

	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	fmt.Println("Successfully connected to the database!")
}