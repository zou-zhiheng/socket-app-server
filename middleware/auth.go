package middleware

import (
	"app-server/global"
	"app-server/model"
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		//判断当前访问的api是否存在

		var apiDB model.Api
		if err := global.ApiColl.FindOne(context.TODO(), bson.M{"method": c.Request.Method, "path": c.Request.URL.Path}).Decode(&apiDB); err != nil {
			//此api不存在
			c.JSON(http.StatusInternalServerError, "此api不存在")
			c.Abort()
			return
		}

		userInter, _ := c.Get("user")
		user := userInter.(model.User)

		var userDB model.User
		if err := global.UserColl.FindOne(context.TODO(), bson.M{"_id": user.Id}).Decode(&userDB); err != nil {
			//用户不存在
			c.JSON(http.StatusInternalServerError, "此用户不存在")
			c.Abort()
			return
		}


		//判断用户是否拥有此权限
		for _, api := range userDB.Auth {
			if api == apiDB.Id {
				c.Next()
				return
			}
		}

		//终止响应
		c.JSON(http.StatusInternalServerError, "此用户无此接口访问权限")
		c.Abort()
		return

	}
}
