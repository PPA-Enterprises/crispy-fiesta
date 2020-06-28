package models

import (
	"context"
	"errors"
	"log"

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

	var foundEmail User

	log.Print(data.Email)

	err := collection.FindOne(context.Background(), bson.D{{"email", data.Email}}).Decode(&foundEmail)
	log.Print(foundEmail)

	if &foundEmail != nil {
		return nil, errors.New("Email already exists")
	}

	hash, err := helpers.GenerateFromPassword(data.Password, helpers.P)

	if err != nil {
		panic(err)
	}

	result, err := collection.InsertOne(context.Background(), bson.D{
		{"name", data.Name},
		{"email", data.Email},
		{"password", hash},
		{"is_verified", false},
	})
	return result, err
}
