package controller

import (
	"app-server/model"
	"app-server/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetApi(c *gin.Context) {
	currPage := c.DefaultQuery("currPage", "1")
	pageSize := c.DefaultQuery("pageSize", "10")
	startTime := c.Query("startTime")
	name := c.Query("name")
	endTime := c.Query("endTime")

	c.JSON(http.StatusOK, service.GetApi(name, currPage, pageSize, startTime, endTime))
}

func CreateApi(c *gin.Context) {

	var api model.Api
	if err := c.Bind(&api); err != nil {
		c.JSON(http.StatusInternalServerError, "参数错误")
		return
	}

	c.JSON(http.StatusOK, service.CreateApi(api))

}

func UpdateApi(c *gin.Context) {
	var api model.Api
	if err := c.Bind(&api); err != nil {
		c.JSON(http.StatusInternalServerError, "参数错误")
		return
	}

	c.JSON(http.StatusOK, service.UpdateApi(api))
}

func DeleteApi(c *gin.Context) {

	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusInternalServerError, "参数错误")
		return
	}

	c.JSON(http.StatusOK, service.DeleteApi(id))
}
