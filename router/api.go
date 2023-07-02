package router

import (
	"app-server/controller"
	"github.com/gin-gonic/gin"
)

func ApiRouter(engine *gin.Engine) {
	api := engine.Group("api")
	{
		api.GET("/get", controller.GetApi)
		api.POST("/create", controller.CreateApi)
		api.PUT("/update", controller.UpdateApi)
		api.DELETE("/delete", controller.DeleteApi)
	}
}
