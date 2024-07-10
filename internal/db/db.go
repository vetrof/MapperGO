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

type Place struct {
	ID   uint
	Name string
	Geom string
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
		schema   = os.Getenv("DB_SCHEMA")
	)
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable search_path=%s",
		host,
		port,
		dbUser,
		dbName,
		password,
		schema,
	)

	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// Проверка подключения к базе данных
	if err = db.Ping(); err != nil {
		log.Fatal("Cannot connect to database:", err)
	} else {
		log.Println("connect to database OK")
	}

	// Создание таблиц, если они не существуют
	createPostGIS()
	createSchema(schema)
	createTables()
}

func createPostGIS() {
	_, err := db.Exec("CREATE EXTENSION IF NOT EXISTS postgis;")
	if err != nil {
		log.Fatal("Error creating PostGIS extension:", err)
	}
}

func createSchema(schema string) {
	schemaQuery := fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS %s;`, schema)
	_, err := db.Exec(schemaQuery)
	if err != nil {
		log.Fatal("Error creating schema:", err)
	}
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

	placeTableQuery := `
	CREATE TABLE IF NOT EXISTS places (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100),
		geom GEOGRAPHY(Point, 4326)
	);`

	_, err := db.Exec(userTableQuery)
	if err != nil {
		log.Fatal("Error creating users table:", err)
	}

	_, err = db.Exec(userPlaceTableQuery)
	if err != nil {
		log.Fatal("Error creating user_places table:", err)
	}

	_, err = db.Exec(placeTableQuery)
	if err != nil {
		log.Fatal("Error creating places table:", err)
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

func CreatePlace(place *Place) (*Place, error) {
	query := `INSERT INTO places (name, geom) VALUES ($1, ST_GeogFromText($2)) RETURNING id`
	err := db.QueryRow(query, place.Name, place.Geom).Scan(&place.ID)
	if err != nil {
		return nil, err
	}
	return place, nil
}
