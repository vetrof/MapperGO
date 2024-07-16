package db

import (
	"database/sql"
	"fmt"
	"gomap/internal/gps_utils"
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
	Info string
	Geom string
}

type Place struct {
	ID       uint
	Name     string
	Geom     string
	Distance float64
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
		log.Println("Connected to database OK")
	}

	// Создание схемы и расширения PostGIS
	createSchema(schema)
	createPostGIS()

	// Создание таблиц, если они не существуют
	createTables(schema)
}

func createSchema(schema string) {
	schemaQuery := fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS %s;`, schema)
	_, err := db.Exec(schemaQuery)
	if err != nil {
		log.Fatal("Error creating schema:", err)
	}
}

func createPostGIS() {
	_, err := db.Exec("CREATE EXTENSION IF NOT EXISTS postgis;")
	if err != nil {
		log.Fatal("Error creating PostGIS extension:", err)
	}
}

func createTables(schema string) {
	userTableQuery := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s.users (
		id UUID PRIMARY KEY,
		imei TEXT,
		telegram_id TEXT,
		email TEXT,
		phone TEXT
	);`, schema)

	userPlaceTableQuery := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s.user_places (
		id SERIAL PRIMARY KEY,
		info TEXT,
		geom GEOGRAPHY(Point, 4326)
	);`, schema)

	placeTableQuery := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s.places (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100),
		geom GEOGRAPHY(Point, 4326)
	);`, schema)

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
	query := fmt.Sprintf(`INSERT INTO %s.users (id, imei, telegram_id, email, phone) VALUES ($1, $2, $3, $4, $5)`, os.Getenv("DB_SCHEMA"))
	_, err := db.Exec(query, user.ID, user.Imei, user.TelegramId, user.Email, user.Phone)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func CreateUserPlace(place *UserPlace) (*UserPlace, error) {
	query := fmt.Sprintf(`INSERT INTO %s.user_places (info, geom) VALUES ($1, ST_GeogFromText($2)) RETURNING id`, os.Getenv("DB_SCHEMA"))
	err := db.QueryRow(query, place.Info, place.Geom).Scan(&place.ID)
	if err != nil {
		return nil, err
	}
	return place, nil
}

func CreatePlace(place *Place) (*Place, error) {
	query := fmt.Sprintf(`INSERT INTO %s.places (name, geom) VALUES ($1, ST_GeogFromText($2)) RETURNING id`, os.Getenv("DB_SCHEMA"))
	err := db.QueryRow(query, place.Name, place.Geom).Scan(&place.ID)
	if err != nil {
		return nil, err
	}
	return place, nil
}

func GetNearPlaces(myPoint gps_utils.GpsCoordinates) ([]Place, error) {
	lat := myPoint.Lat
	lng := myPoint.Lng
	fmt.Println("GetNearPlaces --->>> ", lat, lng)

	//query := fmt.Sprintf(`
	//	SELECT id, name, ST_AsText(geom), ST_Distance(geom, ST_SetSRID(ST_MakePoint($1, $2), 4326)) AS distance
	//	FROM %s.places
	//	ORDER BY distance ASC`, os.Getenv("DB_SCHEMA"))

	query := fmt.Sprintf(`
		SELECT id, name, ST_AsText(geom), ST_Distance(geom::geography, ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography) AS distance
		FROM %s.places
		ORDER BY distance ASC`, os.Getenv("DB_SCHEMA"))

	rows, err := db.Query(query, lng, lat) // координаты подставляются как параметры

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var places []Place
	for rows.Next() {
		var place Place
		err := rows.Scan(&place.ID, &place.Name, &place.Geom, &place.Distance)
		if err != nil {
			return nil, err
		}
		places = append(places, place)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return places, nil
}

//func GetNearPlaces(myPoint gps_utils.GpsCoordinates) ([]Place, error) {
//	lat := myPoint.Lat
//	lng := myPoint.Lng
//	fmt.Println("GetNearPlaces --->>> ", lat, lng)
//
//	//query := fmt.Sprintf(`
//	//	SELECT id, name, ST_AsText(geom)
//	//	FROM %s.places
//	//	ORDER BY ST_Distance(geom, ST_SetSRID(ST_MakePoint($1, $2), 4326)) ASC`, os.Getenv("DB_SCHEMA"))
//
//	query := fmt.Sprintf(`
//		SELECT id, name, ST_AsText(geom), ST_Distance(geom, ST_SetSRID(ST_MakePoint($1, $2), 4326)) AS distance
//		FROM %s.places
//		ORDER BY distance ASC`, os.Getenv("DB_SCHEMA"))
//
//	rows, err := db.Query(query, lng, lat) // координаты подставляются как параметры
//
//	if err != nil {
//		return nil, err
//	}
//	defer rows.Close()
//
//	var places []Place
//	for rows.Next() {
//		var place Place
//		err := rows.Scan(&place.ID, &place.Name, &place.Geom)
//		if err != nil {
//			return nil, err
//		}
//		places = append(places, place)
//	}

//	if err = rows.Err(); err != nil {
//		return nil, err
//	}
//
//	return places, nil
//}
