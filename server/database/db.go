package database

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/davidalvarez305/soflo_go/server/models"
)

type DBInstance = *gorm.DB

var DB DBInstance

type connection struct {
	Host     string
	Port     string
	User     string
	Password string
	DB       string
}

func Connect() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading dot env file")
	}

	conn := connection{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DB:       os.Getenv("POSTGRES_DB"),
	}

	db, err := gorm.Open(postgres.Open(connToString(conn)), &gorm.Config{})

	if err != nil {
		fmt.Printf("Error connecting to the DB: %s\n", err.Error())
		return
	} else {
		fmt.Printf("Connected to Database.")
	}

	db.AutoMigrate(&models.Product{}, &models.Category{}, &models.ReviewPost{})

	DB = db
}

func connToString(info connection) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		info.Host, info.Port, info.User, info.Password, info.DB)
}
