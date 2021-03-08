package mongo
import (
	"PPA"
	"net/http"
	"context"
	"pkg/common/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	NotFound = http.StatusNotFound
	UserCollection = "users"
	Coll = "eventlogs"
)

type Eventlog struct{}
func(ev Eventlog) List(db *mongo.DBConnection, ctx context.Context, fetchOpts PPA.BulkFetchOptions) (*[]PPA.LogEvent, error) {

	if fetchOpts.All {
		return fetchAll(db, ctx, fetchOpts.Sort, Coll)
	}

	coll := db.Use(Coll)
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

	var evlogs []PPA.LogEvent
	if err = cursor.All(ctx, &evlogs); err != nil {
		return nil, PPA.InternalError
	}
	return &evlogs, nil
}

func fetchAll(db *mongo.DBConnection, ctx context.Context, sort bool, collection string) (*[]PPA.LogEvent, error) {
	coll := db.Use(collection)
	opts := options.Find()

	if sort {
		opts.SetSort(bson.D{{"_id", -1}})
	}

	cursor, err := coll.Find(ctx, bson.D{{}}, opts); if err != nil {
		return nil, PPA.InternalError
	}
	defer cursor.Close(ctx)

	var evlogs []PPA.LogEvent
	if err = cursor.All(ctx, &evlogs); err != nil {
		// TODO: check for err no docs?
		return nil, PPA.InternalError
	}
	return &evlogs, nil
}
