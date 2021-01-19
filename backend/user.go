package PPA

import(
"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id"`
	Name string `json:"name" bson:"name"`
	Email string `json:"email" bson:"email"`
	Password string `json:"-" bson:"password"`
	IsDeleted bool `json:"is_deleted" bson:"is_deleted"`
}

type AuthUser struct {

}
