package PPA

import (
	"bytes"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Tinter struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" m:"Database ID"`
	Name string `json:"name" bson:"name,omitempty" m:"Name"`
	Phone string `json:"phone" bson:"phone,omitempty" m:"Phone Number"`
	Jobs []primitive.ObjectID `json:"jobs" bson:"jobs,omitempty"`
	History []LogEvent `json:"history" bson:"history,omitempty"`
}


func (c *Tinter) AttatchJob(jobOid primitive.ObjectID) {
	const matched int = 0
	var exists = false

	for _, oid := range c.Jobs {
		result := bytes.Compare([]byte(jobOid.Hex()), []byte(oid.Hex()))
		if result == matched {
			exists = true
			break
		}
	}
	if !exists {
		c.Jobs = append(c.Jobs, jobOid)
	}
}

func (c *Tinter) FindAndRemoveJob(jobOid primitive.ObjectID) {
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

func (c *Tinter) AppendLog(event LogEvent) {
	c.History = append(c.History, event)
}
