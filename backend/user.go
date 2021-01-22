package PPA

import(
"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name,omitempty"`
	Email string `json:"email" bson:"email,omitempty"`
	Password string `json:"-" bson:"password,omitempty"`
}

type AuthUser struct {

}
