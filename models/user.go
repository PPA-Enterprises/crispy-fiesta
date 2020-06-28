package models

import (
	"github.com/PPA-Enterprises/crispy-fiesta/forms"
	"github.com/PPA-Enterprises/crispy-fiesta/helpers"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string             `json:"name" bson:"name"`
	Email      string             `json:"email" bson:"email"`
	Password   string             `json:"password" bson:"password"`
	IsVerified bool               `json:"is_verified" bson:"is_verified"`
}

type UserModel struct{}

func (u *UserModel) Signup(data forms.SignupUserCommand) (*mongo.InsertOneResult, error) {
	collection := dbConnect.Use(databaseName, "user")

	hash, err := helpers.GenerateFromPassword(data.Password, helpers.P)

	if err != nil {
		panic(err)
	}

	result, err := collection.Insert(bson.D{
		{"name", data.Name},
		{"email", data.Email},
		{"password", hash},
		{"is_verified", false},
	}, nil)

	return result, err
}
