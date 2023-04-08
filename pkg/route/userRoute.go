package route

import (
	"jwt/pkg/controller"
	"jwt/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoute(c *gin.Engine) {

	c.POST("/signup", controller.SignUpUser)
	c.POST("/login", controller.LoginUser)
	c.GET("/validate", middleware.UserAuth, controller.Validation)
}
