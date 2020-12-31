package users

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	passwordUtils "internal/common/password"
	"internal/common/errors"
	jwtUtils "internal/common/token"
	"internal/db"
	"internal/uid"
	"internal/users/types"
)

type userModel struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name"`
	Email string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	IsVerified bool `json:"is_verified" bson:"is_verified"`
	IsDeleted bool `json:"is_deleted" bson:"is_deleted"`
}

type userUpdateModel struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name,omitempty"`
	Email string `json:"email" bson:"email,omitempty"`
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
		IsVerified: true,
		IsDeleted: false,
	}, nil
}

func tryFromUpdateUserCmd(data *userUpdateCommand, id string) (*userUpdateModel, *errors.ResponseError) {
	oid, err := primitive.ObjectIDFromHex(id); if err != nil {
		return nil, errors.InvalidOID()
	}
	return &userUpdateModel {
		ID: oid,
		Name: data.Name,
		Email: data.Email,
	}, nil
}

func (self *userModel) signup(ctx context.Context) (UID.ID, *errors.ResponseError) {
	coll := db.Connection().Use(db.DefaultDatabase, "users")

	if EmailExists(ctx, self.Email) {
		//user already in use
		return nil, errors.EmailAlreadyExistsError()
	}

	res, err := coll.InsertOne(ctx, self); if err != nil {
		return nil, errors.DatabaseError(err)
	}
	return UID.TryFromInterface(res.InsertedID)
}

//TODO: In case we want to add anything else to the jwt
func (self *userModel) jwt() (string, error) {
	return jwtUtils.CreateToken(self.ID.Hex())
}

func authenticate(ctx context.Context, credentials loginUserCommand) (string, *errors.ResponseError){
	user, err := UserByEmail(ctx, credentials.Email); if err != nil {
		//email doesnt exist
		return string(""), errors.EmailDoesNotExistError()
	}

	if user.IsDeleted {
		// User was deleted
		return string(""), errors.EmailDoesNotExistError()
	}

	ok, err := passwordUtils.VerifyPassword(credentials.Password, user.Password)
	if !ok {
		//passwords dont match
		return string(""), errors.InvalidCredentials()
	}
	//get jwt
	jwt, err := user.jwt(); if err != nil {
		//failed to create jwt
		return string(""), errors.JwtError(err)
	}
	return jwt, nil
}

func UserByEmail(ctx context.Context, email string) (userModel, error) {
	coll := db.Connection().Use(db.DefaultDatabase, "users")

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

func fetchUsers(ctx context.Context)([]types.DeliverableUser, *errors.ResponseError) {
	coll := db.Connection().Use(db.DefaultDatabase, "users")

	filter := bson.D{{"is_deleted", false}}
	cursor, err := coll.Find(ctx, filter)
	defer cursor.Close(ctx)

	var users []types.DeliverableUser
	if err = cursor.All(ctx, &users); err != nil {
		return nil, errors.DatabaseError(err)
	}
	return users, nil
}

func (self *userUpdateModel) patch(ctx context.Context, upsert bool) (*types.DeliverableUser, *errors.ResponseError) {
	coll := db.Connection().Use(db.DefaultDatabase, "users")
	opts := options.FindOneAndUpdate().SetUpsert(upsert)

	filter := bson.D{{"_id", self.ID}, {"is_deleted", false}}
	update := bson.D{{"$set", self}}
	var updatedDocument types.DeliverableUser
	err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedDocument)

	if err != nil {
		return nil, errors.PutFailed(err)
	}

	err = coll.FindOne(ctx, filter).Decode(&updatedDocument)
	if err != nil {
		return nil, errors.DatabaseError(err)
	}
	return &updatedDocument, nil
}

func deleteUser(ctx context.Context, id string) (*errors.ResponseError) {
	coll := db.Connection().Use(db.DefaultDatabase, "users")
	oid, err := primitive.ObjectIDFromHex(id); if err != nil {
		return errors.InvalidOID()
	}

	filter := bson.D{{"_id", oid}}
	update := bson.D{{"$set", bson.D{{"is_deleted", true}}}}
	var updatedDocument types.DeliverableUser
	dbErr := coll.FindOneAndUpdate(ctx, filter, update).Decode(&updatedDocument)

	if dbErr != nil {
		return errors.DeleteFailed()
	}
	return nil
}
