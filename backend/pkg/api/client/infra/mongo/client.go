package mongo
import (
	"PPA"
	"net/http"
	"context"
	"pkg/common/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodb "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	AlreadyExists = http.StatusConflict
	NotFound = http.StatusNotFound
)

type Client struct{}
func (c Client) Create(db *mongo.DBConnection, ctx context.Context, client *PPA.Client) (*PPA.Client, error) {
	coll := db.Use("clients")

	if(c.phoneExists(db, ctx, client.Phone)) {
		return nil, PPA.NewAppError(AlreadyExists, "Phone Number Already In Use")
	}

	if _, err := coll.InsertOne(ctx, client); err != nil {
		return nil, PPA.InternalError
	}
	return client, nil
}

func(c Client) ViewById(db *mongo.DBConnection, ctx context.Context, oid primitive.ObjectID) (*PPA.Client, error) {
	coll := db.Use("clients")

	var client PPA.Client
	if err := coll.FindOne(ctx, bson.D{{"_id", oid}}).Decode(&client); err != nil {
		if err == mongodb.ErrNoDocuments {
			return nil, PPA.NewAppError(NotFound, "Client Not Found")
		}
		if err != nil {
			return nil, PPA.InternalError
		}
	}
	return &client, nil
}

func(c Client) ViewByPhone(db *mongo.DBConnection, ctx context.Context, phone string) (*PPA.Client, error) {
	coll := db.Use("clients")

	var client PPA.Client
	err := coll.FindOne(ctx, bson.D{{"phone", phone}}).Decode(&client)
		if err == mongodb.ErrNoDocuments {
			return nil, PPA.NewAppError(NotFound, "Client Not Found")
		}
		if err != nil {
			return nil, PPA.InternalError
		}
	return &client, nil
}

func(c Client) List(db *mongo.DBConnection, ctx context.Context, fetchOpts PPA.BulkFetchOptions) (*[]PPA.Client, error) {
	if fetchOpts.All {
		return fetchAll(db, ctx, fetchOpts.Sort)
	}

	coll := db.Use("clients")
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

	var clients []PPA.Client
	if err = cursor.All(ctx, &clients); err != nil {
		return nil, PPA.InternalError
	}
	return &clients, nil
}

func(c Client) Delete(db *mongo.DBConnection, ctx context.Context, oid primitive.ObjectID) error {
	coll := db.Use("deleted_clients")

	fetched, err := c.ViewById(db, ctx, oid); if err != nil {
		return PPA.NewAppError(NotFound, "Client Not Found")
	}

	if _, insertErr := coll.InsertOne(ctx, fetched); insertErr != nil {
		return PPA.InternalError //insert err, db err
	}

	coll = db.Use("clients")
	if _, delErr := coll.DeleteOne(ctx, bson.D{{"_id", oid}}); delErr != nil {
		return PPA.InternalError //db error
	}
	return nil
}

func (c Client) Update(db *mongo.DBConnection, ctx context.Context, oid primitive.ObjectID, update *PPA.Client) error {
	coll := db.Use("clients")

	filter := bson.D{{"_id", oid}}
	updateDoc := bson.D{{"$set", update}}

	var oldDoc PPA.Client
	err := coll.FindOneAndUpdate(ctx, filter, updateDoc).Decode(&oldDoc)

	if err == mongodb.ErrNoDocuments {
		return PPA.NewAppError(NotFound, "Client Not Found")
	}

	if err != nil {
		return PPA.InternalError
	}
	return nil
}

func (c Client) Populate(db *mongo.DBConnection, ctx context.Context, oids []primitive.ObjectID) ([]PPA.Job, error) {
	coll := db.Use("jobs")
	cursor, err := db.Populate(ctx, coll, oids); if err != nil {
		return []PPA.Job{}, PPA.InternalError
	}
	defer cursor.Close(ctx)

	var jobs []PPA.Job
	if err = cursor.All(ctx, &jobs); err != nil {
		return []PPA.Job{}, PPA.InternalError
	}
	return jobs, nil
}

func (c Client) RemoveJob(db *mongo.DBConnection, ctx context.Context, clientPhone string, jobOid primitive.ObjectID) error {
	fetched, err := c.ViewByPhone(db, ctx, clientPhone); if err != nil {
		return err
	}
	fetched.FindAndRemoveJob(jobOid)

	coll := db.Use("clients")

	filter := bson.D{{"_id", fetched.ID}}
	updateDoc := bson.D{{"$set", fetched}}

	var oldDoc PPA.Client
	err = coll.FindOneAndUpdate(ctx, filter, updateDoc).Decode(&oldDoc)

	if err == mongodb.ErrNoDocuments {
		return PPA.NewAppError(NotFound, "Client Not Found")
	}

	if err != nil {
		return PPA.InternalError
	}
	return nil
}

func (c Client) phoneExists(db *mongo.DBConnection, ctx context.Context, phone string) bool {
	if _, err := c.ViewByPhone(db, ctx, phone); err != nil {
		return false
	}
	return true
}

func fetchAll(db *mongo.DBConnection, ctx context.Context, sort bool) (*[]PPA.Client, error) {
	coll := db.Use("clients")
	opts := options.Find()

	if sort {
		opts.SetSort(bson.D{{"_id", -1}})
	}

	cursor, err := coll.Find(ctx, bson.D{{}}, opts); if err != nil {
		return nil, PPA.InternalError
	}
	defer cursor.Close(ctx)

	var clients []PPA.Client
	if err = cursor.All(ctx, &clients); err != nil {
		// TODO: check for err no docs?
		return nil, PPA.InternalError
	}
	return &clients, nil
}
