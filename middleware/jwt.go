package middleware

import (
	"app-server/model"
	"app-server/utils"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type CustomClaims struct {
	User model.User `json:"user"`
	jwt.StandardClaims
}

var SigningKey = []byte("zouzhiheng@app")

func CreateToken(user model.User) (string, error) {
	//获取token，前两部分
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{User: user,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(), //签名生效时间
			//ExpiresAt: time.Now().Unix() + 60*60*24*2, //2小时过期
			Issuer:    "zzh", //签发人，
		},
	})
	//根据密钥生成加密token，token完整三部分
	tokenString, err := token.SignedString(SigningKey)
	if err != nil {
		return "", err
	}
	return tokenString, err

}

func ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return SigningKey, nil
	})
	if err != nil {
		return nil, err
	}
	if token != nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, errors.New("token无效")
	}
	return nil, errors.New("token无效")
}

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取token
		token := c.Request.Header.Get("token")
		if token == "" {
			c.JSON(http.StatusOK, utils.Response{Code: 401, Message: "token缺失", Data: ""})
			//终止
			c.Abort()
			return
		}
		claims, err := ParseToken(token)
		if err != nil {
			c.JSON(http.StatusOK, utils.Response{Code: 401, Message: "token过期", Data: ""})
			//终止
			c.Abort()
			return
		}
		//将用户信息储存再上下文
		c.Set("user", claims.User)
		//继续下面的操作
		c.Next()
	}
}
