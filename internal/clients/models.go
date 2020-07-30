package clients

import (
	"context"
	"bytes"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	errors "internal/common"
	"internal/db"
	"internal/uid"
	"internal/clients/types"
	jobTypes "internal/jobs/types"
)

type clientModel struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name string `json:"name" bson:"name"`
	Phone string `json:"phone" bson:"phone"`
	Jobs []primitive.ObjectID `json:"jobs" bson:"jobs"`
}

type populatedClientModel struct {
	ID primitive.ObjectID `json:"_id"`
	Name string `json:"name"`
	Phone string `json:"phone"`
	Jobs []jobTypes.Job `json:"jobs"`
}

func NewClient(name, phone string) types.Client {
	return &clientModel{
		ID: primitive.NewObjectID(),
		Name: name,
		Phone: phone,
		Jobs: []primitive.ObjectID{},
	}
}

func ClientByPhone(ctx context.Context, phone string) types.Client {
	coll := db.Connection().Use(db.DefaultDatabase, "clients")

	var foundClient clientModel
	err := coll.FindOne(ctx, bson.D{{"phone", phone}}).Decode(&foundClient)
	if err != nil {
		return nil
	}
	return &foundClient
}

func (self *clientModel) AttatchJobID(oid primitive.ObjectID) {
	//search for id, insert if not already in the array
	// linear search for now
	const matched int = 0

	for _, id := range self.Jobs {
		result := bytes.Compare([]byte(oid.String()), []byte(id.String()))
		if result == matched {
			return
		}
	}
	self.Jobs = append(self.Jobs, oid)
}

func (self *clientModel) create(ctx context.Context) (UID.ID, *errors.ResponseError) {
	coll := db.Connection().Use(db.DefaultDatabase, "clients")
	res, err := coll.InsertOne(ctx, self); if err != nil {
		return nil, errors.DatabaseError(err)
	}
	return UID.TryFromInterface(res.InsertedID)
}

func (self *clientModel) Put(ctx context.Context) *errors.ResponseError {
	coll := db.Connection().Use(db.DefaultDatabase, "clients")
	opts := options.FindOneAndReplace()
	opts = opts.SetUpsert(true)

	err := coll.FindOneAndReplace(ctx, bson.D{{"_id", self.ID}}, self, opts).Err()
	if err == mongo.ErrNoDocuments {
		return nil
	}

	if err != nil {
		return errors.PutFailed(err)
	}
	return nil
}

func clientByPhone(phone string, ctx context.Context) (*clientModel, *errors.ResponseError) {
	coll := db.Connection().Use(db.DefaultDatabase, "clients")

	var client clientModel
	err := coll.FindOne(ctx, bson.D{{"phone", phone}}).Decode(&client)

	if err != nil {
		return nil, errors.DatabaseError(err)
	}
	return &client, nil
}

func (self *clientModel) populate(ctx context.Context) (*populatedClientModel, *errors.ResponseError) {
	coll := db.Connection().Use(db.DefaultDatabase, "jobs")
	cursor, err := db.Populate(ctx, coll, self.Jobs); if err != nil {
		return nil, errors.DatabaseError(err)
	}
	defer cursor.Close(ctx)

	var jobs []jobTypes.Job
	if err = cursor.All(ctx, &jobs); err != nil {
		return nil, errors.DatabaseError(err)
	}

	return &populatedClientModel{
		ID: self.ID,
		Name: self.Name,
		Phone: self.Phone,
		Jobs: jobs,
	}, nil
}
