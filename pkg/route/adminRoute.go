package route

import (
	"jwt/pkg/controller"
	"jwt/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRoute(c *gin.Engine) {

	c.POST("/adminlogin", controller.AdminLogin)
	c.POST("/addadmin", controller.AddAdmin)
	c.POST("adminlogout", middleware.AdminAuth, controller.AdminLogout)
	c.GET("/userview", middleware.AdminAuth, controller.UserView)
	c.GET("/deleteuser", middleware.AdminAuth, controller.UserDelete)
	c.PUT("/edituser", middleware.AdminAuth, controller.EditUser)
	c.PATCH("/userblock", middleware.AdminAuth, controller.BlockUser)
}
