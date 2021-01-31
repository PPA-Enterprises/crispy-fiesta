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
)
type Job struct{}

func (j Job) Create(db *mongo.DBConnection, ctx context.Context, job *PPA.Job) (*PPA.Job, error) {
	coll := db.Use("jobs")

	if _, err := coll.InsertOne(ctx, job); err != nil {
		fmt.Println(err)
		return nil, PPA.InternalError
	}
	return job, nil
}

func(j Job) ViewById(db *mongo.DBConnection, ctx context.Context, oid primitive.ObjectID) (*PPA.Job, error) {
	coll := db.Use("jobs")

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

	coll := db.Use("jobs")
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
	coll := db.Use("deleted_jobs")

	fetched, err := j.ViewById(db, ctx, oid); if err != nil {
		return PPA.NewAppError(NotFound, "Job Not Found")
	}

	if _, insertErr := coll.InsertOne(ctx, fetched); insertErr != nil {
		return PPA.InternalError //insert err, db err
	}

	coll = db.Use("jobs")
	if _, delErr := coll.DeleteOne(ctx, bson.D{{"_id", oid}}); delErr != nil {
		return PPA.InternalError //db error
	}

	return nil
}

func (j Job) Update(db *mongo.DBConnection, ctx context.Context, oid primitive.ObjectID, update *PPA.Job) error {
	coll := db.Use("jobs")

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

func fetchAll(db *mongo.DBConnection, ctx context.Context, sort bool) (*[]PPA.Job, error) {
	coll := db.Use("jobs")
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
