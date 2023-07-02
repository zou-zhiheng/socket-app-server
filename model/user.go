package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id         primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string               `json:"name" bson:"name"`
	Account    string               `json:"account" bson:"account"`
	Password   string               `json:"password" bson:"password"`
	Salt       string               `json:"salt" bson:"salt"`
	Auth       []primitive.ObjectID `json:"auth" bson:"auth"`
	IsValid    bool                 `json:"isValid" bson:"isValid"`
	CreateTime string               `json:"createTime,omitempty" bson:"createTime,omitempty"`
	UpdateTime string               `json:"updateTime,omitempty" bson:"updateTime,omitempty"`
}
