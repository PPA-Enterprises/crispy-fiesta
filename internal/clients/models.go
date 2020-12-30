package clients

import (
	"bytes"
	"context"
	"internal/clients/types"
	"internal/common/errors"
	"internal/db"
	jobTypes "internal/jobs/types"
	"internal/uid"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type clientModel struct {
	ID    primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Name  string               `json:"name" bson:"name"`
	Phone string               `json:"phone" bson:"phone"`
	Jobs  []primitive.ObjectID `json:"jobs" bson:"jobs"`
}

func tryFromUpdateClientCmd(data *updateClientCmd) (*clientModel, *errors.ResponseError) {
	clientOID, err := primitive.ObjectIDFromHex(data.ID)
	if err != nil {
		return nil, errors.InvalidOID()
	}
	return &clientModel{
		ID:    clientOID,
		Name:  data.Name,
		Phone: data.Phone,
		Jobs:  normalize(data.Jobs),
	}, nil
}

func normalize(j []jobTypes.Job) []primitive.ObjectID {
	oids := make([]primitive.ObjectID, 0)
	for _, job := range j {
		oids = append(oids, job.ID)
	}
	return oids
}

func NewClient(name, phone string) types.Client {
	return &clientModel{
		ID:    primitive.NewObjectID(),
		Name:  name,
		Phone: phone,
		Jobs:  []primitive.ObjectID{},
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
	res, err := coll.InsertOne(ctx, self)
	if err != nil {
		return nil, errors.DatabaseError(err)
	}
	return UID.TryFromInterface(res.InsertedID)
}

func (self *clientModel) Put(ctx context.Context, upsert bool) *errors.ResponseError {
	coll := db.Connection().Use(db.DefaultDatabase, "clients")
	opts := options.FindOneAndReplace()
	opts = opts.SetUpsert(true)

	err := coll.FindOneAndReplace(ctx, bson.D{{"_id", self.ID}}, self, opts).Err()
	if err == mongo.ErrNoDocuments {
		if upsert {
			return nil
		} else {
			return errors.PutFailed(err)
		}
	}

	if err != nil {
		return errors.PutFailed(err)
	}
	return nil
}

func (self *clientModel) Populate(ctx context.Context) (*types.PopulatedClientModel, *errors.ResponseError) {
	coll := db.Connection().Use(db.DefaultDatabase, "jobs")
	cursor, err := db.Populate(ctx, coll, self.Jobs)
	if err != nil {
		return nil, errors.DatabaseError(err)
	}
	defer cursor.Close(ctx)

	var jobs []jobTypes.Job
	if err = cursor.All(ctx, &jobs); err != nil {
		return nil, errors.DatabaseError(err)
	}

	return &types.PopulatedClientModel{
		ID:    self.ID,
		Name:  self.Name,
		Phone: self.Phone,
		Jobs:  jobs,
	}, nil
}

//TODO: fuzzy search for client
func fuzzySearch(ctx context.Context, query string, quantity int) ([]clientModel, *errors.ResponseError) {
	coll := db.Connection().Use(db.DefaultDatabase, "clients")

	filter := bson.D{{"name", primitive.Regex{Pattern: query, Options: "i"}}}

	cursor, err := coll.Find(ctx, filter)
	defer cursor.Close(ctx)

	var clients []clientModel
	if err = cursor.All(ctx, &clients); err != nil {
		return nil, errors.DatabaseError(err)
	}
	return clients, nil
}

func populateClients(ctx context.Context, clients []clientModel) []types.PopulatedClientModel {

	populatedClients := make([]types.PopulatedClientModel, 0, len(clients))
	for _, c := range clients {
		client, err := c.Populate(ctx)
		//just skip bad records
		if err == nil {
			populatedClients = append(populatedClients, *client)
		}
	}
	return populatedClients
}

func clientByID(ctx context.Context, id string) (*clientModel, *errors.ResponseError) {
	coll := db.Connection().Use(db.DefaultDatabase, "clients")

	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.InvalidOID()
	}

	var foundClient clientModel
	err = coll.FindOne(ctx, bson.D{{"_id", oid}}).Decode(&foundClient)
	if err != nil {
		return nil, errors.DoesNotExist()
	}
	return &foundClient, nil
}

func fetchAll(ctx context.Context, sort bool) ([]types.UnpopulatedClientModel, *errors.ResponseError) {
	coll := db.Connection().Use(db.DefaultDatabase, "clients")
	opts := options.Find()

	if sort {
		opts.SetSort(bson.D{{"_id", -1}})
	}

	cursor, err := coll.Find(ctx, opts)
	defer cursor.Close(ctx)
	var clients []types.UnpopulatedClientModel

	if err = cursor.All(ctx, &clients); err != nil {
		return nil, errors.DatabaseError(err)
	}

	return clients, nil
}

func fetch(ctx context.Context, fetchOpts *BulkFetch) ([]types.UnpopulatedClientModel, *errors.ResponseError) {
	if fetchOpts.All {
		return fetchAll(ctx, fetchOpts.Sort)
	}

	coll := db.Connection().Use(db.DefaultDatabase, "clients")

	findOptions := options.
	Find().
	SetSkip(int64(fetchOpts.Source)).
	SetLimit(int64(fetchOpts.Next))

	if fetchOpts.Sort {
		findOptions.SetSort(bson.D{{"_id", -1}})
	}

	cursor, err := coll.Find(ctx, findOptions)
	defer cursor.Close(ctx)

	var clients []types.UnpopulatedClientModel
	if err = cursor.All(ctx, &clients); err != nil {
		return nil, errors.DatabaseError(err)
	}
	return clients, nil
}
