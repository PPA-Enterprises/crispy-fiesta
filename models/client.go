package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Client struct {
	ID         primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	name       string               `json:"name" bson:"name"`
	InProgress []primitive.ObjectID `json:"inProgress" bson:"in_progress"`
	Completed  []primitive.ObjectID `json:"completed" bson:"in_progress"`
}

type ClientModel struct{}

func (getClient *ClientModel) GetClientById(id string) (Client, error) {
	collection := dbConnect.Use(databaseName, "client")
	var client Client
	err := collection.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&client)

	return client, err
}

//func ClientByName(name string)

/*func FromJob(job *Job) *Client {

}*/
