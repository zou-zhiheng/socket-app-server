package router

import (
	"app-server/controller"
	"app-server/middleware"
	"github.com/gin-gonic/gin"
)

func GetEngine() *gin.Engine {

	engine := gin.Default()
	engine.Use(middleware.Cors())
	engine.POST("/login", controller.UserLogin)

	//权限
	//engine.Use(middleware.JWTAuth(), middleware.Auth())
	//跨域
	//用户
	UserRouter(engine)
	//API
	ApiRouter(engine)
	//设备
	DeviceRouter(engine)

	return engine

}
