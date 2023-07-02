package global

import "net"

var (
	SocketAddr       []string
	SocketChan       map[string]chan bool            //控制socket服务
	SocketListen     map[string]net.Listener         //socket句柄
	GoRouteOpen      map[string]bool                 //对应对口socket服务是否开启
	SocketRoute      map[string]func(address string) //socket服务任务队列
	SocketServerOpen chan bool                       //开启socket服务池
	SocketServerChan map[string]chan bool            //控制开启指定socket服务
)
