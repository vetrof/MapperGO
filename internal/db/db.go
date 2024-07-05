package db

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var db *gorm.DB
var err error

type User struct {
	ID         string `gorm:"primarykey"`
	Imei       string
	TelegramId string
	Email      string
	Phone      string
}

type UserPlace struct {
	ID   uint `gorm:"primaryKey;autoIncrement"`
	N    string
	E    string
	Info string
}

func InitPostgresDB() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	var (
		host     = os.Getenv("DB_HOST")
		port     = os.Getenv("DB_PORT")
		dbUser   = os.Getenv("DB_USER")
		dbName   = os.Getenv("DB_NAME")
		password = os.Getenv("DB_PASSWORD")
	)
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		host,
		port,
		dbUser,
		dbName,
		password,
	)

	db, err = gorm.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(User{})
	db.AutoMigrate(UserPlace{})
}

func CreateUser(user *User) (*User, error) {
	user.ID = uuid.New().String()
	res := db.Create(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	return user, nil
}

func CreateUserPlace(place *UserPlace) (*UserPlace, error) {
	//place.E = "123"
	//place.N = "456"
	//place.Info = "456"
	res := db.Create(&place)
	if res.Error != nil {
		return nil, res.Error
	}
	return place, nil
}
