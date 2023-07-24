package main

import (
	"app-server/initialize"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// 项目运行初始配置
func init() {
	initialize.Init()
}

func main() {
	//开启socket服务
	//go service.SocketServer()
	//engine := router.GetEngine()
	//if err := engine.Run(":7001"); err != nil {
	//	panic(err)
	//}

	token := CreateToken()
	fmt.Println(token)
}

func CreateToken() string {
	// 创建一个新的token对象，指定签名方法和声明
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"foo": "bar",
		"nbf": time.Now().Unix(),
	})

	// 使用私钥进行签名并获取完整的编码后的字符串token
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte("aaaaaaaaaaa"))
	if err != nil {
		fmt.Println(err, "test")
		return ""
	}
	tokenString, _ := token.SignedString(privateKey)
	return tokenString
}

func ParseToken() {
	// 解析token
	token, err := jwt.Parse("YOUR TOKEN STRING", func(token *jwt.Token) (interface{}, error) {
		// 确保token方法符合"SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// 返回公钥，用于验证签名
		return jwt.ParseRSAPublicKeyFromPEM([]byte("YOUR PUBLIC KEY"))
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["foo"], claims["nbf"])
	} else {
		fmt.Println(err)
	}
}
