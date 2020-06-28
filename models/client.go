package models

import (
//	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"

)

type Client struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
	name string `bson:"name"`
	InProgress []primitive.ObjectID `bson:"in_progress"`
	Completed []primitive.ObjectID `bson:"in_progress"`
}

func FromJob(job *Job) *Client {

}
