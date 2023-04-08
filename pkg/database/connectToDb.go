package database

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {

	var err error

	dsn := os.Getenv("DATABASE")

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {

		panic("failed to connect to database")
	}
}
