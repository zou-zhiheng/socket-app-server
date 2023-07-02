package router

import (
	"app-server/controller"
	"github.com/gin-gonic/gin"
)

func UserRouter(engine *gin.Engine) {
	user := engine.Group("user")
	{
		user.GET("/get", controller.GetUserData)
		user.POST("/create", controller.CreateUser)
		user.PUT("/update", controller.UpdateUser)
		user.DELETE("/delete", controller.DeleteUser)
	}
}
