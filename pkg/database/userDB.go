package database

import "jwt/pkg/models"

func UserDB() {
	DB.AutoMigrate(&models.User{})
}
