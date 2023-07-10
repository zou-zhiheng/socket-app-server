package main

import (
	"app-server/initialize"
	"app-server/router"
	"app-server/service"
	"fmt"
)

func init() {
	initialize.Init()
}

func main() {
	//开启socket服务
	fmt.Println("test")
	go service.SocketServer()
	engine := router.GetEngine()
	if err := engine.Run(":7001"); err != nil {
		panic(err)
	}
}
