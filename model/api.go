package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Api struct {
	Id     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name   string             `json:"name" bson:"name"`
	Path   string             `json:"path" bson:"path"`
	Method string             `json:"method" bson:"method"`
}
