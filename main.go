package main

import (
	"app-server/initialize"
	"app-server/router"
	"app-server/service"
	"fmt"
)

// 项目运行初始配置
func init() {
	initialize.Init()
}

func main() {
	fmt.Println("test")
	//开启socket服务
	go service.SocketServer()
	engine := router.GetEngine()
	if err := engine.Run(":7001"); err != nil {
		panic(err)
	}
}
