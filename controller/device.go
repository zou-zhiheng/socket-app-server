package controller

import (
	"app-server/model"
	"app-server/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetDeviceEchartsPie 饼图
func GetDeviceEchartsPie(c *gin.Context) {
	c.JSON(http.StatusOK, service.GetDeviceEchartsPie())
}

// RecDeviceSignal 接收信号
func RecDeviceSignal(c *gin.Context) {
	port := c.Query("port")
	signal := c.Query("signal")
	if port == "" || signal == "" {
		c.JSON(http.StatusInternalServerError, "端口号或信号不能为空")
		return
	}

	c.JSON(http.StatusOK, service.RecDeviceSignal(port, signal))

}

func CreateDevice(c *gin.Context) {
	var device model.Device
	if err := c.Bind(&device); err != nil {
		c.JSON(http.StatusInternalServerError, "参数错误")
		return
	}


	c.JSON(http.StatusOK, service.CreateDevice(device))
}

func UpdateDevice(c *gin.Context) {
	var device model.Device
	if err := c.Bind(&device); err != nil {
		c.JSON(http.StatusInternalServerError, "参数错误")
		return
	}

	c.JSON(http.StatusOK, service.UpdateDevice(device))
}

func GetDevice(c *gin.Context) {
	code := c.Query("code")
	currPage := c.DefaultQuery("currPage", "1")
	pageSize := c.DefaultQuery("pageSize", "10")
	flag := c.DefaultQuery("flag", "true")
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")

	c.JSON(http.StatusOK, service.GetDevice(code, flag, currPage, pageSize, startTime, endTime))
}

func DeleteDevice(c *gin.Context) {
	id := c.Query("id")

	if id == "" {
		c.JSON(http.StatusInternalServerError, "参数为空")
		return
	}

	c.JSON(http.StatusOK, service.DeleteDevice(id))
}
