package database

import "jwt/pkg/models"

func AdminDB() {
	DB.AutoMigrate(&models.Admin{})
}
