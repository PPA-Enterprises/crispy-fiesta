package models

import (
	"github.com/PPA-Enterprises/crispy-fiesta/forms"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	ID         bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string        `json:"name" bson:"name"`
	Email      string        `json:"email" bson:"email"`
	Password   string        `json:"password" bson:"password"`
	IsVerified bool          `json:"is_verified" bson:"is_verified"`
}

type UserModel struct{}

func (u *UserModel) Signup(data forms.SignupUserCommand) error {
	collection := dbConnect.Use(databaseName, "user")

	err := collection.Insert(bson.M{
		"name":        data.Name,
		"email":       data.Email,
		"password":    data.Password,
		"is_verified": false,
	})

	return err
}
