package mongo

import (
	"PPA"
	"context"
	"fmt"
	"net/http"
	"pkg/common/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodb "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	AlreadyExists = http.StatusConflict
	NotFound = http.StatusNotFound
	Collection = "jobs"
	DeletedJobsCollection = "deleted_jobs"
)
type Job struct{}

func (j Job) Create(db *mongo.DBConnection, ctx context.Context, job *PPA.Job) (*PPA.Job, error) {
	coll := db.Use(Collection)

	if _, err := coll.InsertOne(ctx, job); err != nil {
		fmt.Println(err)
		return nil, PPA.InternalError
	}
	return job, nil
}

func(j Job) ViewById(db *mongo.DBConnection, ctx context.Context, oid primitive.ObjectID) (*PPA.Job, error) {
	coll := db.Use(Collection)

	var job PPA.Job
	if err := coll.FindOne(ctx, bson.D{{"_id", oid}}).Decode(&job); err != nil {
		if err == mongodb.ErrNoDocuments {
			return nil, PPA.NewAppError(NotFound, "Job Not Found")
		}
		if err != nil {
			return nil, PPA.InternalError
		}
	}

	return &job, nil
}

func(j Job) List(db *mongo.DBConnection, ctx context.Context, fetchOpts PPA.BulkFetchOptions) (*[]PPA.Job, error) {
	if fetchOpts.All {
		return fetchAll(db, ctx, fetchOpts.Sort)
	}

	coll := db.Use(Collection)
	findOpts := options.
	Find().
	SetSkip(int64(fetchOpts.Source)).
	SetLimit(int64(fetchOpts.Next))

	if fetchOpts.Sort {
		findOpts.SetSort(bson.D{{"_id", -1}})
	}

	cursor, err := coll.Find(ctx, bson.D{{}}, findOpts); if err != nil {
		//perhaps a better error like no docs match
		return nil, PPA.InternalError
	}
	defer cursor.Close(ctx)

	var jobs []PPA.Job
	if err = cursor.All(ctx, &jobs); err != nil {
		return nil, PPA.InternalError
	}
	return &jobs, nil
}

func(j Job) Delete(db *mongo.DBConnection, ctx context.Context, oid primitive.ObjectID) error {
	coll := db.Use(DeletedJobsCollection)

	fetched, err := j.ViewById(db, ctx, oid); if err != nil {
		return PPA.NewAppError(NotFound, "Job Not Found")
	}

	if _, insertErr := coll.InsertOne(ctx, fetched); insertErr != nil {
		return PPA.InternalError //insert err, db err
	}

	coll = db.Use(Collection)
	if _, delErr := coll.DeleteOne(ctx, bson.D{{"_id", oid}}); delErr != nil {
		return PPA.InternalError //db error
	}

	return nil
}

func (j Job) Update(db *mongo.DBConnection, ctx context.Context, oid primitive.ObjectID, update *PPA.Job) error {
	coll := db.Use(Collection)

	filter := bson.D{{"_id", oid}}
	updateDoc := bson.D{{"$set", update}}

	var oldDoc PPA.Job
	err := coll.FindOneAndUpdate(ctx, filter, updateDoc).Decode(&oldDoc)

	if err == mongodb.ErrNoDocuments {
		return PPA.NewAppError(NotFound, "Job Not Found")
	}

	if err != nil {
		return PPA.InternalError
	}
	return nil
}

func (j Job) UnassignTinter(db *mongo.DBConnection, ctx context.Context, update *PPA.Job) error {
	coll := db.Use(Collection)

	filter := bson.D{{"_id", update.ID}}
	updateDoc := bson.D{{"$set", bson.D{{"assigned_worker", primitive.NilObjectID}} }}

	var oldDoc PPA.Job
	err := coll.FindOneAndUpdate(ctx, filter, updateDoc).Decode(&oldDoc)
	fmt.Println(oldDoc)

	if err == mongodb.ErrNoDocuments {
		return PPA.NewAppError(NotFound, "Job Not Found")
	}

	if err != nil {
		return PPA.InternalError
	}
	return nil
}

func (j Job) LogEvent(db *mongo.DBConnection, ctx context.Context, update *PPA.Job) {
	if err := j.Update(db, ctx, update.ID, update); err != nil {
		fmt.Println(err)
	}
}

func (j Job) Stream(db *mongo.DBConnection, ctx context.Context, tx chan *PPA.StreamResult) {
	coll := db.Use(Collection)

	changeStream, streamErr := coll.Watch(ctx, mongodb.Pipeline{}); if streamErr != nil {
		fmt.Println(streamErr)
		return
	}

	defer changeStream.Close(ctx)

	for changeStream.Next(ctx) {
		var data bson.M
		if err := changeStream.Decode(&data); err != nil {
			fmt.Println(err)
			return
		}
		oid := data["documentKey"].(primitive.M)["_id"].(primitive.ObjectID)

		var job PPA.Job
		_ = coll.FindOne(ctx, bson.D{{"_id", oid}}).Decode(&job)

		tx <-&PPA.StreamResult {
			EventType:data["operationType"].(string),
			Data: job,
		}
	}
}

func fetchAll(db *mongo.DBConnection, ctx context.Context, sort bool) (*[]PPA.Job, error) {
	coll := db.Use(Collection)
	opts := options.Find()

	if sort {
		opts.SetSort(bson.D{{"_id", -1}})
	}

	cursor, err := coll.Find(ctx, bson.D{{}}, opts); if err != nil {
		return nil, PPA.InternalError
	}
	defer cursor.Close(ctx)

	var jobs []PPA.Job
	if err = cursor.All(ctx, &jobs); err != nil {
		// TODO: check for err no docs?
		return nil, PPA.InternalError
	}
	return &jobs, nil
}
