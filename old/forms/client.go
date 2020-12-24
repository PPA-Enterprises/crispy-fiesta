package forms

import "go.mongodb.org/mongo-driver/bson/primitive"

type CreateClientCommand struct {
	Name       string               `json:"name" binding:"required"`
	InProgress []primitive.ObjectID `json:"inProgress"`
	Completed  []primitive.ObjectID `json:"completed"`
}
