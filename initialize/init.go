package initialize

import "fmt"

func Init() {
	fmt.Println("init")
	MongoInit()
	SocketInit()
}
