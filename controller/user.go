package controller

import (
	"app-server/model"
	"app-server/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUserData(c *gin.Context) {
	name := c.Query("name")
	currPage := c.DefaultQuery("currPage", "1")
	pageSize := c.DefaultQuery("pageSize", "0")
	startTime := c.Query("startTime")
	endTime := c.Query("endTime")

	c.JSON(http.StatusOK, service.GetUserData(name, currPage, pageSize, startTime, endTime))
}

func CreateUser(c *gin.Context) {

	var user model.User
	if err := c.Bind(&user); err != nil {
		c.JSON(http.StatusInternalServerError, "参数错误")
		return
	}

	c.JSON(http.StatusOK, service.CreateUser(user))

}

func UpdateUser(c *gin.Context) {

	flag := c.DefaultQuery("flag", "false")
	if flag != "false" && flag != "true" {
		c.JSON(http.StatusInternalServerError, "参数错误")
		return
	}

	var user model.User
	if err := c.Bind(&user); err != nil {
		c.JSON(http.StatusInternalServerError, "参数错误")
		return
	}

	c.JSON(http.StatusOK, service.UpdateUser(user,flag))

}

func DeleteUser(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusInternalServerError, "参数错误")
		return
	}

	c.JSON(http.StatusOK, service.DeleteUser(id))
}

func UserLogin(c *gin.Context) {
	var user model.User
	if err := c.Bind(&user); err != nil {
		c.JSON(http.StatusInternalServerError, "参数错误")
		return
	}

	c.JSON(http.StatusOK, service.UserLogin(user))
}
