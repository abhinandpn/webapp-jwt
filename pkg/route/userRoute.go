package route

import (
	"jwt/pkg/controller"
	"jwt/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoute(c *gin.Engine) {

	c.POST("/usersignup", controller.SignUpUser)
	c.POST("/userlogin", controller.LoginUser)
	c.GET("/userprofilr", middleware.UserAuth, controller.UserProfile)
	c.PUT("/useredit", middleware.UserAuth, controller.UserEdit)
}
