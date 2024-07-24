package db

import (
	"fmt"
	"log"
	"os"

	"gomap/internal/gps_utils"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

type City struct {
	Name string `gorm:"primaryKey;uniqueIndex"`
}

type User struct {
	ID         string `gorm:"primaryKey"`
	Imei       string
	TelegramId string
	Email      string
	Phone      string
}

type UserPlace struct {
	ID   uint `gorm:"primaryKey"`
	Info string
	Geom string `gorm:"type:geography(Point,4326)"`
}

type Place struct {
	ID       uint `gorm:"primaryKey"`
	CityName string
	City     City `gorm:"foreignKey:CityName;references:Name;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Name     string
	Geom     string `gorm:"type:geography(Point,4326)"`
	Desc     string
	Distance float64 `gorm:"-"`
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

	var openErr error
	db, openErr = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if openErr != nil {
		log.Fatal(openErr)
	}

	// Migrate the schema
	if err := createSchema(schema); err != nil {
		log.Fatal("Error creating schema:", err)
	}

	if err := createPostGIS(); err != nil {
		log.Fatal("Error creating PostGIS extension:", err)
	}

	// Auto-migrate the tables
	if err := db.AutoMigrate(&City{}, &User{}, &UserPlace{}, &Place{}); err != nil {
		log.Fatal("Error migrating tables:", err)
	}
}

func createSchema(schema string) error {
	schemaQuery := fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS %s;`, schema)
	return db.Exec(schemaQuery).Error
}

func createPostGIS() error {
	return db.Exec("CREATE EXTENSION IF NOT EXISTS postgis;").Error
}

func CreateUser(user *User) (*User, error) {
	user.ID = uuid.New().String()
	result := db.Create(user)
	return user, result.Error
}

func CreateUserPlace(place *UserPlace) (*UserPlace, error) {
	result := db.Create(place)
	return place, result.Error
}

func CreatePlace(place *Place) (*Place, error) {
	result := db.Create(place)
	return place, result.Error
}

func GetNearPlaces(myPoint gps_utils.GpsCoordinates) ([]Place, error) {
	lat := myPoint.Lat
	lng := myPoint.Lng

	query := fmt.Sprintf(`
		SELECT id, city_name, name, ST_AsText(geom) as geom, description, ST_Distance(geom::geography, ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography) AS distance
		FROM %s.places
		ORDER BY distance ASC`, os.Getenv("DB_SCHEMA"))

	rows, err := db.Raw(query, lng, lat).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var places []Place
	for rows.Next() {
		var place Place
		err := rows.Scan(&place.ID, &place.CityName, &place.Name, &place.Geom, &place.Desc, &place.Distance)
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
