package service

import (
	"app-server/global"
	"app-server/model"
	"app-server/utils"
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
)

//防止对应端口启动的协程过多
var goRouteOpenCount map[string]int

func SocketServer() {

	fmt.Println("SocketServer")
	for {
		select {
		case work := <-global.SocketServerOpen:
			if work != true {
				break
			}

			fmt.Println("socket监听池开启")
			goRouteOpenCount = make(map[string]int)

			//开启对应监听服务
			for i := range global.SocketAddr {
				//变量覆盖
				i := i
				if goRouteOpenCount[getPort(global.SocketAddr[i])] == 0 && !global.GoRouteOpen[getPort(global.SocketAddr[i])] { //对应端口服务未开启
					goRouteOpenCount[getPort(global.SocketAddr[i])]++
					fmt.Println(getPort(global.SocketAddr[i]), "port")
					go func() {
						fmt.Println("GoRouteOpen")
						//持续监听
						for {
							select {
							case isOpen := <-global.SocketServerChan[getPort(global.SocketAddr[i])]:
								if isOpen { //是否开启
									fu := global.SocketRoute[getPort(global.SocketAddr[i])]
									go fu(global.SocketAddr[i])
								} else {

								}
							}
						}
					}()
				}

			}

		}
	}

}

// Socket socket服务
func Socket(address string) {
	fmt.Println(address, "Socket")
	for {
		fmt.Println("running")
		select {
		case job := <-global.SocketChan[getPort(address)]:
			if job == true {
				fmt.Println("socket open")
				if global.GoRouteOpen[getPort(address)] { //对应端口服务已开启
					break
				}
				global.GoRouteOpen[getPort(address)] = true

				go func() {
					//开启服务
					listen, err := net.Listen("tcp", address) //代表监听的地址端口
					global.SocketListen[getPort(address)] = listen
					if err != nil {
						fmt.Println("listen failed, err:", err)
						return
					}
					fmt.Println("正在等待建立连接.....", listen.Addr())

					for {
						conn, err := listen.Accept() //请求建立连接，客户端未连接就会在这里一直等待
						if err != nil {
							fmt.Println("accept failed, err:", err)
							return
						}
						fmt.Println(conn.LocalAddr(), "连接建立成功.....")
						go process(conn)
					}
				}()

				//更新设备状态
				updateDeviceStatus(getPort(address), true)

			} else {
				if global.GoRouteOpen[getPort(address)] {
					fmt.Println("exit")
					global.GoRouteOpen[getPort(address)] = false
					goRouteOpenCount[getPort(address)]--
					//结束指定socket端口监听
					if global.SocketListen[getPort(address)] != nil {
						_ = global.SocketListen[getPort(address)].Close()
					}
					//更新设备状态
					updateDeviceStatus(getPort(address), false)
					return
				}
			}

		}
	}

}

//连接处理函数
func process(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}(conn)
	for {
		var buf [2048]byte
		n, err := conn.Read(buf[:]) //定义为切片 相当于buf[0:len(buf)]
		if err != nil {             //一直在读取,读取失败break
			log.Println("read from client failed, err:", err)
			break
		}
		log.Println(conn.LocalAddr(), "收到", conn.RemoteAddr(), "发来的数据")
		go storeData(buf[:], n, conn.LocalAddr().String())
	}
}

//数据存储
func storeData(buf []byte, n int, port string) {

	var payload []string //原始16进制数据串
	for i := 0; i < n; i++ {
		//转换为16进制
		payload = append(payload, strconv.FormatInt(int64(buf[i]), 16))

	}

	var device model.Device
	for i := range payload {
		if i == 0 {
			device.Data += payload[i]
		} else {
			device.Data += " " + payload[i]
		}
	}

	device.CreateTime = utils.TimeFormat(time.Now())
	device.Port = port[len(port)-5:]

	_, err := global.DeviceDataColl.InsertOne(context.TODO(), device)
	if err != nil {
		fmt.Println(err)
		return
	}

}
