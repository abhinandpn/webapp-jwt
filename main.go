package main

import (
	"jwt/pkg/database"
	"jwt/pkg/initializer"
	"jwt/pkg/route"

	"github.com/gin-gonic/gin"
)

func init() {

	initializer.LoadEnvVariable()
	database.ConnectToDB()
	database.UserDB()
	database.AdminDB()

}
func main() {

	router := gin.Default()

	router.Use(gin.Logger())

	route.UserRoute(router)
	route.AdminRoute(router)

	router.Run()
}
