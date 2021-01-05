package platform_exercise

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init() *gorm.DB {
	postgresURL := os.Getenv("postgresURL")
	log.Printf("\npostgresURL: %v\n", postgresURL)
	db, err := gorm.Open(postgres.Open(postgresURL), &gorm.Config{})

	if err != nil {
		log.Printf("\nCould not connect to postgres\n%v\n", err)
		panic("Could not connect to postgres")
	}

	return db
}
