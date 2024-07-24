package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/qor/admin"
)

type Place struct {
	ID          uint   `gorm:"primary_key"`
	Name        string `gorm:"type:varchar(100)"`
	Description string `gorm:"type:text"`
	Geom        string `gorm:"type:geography(Point,4326)"`
}

func main() {
	// Define the PostgreSQL connection string
	dsn := "host=localhost user=postgres dbname=gomap sslmode=disable password=159753 port=5433 search_path=vetrof"
	DB, err := gorm.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer DB.Close()

	// Ensure the 'places' table exists and is correctly migrated
	if err := DB.AutoMigrate(&Place{}).Error; err != nil {
		log.Fatalf("failed to migrate database schemas: %v", err)
	}

	Admin := admin.New(&admin.AdminConfig{DB: DB})

	Admin.AddResource(&Place{})

	mux := http.NewServeMux()

	Admin.MountTo("/admin", mux)

	fmt.Println("Listening on: 9595")
	if err := http.ListenAndServe(":9595", mux); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
