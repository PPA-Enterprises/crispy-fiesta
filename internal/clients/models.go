package clients

import (
	"context"
	"bytes"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	errors "internal/common"
	"internal/db"
	"internal/uid"
	"internal/jobs"
)

type Client interface {
	AttatchJobID(primitive.ObjectID) (*errors.ResponseError)
}

type clientModel struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name"`
	Phone string `json:"phone" bson:"phone"`
	Jobs []primitive.ObjectID `json:"jobs" bson:"jobs"`
}

func NewClient(name, phone string) *Client {
	return &clientModel{
		Name: name,
		Phone: phone,
		InProgress: []primitive.ObjectID{},
		Completed: []primitive.ObjectID{},
	}
}

func ClientByPhone(ctx context.Context, phone string) (*Client, error) {
	coll := db.Connection().Use(db.DefaultDatabase, "client")

	var foundClient clientModel
	err := coll.FindOne(ctx, bson.D{{"phone", phone}}).Decode(&foundClient)
	if err != nil {
		return nil, err
	}
	return foundClient, nil
}
/*
func (self *clientModel) AttatchJobID(oid primitive.ObjectID) {
	//search for id, insert if not already in the array
	// linear search for now
	for _, id
}*/
