package initialize

import (
	"app-server/global"
	"app-server/model"
	"app-server/service"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"net"
)

func SocketInit() {
	//在socket表中载入需要监听的地址
	res, err := global.DeviceColl.Find(context.TODO(), bson.M{})
	if err != nil {
		fmt.Println(err)
		return
	}

	var dataDB []model.Device
	if err = res.All(context.TODO(), &dataDB); err != nil {
		fmt.Println(err)
		return
	}

	//全局变量初始化
	global.GoRouteOpen = make(map[string]bool)
	global.SocketChan = make(map[string]chan bool)
	global.SocketListen = make(map[string]net.Listener)
	global.SocketRoute = make(map[string]func(address string)) //初始化socket任务队列
	global.SocketServerOpen = make(chan bool, 1)
	global.SocketServerChan = make(map[string]chan bool)

	//初始化socket配置
	for i := range dataDB {
		global.SocketAddr = append(global.SocketAddr, dataDB[i].Addr)
		if dataDB[i].Port == "" {
			continue
		}
		//初始化管道
		global.SocketChan[dataDB[i].Port] = make(chan bool, 1)
		global.SocketServerChan[dataDB[i].Port] = make(chan bool, 1)
		//添加socket任务,由SocketChan控制服务是否开启
		global.SocketRoute[dataDB[i].Port] = func(address string) {
			service.Socket(address)
		}
	}
	global.SocketServerOpen <- true

	//初始化设备状态
	//更新各个设备的情况
	_, err = global.DeviceColl.UpdateMany(context.TODO(), bson.M{}, bson.M{"$set": bson.M{"isListen": false}})
	if err != nil {
		return
	}

}
