package platform_exercise

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init(postgresURL string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(postgresURL), &gorm.Config{})

	if err != nil {
		panic("Could not connect to postgres")
	}

	return db
}
