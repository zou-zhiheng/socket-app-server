package service

import (
	"app-server/global"
	"app-server/model"
	"app-server/utils"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"time"
)

// GetDeviceEchartsPie  饼图
func GetDeviceEchartsPie() utils.Response {

	res, err := global.DeviceColl.Find(context.TODO(), bson.M{})
	if err != nil {
		return utils.ErrorMess("失败", err.Error())
	}

	var dataDB []model.Device
	if err = res.All(context.TODO(), &dataDB); err != nil {
		return utils.ErrorMess("失败", err.Error())
	}

	type InfoDetail struct {
		Name  string `json:"name" bson:"name"`
		Value int    `json:"value" bson:"value"`
	}

	type Data struct {
		Name []string     `json:"name" bson:"name"`
		Info []InfoDetail `json:"info" bson:"info"`
	}

	deviceMap := make(map[string]int)
	//数据处理
	for i := range dataDB {
		if dataDB[i].Port != "" {
			deviceMap[dataDB[i].Port]++
		}
	}

	var data Data
	for key := range deviceMap {
		//name
		data.Name = append(data.Name, key)
		//detail
		data.Info = append(data.Info, InfoDetail{
			Name:  key,
			Value: deviceMap[key],
		})
	}

	return utils.SuccessMess("成功", data)

}

func updateDeviceStatus(port string, isListen bool) {
	//更新各个设备的情况
	_, _ = global.DeviceColl.UpdateMany(context.TODO(), bson.M{"port": port}, bson.M{"$set": bson.M{"isListen": isListen}})
}

func RecDeviceSignal(port, signal string) utils.Response {

	//判断指定端口服务是否存在]
	var flag bool
	for key := range global.SocketRoute {
		if port == key {
			flag = true
		}
	}

	if !flag {
		return utils.ErrorMess("失败，该端口服务不存在", nil)
	}

	if signal != "true" && signal != "false" {
		return utils.ErrorMess("失败,参数错误", nil)
	} else {
		global.SocketServerChan[port] <- true
		if signal == "true" {
			if !global.GoRouteOpen[port] {
				global.SocketChan[port] <- true
			}
		} else {

			//对应服务已开启
			if global.GoRouteOpen[port] {
				//关闭对应socket服务
				global.SocketChan[port] <- false
			}

		}
	}

	return utils.SuccessMess("成功", nil)
}

func CreateDevice(device model.Device) utils.Response {

	if device.Code == "" {
		return utils.ErrorMess("失败，设备编号不能为空", nil)
	}

	fmt.Println(device)

	err := global.DeviceColl.FindOne(context.TODO(), bson.M{"code": device.Code}).Decode(&bson.M{})
	if err == nil {
		return utils.ErrorMess("失败，该设备已存在", err)
	}

	if len(device.Addr) >= 12 {
		device.Port = device.Addr[len(device.Addr)-4:]
		device.IP, device.Port = strSplitOrder(device.Addr, ":")

		flag := true
		for i := range global.SocketAddr {
			if global.SocketAddr[i] == device.Addr {
				flag = false
			}
		}

		if flag {
			//更新配置
			//初始化管道
			global.SocketChan[device.Port] = make(chan bool, 1)
			global.SocketServerChan[device.Port] = make(chan bool, 1)
			//添加socket任务,由SocketChan控制服务是否开启
			global.SocketRoute[device.Port] = func(address string) {
				Socket(address)
			}
			//更新socket监听池
			global.SocketAddr = append(global.SocketAddr, device.Addr)
			//更新协议池配置
			global.SocketServerOpen <- true
		}

	}

	device.CreateTime = utils.TimeFormat(time.Now())
	res, err := global.DeviceColl.InsertOne(context.TODO(), device)
	if err != nil {
		fmt.Println(err)
		return utils.ErrorMess("失败", err.Error())
	}
	fmt.Println(res.InsertedID)

	return utils.SuccessMess("成功", res.InsertedID)
}

