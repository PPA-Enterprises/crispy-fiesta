package users

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
	passwordUtils "internal/common"
	"internal/db"
	"internal/uid"
)

type userModel struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name"`
	Email string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	IsVerified bool `json:"is_verified" bson:"is_verified"`
}

func tryFromSignupUserCmd(data *signupUserCommand) (*userModel, error) {
	encrypted, err := passwordUtils.HashFromPlaintext(data.Password)
	if err != nil {
		return nil, err
	}

	return &userModel{
		Name: data.Name,
		Email: data.Email,
		Password: encrypted,
	}, nil
}

func (u *userModel) signup(ctx context.Context) (UID.ID, error) {
	coll := db.Connection().Use(db.DefaultDatabase, "user")

	if EmailExists(ctx, u.Email) {
		//user already in use
		return nil, errors.New("Email already exists")
	}

	res, err := coll.InsertOne(ctx, u); if err != nil {
		return nil, errors.New("Server Error")
	}
	return UID.IdFromInterface(res.InsertedID)
}

func UserByEmail(ctx context.Context, email string) (userModel, error) {
	coll := db.Connection().Use(db.DefaultDatabase, "user")

	var foundUser userModel
	err := coll.FindOne(ctx, bson.D{{"email", email}}).Decode(&foundUser)
	return foundUser, err
}

func EmailExists(ctx context.Context, email string) bool {
	_, err := UserByEmail(ctx, email); if err != nil {
		return false
	}
	return true
}
