package PPA

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Client struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name,omitempty"`
	Phone string `json:"phone" bson:"phone,omitempty"`
	Jobs []primitive.ObjectID `json:"jobs" bson:"jobs,omitempty"`
}

func (c *Client) AttatchJob(jobOid primitive.ObjectID) {
	c.Jobs = append(c.Jobs, jobOid)
}
