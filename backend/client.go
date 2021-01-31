package PPA

import (
	"bytes"
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

func (c *Client) FindAndRemoveJob(jobOid primitive.ObjectID) {
	const matched int = 0
	for i, oid := range c.Jobs {
		result := bytes.Compare([]byte(jobOid.Hex()), []byte(oid.Hex()))
		if result == matched {
			// preserve the order. Idiomatic way
			c.Jobs = append(c.Jobs[:i], c.Jobs[i+1:]...)
			return
		}
	}

}
