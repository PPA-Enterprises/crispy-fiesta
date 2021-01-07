package users

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"internal/common/errors"
	"internal/db"
	"internal/users/types"
)

type updateableUser struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name,omitempty"`
	Email string `json:"email" bson:"email,omitempty"`
}

func tryFromUpdateUserCmd(data *userUpdateCommand, id string) (*updateableUser, *errors.ResponseError) {
	oid, err := primitive.ObjectIDFromHex(id); if err != nil {
		return nil, errors.InvalidOID()
	}
	return &updateableUser {
		ID: oid,
		Name: data.Name,
		Email: data.Email,
	}, nil
}

func (self *updateableUser) patch(ctx context.Context, upsert bool) (*types.DeliverableUser, *errors.ResponseError) {
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
