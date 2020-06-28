package models

import (
	"context"
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

	var foundEmail bson.M
	err := collection.FindOne(context.Background(), bson.D{{"email", data.Email}}).Decode(&foundEmail)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return
		}
		log.Fatal(err)
	}

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