func GetDevice(code, flag string, currPage, pageSize, startTime, endTime string) utils.Response {

	var coll *mongo.Collection
	switch flag {
	case "true":
		coll = global.DeviceDataColl
	case "false":
		coll = global.DeviceColl
	default:
		return utils.ErrorMess("失败,参数错误", nil)
	}

	size, err := strconv.ParseInt(pageSize, 10, 64)
	if err != nil {
		return utils.ErrorMess("行数字段过长", size)
	}
	curr, err := strconv.ParseInt(currPage, 10, 64)
	if err != nil {
		return utils.ErrorMess("指定页面字段过长", nil)
	}

	skip := (curr - 1) * size
	opt := options.FindOptions{
		Limit: &size,
		Skip:  &skip,
		Sort:  bson.M{"_id": -1},
	}

	filter := bson.M{}

	if code != "" {
		filter["code"] = code
	}

	if startTime != "" && endTime != "" {
		filter = bson.M{
			"createTime": bson.M{
				"$gte": startTime,
				"$lte": endTime,
			},
		}
	}

	count, _ := coll.CountDocuments(context.TODO(), bson.M{})
	res, err := coll.Find(context.TODO(), bson.M{}, &opt)
	if err != nil {
		return utils.ErrorMess("失败", err.Error())
	}

	var deviceDB []model.Device
	if err = res.All(context.TODO(), &deviceDB); err != nil {
		return utils.ErrorMess("失败", err.Error())
	}

	for i := range deviceDB {
		deviceDB[i].CreateTime = deviceDB[i].CreateTime[len(deviceDB[i].CreateTime)-8:]
	}

	return utils.SuccessMess("成功", struct {
		Count int64          `json:"count" bson:"count"`
		Data  []model.Device `json:"data" bson:"data"`
	}{
		Count: count,
		Data:  deviceDB,
	})
}

func UpdateDevice(device model.Device) utils.Response {

	if device.Id.IsZero() {
		return utils.ErrorMess("失败,禁止访问!", nil)
	}

	var deviceDB model.Device
	err := global.DeviceColl.FindOne(context.TODO(), bson.M{"_id": device.Id}).Decode(&deviceDB)
	if err != nil {
		return utils.ErrorMess("失败", err)
	}

	if len(device.Addr) >= 12 {
		device.Port = device.Addr[len(device.Addr)-4:]
		device.IP, device.Port = strSplitOrder(device.Addr, ":")

		flag := true
		for i := range global.SocketAddr {
			if global.SocketAddr[i] == device.Addr {
				flag = false
			}
		}

		if flag {
			//更新配置
			//初始化管道
			global.SocketChan[device.Port] = make(chan bool, 1)
			global.SocketServerChan[device.Port] = make(chan bool, 1)
			//添加socket任务,由SocketChan控制服务是否开启
			global.SocketRoute[device.Port] = func(address string) {
				Socket(address)
			}
			//更新socket监听池
			global.SocketAddr = append(global.SocketAddr, device.Addr)
			//更新协议池配置
			global.SocketServerOpen <- true
		}

	}
	device.CreateTime = deviceDB.CreateTime
	device.UpdateTime = utils.TimeFormat(time.Now())
	device.IsListen = deviceDB.IsListen

	res, err := global.DeviceColl.UpdateOne(context.TODO(), bson.M{"_id": device.Id}, bson.M{"$set": device})
	if err != nil {
		return utils.ErrorMess("失败", err)
	}

	return utils.SuccessMess("成功", res)
}

func DeleteDevice(idStr string) utils.Response {
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return utils.ErrorMess("失败", err.Error())
	}

	res, err := global.DeviceColl.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		return utils.ErrorMess("失败", err.Error())
	}

	return utils.SuccessMess("成功", res.DeletedCount)

}

//根据第一个指定字符分割字符串并舍去指定字符，可自定义修改
func strSplitOrder(payload string, index string) (string, string) {

	var l, r string
	if payload == "" || index == "" {
		return l, r
	}

	index = index[:]
	for i := range payload {
		//找到指定位置
		if payload[i] == index[0] {
			l = payload[0:i]
			r = payload[i+1:]
		}
	}

	return l, r
}

func getPort(address string) string {

	_, r := strSplitOrder(address, ":")

	return r
}
