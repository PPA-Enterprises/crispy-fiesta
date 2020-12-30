package types

import ("go.mongodb.org/mongo-driver/bson/primitive")

type DeliverableUser struct {
	ID primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name"`
}
