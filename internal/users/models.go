package users

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
	passwordUtils "internal/common"
	errors "internal/common"
	//jwtUtils "internal/common"
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

func tryFromSignupUserCmd(data *signupUserCommand) (*userModel, *errors.ResponseError) {
	encrypted, err := passwordUtils.HashFromPlaintext(data.Password)
	if err != nil {
		return nil, errors.ArgonHashError(err)
	}

	return &userModel{
		Name: data.Name,
		Email: data.Email,
		Password: encrypted,
	}, nil
}

func (self *userModel) signup(ctx context.Context) (UID.ID, *errors.ResponseError) {
	coll := db.Connection().Use(db.DefaultDatabase, "user")

	if EmailExists(ctx, self.Email) {
		//user already in use
		return nil, errors.EmailAlreadyExistsError()
	}

	res, err := coll.InsertOne(ctx, self); if err != nil {
		return nil, errors.DatabaseError(err)
	}
	return UID.IdFromInterface(res.InsertedID)
}

/*func (self *userModel) jwt() (string, error) {

}*/

func authenticate(ctx context.Context, credentials loginUserCommand) {
	user, err := UserByEmail(ctx, credentials.Email); if err != nil {
		//email doesnt exist
	}

	ok, err := passwordUtils.VerifyPassword(credentials.Password, user.Password)
	if !ok {
		//passwords dont match
	}
	//get jwt
	//jwt, err := jwtUtils.CreateToken
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

