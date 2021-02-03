package clients

import (
	"bytes"
	"context"
	"internal/clients/types"
	"internal/common/errors"
	"internal/db"
	"internal/event_log"
	eventLogTypes "internal/event_log/types"
	jobTypes "internal/jobs/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type clientModel struct {
	ID		primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Name	string               `json:"name" bson:"name"`
	Phone	string               `json:"phone" bson:"phone"`
	Jobs	[]primitive.ObjectID `json:"jobs" bson:"jobs"`
	Log		[]eventLogTypes.NormalizedLoggedEvent `json:"log" bson:"log"`
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

func (self *clientModel) AttatchJobID(ctx context.Context, oid primitive.ObjectID, editor *eventLogTypes.Editor) *errors.ResponseError {
	//search for id, insert if not already in the array
	// linear search for now
	const matched int = 0

	for _, id := range self.Jobs {
		result := bytes.Compare([]byte(oid.String()), []byte(id.String()))
		if result == matched {
			return errors.JobAlreadyExistsError()
		}
	}
	self.Jobs = append(self.Jobs, oid)
	return self.put(ctx, true, editor)
}

func (self *clientModel) put(ctx context.Context, upsert bool, editor *eventLogTypes.Editor) *errors.ResponseError {
	coll := db.Connection().Use(db.DefaultDatabase, "clients")
	opts := options.FindOneAndReplace()
	opts = opts.SetUpsert(true)

	err := coll.FindOneAndReplace(ctx, bson.D{{"_id", self.ID}}, self, opts).Err()
	if err == mongo.ErrNoDocuments {
		if upsert {
			//client was created
			loggedClient := event_log.LogCreated(ctx, self.logable(), editor)
			_ = appendLog(ctx, self, loggedClient)
			return nil
		} else {
			return errors.PutFailed(err)
		}
	}

	if err != nil {
		return errors.PutFailed(err)
	}
	//client was updated with an appended job
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

func (self *clientModel) logable() *types.LogableClient {
	return &types.LogableClient {
		ID: self.ID.Hex(),
		Name: self.Name,
		Phone: self.Phone,
	}
}

/*
//https://github.com/mongodb/mongo-go-driver/blob/51421e413403fe3c9b0097147841f752421133e4/examples/documentation_examples/examples.go#L293
func fuzzySearch(ctx context.Context, opts *FuzzySearch) ([]types.UnpopulatedClientModel, *errors.ResponseError) {
	coll := db.Connection().Use(db.DefaultDatabase, "clients")

	findOptions := options.
	Find().
	SetSkip(int64(opts.Source)).
	SetLimit(int64(opts.Next))

	regexQuery := powersetRegex(opts.Term)
	fmt.Println(regexQuery)
	filter := bson.D{{"name", primitive.Regex{Pattern: regexQuery, Options:"i"}}}
	cursor, err := coll.Find(ctx, filter, findOptions)
	defer cursor.Close(ctx)

	var clients []types.UnpopulatedClientModel
	if err = cursor.All(ctx, &clients); err != nil {
		return nil, errors.DatabaseError(err)
	}
	return clients, nil
}

func powersetRegex(term string) string {
	var termArr = make([]string, 0, len(term))
	termArr = append(termArr, "..")
	for i:=0; i<len(term); i++ {
		termArr = append(termArr, string(term[i]))
	}
	powerset := combinations.All(termArr)
	powerset = powerset[1:]

	regex := "^" + term + "$" + "|"
	regex = regex + "("
	for i:=len(powerset)-1; i>=0; i-- {
		regexTerm := strings.Join(powerset[i], "")
		regexTerm = "^" + regexTerm + "$"
		regex = regex + "(" + regexTerm + ")" //+ "|"
	}
	regex = regex + "){1}" //+ "|"

	return regex
}
*/

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

	cursor, err := coll.Find(ctx, bson.D{{}}, opts)
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

	cursor, err := coll.Find(ctx, bson.D{{}}, findOptions)
	defer cursor.Close(ctx)

	var clients []types.UnpopulatedClientModel
	if err = cursor.All(ctx, &clients); err != nil {
		return nil, errors.DatabaseError(err)
	}
	return clients, nil
}

func RemoveJob(ctx context.Context, clientID, jobID string) *errors.ResponseError {
	const matched int = 0

	client, err := clientByID(ctx, clientID); if err != nil {
		return err
	}

	for i, oid := range client.Jobs {
		result := bytes.Compare([]byte(jobID), []byte(oid.Hex()))
		if result == matched {
			// preserve the order. Idiomatic way
			client.Jobs = append(client.Jobs[:i], client.Jobs[i+1:]...)
		}
	}

	coll := db.Connection().Use(db.DefaultDatabase, "clients")
	filter := bson.D{{"_id", client.ID}}
	update := bson.D{{"$set", bson.D{{"jobs", client.Jobs}}}}

	var updatedDoc clientModel
	updateErr := coll.FindOneAndUpdate(ctx, filter, update).Decode(&updatedDoc)
	if updateErr != nil {
		return errors.DatabaseError(updateErr)
	}
	return nil
}

func deleteByID(ctx context.Context, clientID string) *errors.ResponseError {
	coll := db.Connection().Use(db.DefaultDatabase, "deleted_clients")
	client, err := clientByID(ctx, clientID); if err != nil {
		return err
	}

	_, insertErr := coll.InsertOne(ctx, client); if insertErr != nil {
		return errors.DatabaseError(insertErr)
	}

	jobDestroyer := jobTypes.DeletorFactory()

	for _, oid := range client.Jobs {
		_ := jobDestroyer.DeleteByID(ctx, oid)
	}

	coll = db.Connection().Use(db.DefaultDatabase, "clients")
	_, delErr := coll.DeleteOne(ctx, bson.D{{"_id", client.ID}})
	if delErr != nil {
		return errors.DatabaseError(delErr)
	}
	return nil
}
