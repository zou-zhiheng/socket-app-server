package router

import (
	"app-server/controller"
	"github.com/gin-gonic/gin"
)

func DeviceRouter(engine *gin.Engine) {
	device := engine.Group("device")
	{
		device.GET("/getPie", controller.GetDeviceEchartsPie)
		device.GET("/get", controller.GetDevice)
		device.POST("/create", controller.CreateDevice)
		device.PUT("/update", controller.UpdateDevice)
		device.DELETE("/delete", controller.DeleteDevice)
		device.GET("/signal", controller.RecDeviceSignal)
	}
}
