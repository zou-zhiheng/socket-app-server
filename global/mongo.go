package global

import "go.mongodb.org/mongo-driver/mongo"

var (
	MongoClient    *mongo.Client
	UserColl       *mongo.Collection
	ApiColl        *mongo.Collection
	DeviceColl     *mongo.Collection
	DeviceDataColl *mongo.Collection
)
