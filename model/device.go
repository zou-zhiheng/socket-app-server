package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Device struct {
	Id         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`               //Id
	Code       string             `json:"code" bson:"code"`                                 //编码
	Addr       string             `json:"addr" bson:"addr"`                                 //addr
	Protocol   string             `json:"protocol" bson:"protocol"`                         //协议
	IP         string             `json:"ip" bson:"ip"`                                     //ip地址
	Port       string             `json:"port" bson:"port"`                                 //端口
	IsListen   bool               `json:"isListen" bson:"isListen"`                         //是否监听
	Data       string             `json:"data" bson:"data"`                                 //原始数据,16进制
	CreateTime string             `json:"createTime,omitempty" bson:"createTime,omitempty"` //创建时间
	UpdateTime string             `json:"updateTime,omitempty" bson:"updateTime,omitempty"` //更新时间
}
