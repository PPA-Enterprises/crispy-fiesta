package PPA

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Client struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name"`
	Phone string `json:"phone" bson:"phone"`
	Jobs []primitive.ObjectID `json:"-" bson:"jobs"`
}
