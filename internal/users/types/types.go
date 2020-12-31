package types

import ("go.mongodb.org/mongo-driver/bson/primitive")

type DeliverableUser struct {
	ID primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name"`
	Email string `json:"email" bson:"email"`
	//IsVerified bool `json:"is_verified" bson:"is_verified"`
	//IsDeleted bool `json:"is_deleted" bson:"is_deleted" binding:"required"`
}
