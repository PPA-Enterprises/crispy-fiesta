package PPA

import (
	"bytes"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Client struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" m:"Database ID"`
	Name string `json:"name" bson:"name,omitempty" m:"Name"`
	Phone string `json:"phone" bson:"phone,omitempty" m:"Phone Number"`
	Jobs []primitive.ObjectID `json:"jobs" bson:"jobs,omitempty"`
	History []LogEvent `json:"history" bson:"history,omitempty"`
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

func (c *Client) AppendLog(event LogEvent) {
	c.History = append(c.History, event)
}

type ClientLabel struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" m:"Database ID"`
	LabelName string `json:"label_name" bson:"label_name,omitempty" m:"Label Name"`
	IsDeleted bool `json:"-" bson:"is_deleted"`
	Clients []primitive.ObjectID `json:"clients" bson:"clients,omitempty"`
	History []LogEvent `json:"history" bson:"history,omitempty"`
}

func(l *ClientLabel) AppendLog(event LogEvent) {
	l.History = append(l.History, event)
}
