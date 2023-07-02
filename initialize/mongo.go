package initialize

import (
	"app-server/global"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func MongoInit() {
	if global.MongoClient == nil {
		global.MongoClient = getMongoClient("mongodb://127.0.0.1:27017/?authSource=admin")
	}
	appManager := global.MongoClient.Database("appManager")
	{
		global.UserColl = appManager.Collection("user")
		global.ApiColl = appManager.Collection("api")
		global.DeviceColl = appManager.Collection("device")
		global.DeviceDataColl = appManager.Collection("deviceData")
	}
}

func getMongoClient(uri string) *mongo.Client {
	clientOptions := options.Client().ApplyURI(uri)

	MongoClient, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		fmt.Println(err)
	}
	if err = MongoClient.Ping(context.TODO(), nil); err != nil {
		fmt.Println(err)
	}
	return MongoClient
}
