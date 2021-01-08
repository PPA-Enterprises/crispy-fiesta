package clients

import (
	"context"
	"internal/clients/types"
	"internal/common/errors"
	"internal/db"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type updateableClient struct {
	ID    primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Name  string               `json:"name" bson:"name,omitempty"`
	Phone string               `json:"phone" bson:"phone,omitempty"`
}

func tryFromUpdateClientCmd(data *updateClientCmd, id string) (*updateableClient, *errors.ResponseError) {
	clientOID, err := primitive.ObjectIDFromHex(id); if err != nil {
		return nil, errors.InvalidOID()
	}
	return &updateableClient{
		ID:    clientOID,
		Name:  data.Name,
		Phone: data.Phone,
	}, nil
}

func (self *updateableClient) patch(ctx context.Context, upsert bool) (*types.PopulatedClientModel, *errors.ResponseError) {
	coll := db.Connection().Use(db.DefaultDatabase, "clients")
	opts := options.FindOneAndUpdate().SetUpsert(upsert)

	filter := bson.D{{"_id", self.ID}}
	update := bson.D{{"$set", self}}
	var updatedDocument clientModel
	err := coll.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedDocument)

	if err != nil {
		return nil, errors.PutFailed(err)
	}

	err = coll.FindOne(ctx, filter).Decode(&updatedDocument)
	if err != nil {
		return nil, errors.DatabaseError(err)
	}
	return updatedDocument.Populate(ctx)
}
