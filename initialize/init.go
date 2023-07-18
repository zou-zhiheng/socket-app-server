package initialize

import "fmt"

func Init() {
	fmt.Println("coding")
	MongoInit()
	SocketInit()
}
