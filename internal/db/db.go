package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

type User struct {
	ID         string
	Imei       string
	TelegramId string
	Email      string
	Phone      string
}

type UserPlace struct {
	ID   uint
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

	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// Проверка подключения к базе данных
	if err = db.Ping(); err != nil {
		log.Fatal("Cannot connect to database:", err)
	}

	// Создание таблиц, если они не существуют
	createTables()
}

func createTables() {
	userTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY,
		imei TEXT,
		telegram_id TEXT,
		email TEXT,
		phone TEXT
	);`

	userPlaceTableQuery := `
	CREATE TABLE IF NOT EXISTS user_places (
		id SERIAL PRIMARY KEY,
		n TEXT,
		e TEXT,
		info TEXT
	);`

	_, err := db.Exec(userTableQuery)
	if err != nil {
		log.Fatal("Error creating users table:", err)
	}

	_, err = db.Exec(userPlaceTableQuery)
	if err != nil {
		log.Fatal("Error creating user_places table:", err)
	}
}

func CreateUser(user *User) (*User, error) {
	user.ID = uuid.New().String()
	query := `INSERT INTO users (id, imei, telegram_id, email, phone) VALUES ($1, $2, $3, $4, $5)`
	_, err := db.Exec(query, user.ID, user.Imei, user.TelegramId, user.Email, user.Phone)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func CreateUserPlace(place *UserPlace) (*UserPlace, error) {
	query := `INSERT INTO user_places (n, e, info) VALUES ($1, $2, $3) RETURNING id`
	err := db.QueryRow(query, place.N, place.E, place.Info).Scan(&place.ID)
	if err != nil {
		return nil, err
	}
	return place, nil
}
