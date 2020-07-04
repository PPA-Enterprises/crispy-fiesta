package models

import (
	"context"
	"log"

	"github.com/PPA-Enterprises/crispy-fiesta/forms"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Client struct {
	ID         primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string               `json:"name" bson:"name"`
	InProgress []primitive.ObjectID `json:"inProgress" bson:"in_progress"`
	Completed  []primitive.ObjectID `json:"completed" bson:"in_completed"`
}

type ClientModel struct{}

func (createClient *ClientModel) CreateClient(data forms.CreateClientCommand) (*mongo.InsertOneResult, error) {
	collection := dbConnect.Use(databaseName, "client")

	result, err := collection.InsertOne(context.Background(), bson.D{
		{"name", data.Name},
		{"inProgress", data.InProgress},
		{"completed", data.Completed},
	})

	return result, err
}

func (getClient *ClientModel) GetClientById(id string) (Client, error) {
	collection := dbConnect.Use(databaseName, "client")
	var client Client
	oid, _ := primitive.ObjectIDFromHex(id)
	err := collection.FindOne(context.TODO(), bson.D{{"_id", oid}}).Decode(&client)

	return client, err
}

func (getAllClients *ClientModel) GetAllClients() ([]Client, error) {
	collection := dbConnect.Use(databaseName, "client")

	cursor, err := collection.Find(context.TODO(), bson.D{})

	var results []Client

	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Fatal(err)
	}

	return results, err

}

/*func FromJob(job *Job) *Client {

}*/
