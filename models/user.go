package models

import (
	"github.com/PPA-Enterprises/crispy-fiesta/forms"
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

	result, err := collection.Insert(bson.D{
		{"name", data.Name},
		{"email", data.Email},
		{"password", data.Password},
		{"is_verified", false},
	}, nil)
	return result, err
}
