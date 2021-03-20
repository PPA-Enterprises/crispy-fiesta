package mongo
import (
	"PPA"
	"net/http"
	"fmt"
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
	Collection = "tinters"
	DeletedTintersCollection = "deleted_tinters"
	JobsCollection = "jobs"
)

type Tinter struct{}
func (t Tinter) Create(db *mongo.DBConnection, ctx context.Context, tinter *PPA.Tinter) (*PPA.Tinter, error) {
	fmt.Println("CALLED")
	coll := db.Use(Collection)

	if(t.phoneExists(db, ctx, tinter.Phone)) {
		return nil, PPA.NewAppError(AlreadyExists, "Phone Number Already In Use")
	}

	if _, err := coll.InsertOne(ctx, tinter); err != nil {
		return nil, PPA.InternalError
	}
	return tinter, nil
}

func(t Tinter) ViewById(db *mongo.DBConnection, ctx context.Context, oid primitive.ObjectID) (*PPA.Tinter, error) {
	coll := db.Use(Collection)

	var tinter PPA.Tinter
	if err := coll.FindOne(ctx, bson.D{{"_id", oid}}).Decode(&tinter); err != nil {
		if err == mongodb.ErrNoDocuments {
			return nil, PPA.NewAppError(NotFound, "Tinter Not Found")
		}
		if err != nil {
			return nil, PPA.InternalError
		}
	}
	return &tinter, nil
}

func(t Tinter) ViewByPhone(db *mongo.DBConnection, ctx context.Context, phone string) (*PPA.Tinter, error) {
	coll := db.Use(Collection)

	var tinter PPA.Tinter
	err := coll.FindOne(ctx, bson.D{{"phone", phone}}).Decode(&tinter)
		if err == mongodb.ErrNoDocuments {
			return nil, PPA.NewAppError(NotFound, "Tinter Not Found")
		}
		if err != nil {
			return nil, PPA.InternalError
		}
	return &tinter, nil
}

func(t Tinter) List(db *mongo.DBConnection, ctx context.Context, fetchOpts PPA.BulkFetchOptions) (*[]PPA.Tinter, error) {
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

	var tinters []PPA.Tinter
	if err = cursor.All(ctx, &tinters); err != nil {
		return nil, PPA.InternalError
	}
	return &tinters, nil
}

func(t Tinter) Delete(db *mongo.DBConnection, ctx context.Context, oid primitive.ObjectID) error {
	coll := db.Use(DeletedTintersCollection)

	fetched, err := t.ViewById(db, ctx, oid); if err != nil {
		return PPA.NewAppError(NotFound, "Tinter Not Found")
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

func (t Tinter) Update(db *mongo.DBConnection, ctx context.Context, oid primitive.ObjectID, update *PPA.Tinter) error {
	coll := db.Use(Collection)

	filter := bson.D{{"_id", oid}}
	updateDoc := bson.D{{"$set", update}}

	var oldDoc PPA.Tinter
	err := coll.FindOneAndUpdate(ctx, filter, updateDoc).Decode(&oldDoc)

	if err == mongodb.ErrNoDocuments {
		return PPA.NewAppError(NotFound, "Tinter Not Found")
	}

	if err != nil {
		return PPA.InternalError
	}
	return nil
}

func (t Tinter) LogEvent(db *mongo.DBConnection, ctx context.Context, update *PPA.Tinter) {
	if err := t.Update(db, ctx, update.ID, update); err != nil {
		fmt.Println(err)
	}
}

func (t Tinter) Populate(db *mongo.DBConnection, ctx context.Context, oids []primitive.ObjectID) ([]PPA.Job, error) {
	coll := db.Use(JobsCollection)
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

/*
func (t Tinter) RemoveJob(db *mongo.DBConnection, ctx context.Context, tinterPhone string, jobOid primitive.ObjectID) error {
	fetched, err := t.ViewByPhone(db, ctx, tinterPhone); if err != nil {
		return err
	}
	fetched.FindAndRemoveJob(jobOid)

	coll := db.Use(Collection)

	filter := bson.D{{"_id", fetched.ID}}
	updateDoc := bson.D{{"$set", fetched}}

	var oldDoc PPA.Tinter
	err = coll.FindOneAndUpdate(ctx, filter, updateDoc).Decode(&oldDoc)

	if err == mongodb.ErrNoDocuments {
		return PPA.NewAppError(NotFound, "Tinter Not Found")
	}

	if err != nil {
		return PPA.InternalError
	}
	return nil
}*/

func (t Tinter) phoneExists(db *mongo.DBConnection, ctx context.Context, phone string) bool {
	if _, err := t.ViewByPhone(db, ctx, phone); err != nil {
		return false
	}
	return true
}

func fetchAll(db *mongo.DBConnection, ctx context.Context, sort bool) (*[]PPA.Tinter, error) {
	coll := db.Use(Collection)
	opts := options.Find()

	if sort {
		opts.SetSort(bson.D{{"_id", -1}})
	}

	cursor, err := coll.Find(ctx, bson.D{{}}, opts); if err != nil {
		return nil, PPA.InternalError
	}
	defer cursor.Close(ctx)

	var tinters []PPA.Tinter
	if err = cursor.All(ctx, &tinters); err != nil {
		// TODO: check for err no docs?
		return nil, PPA.InternalError
	}
	return &tinters, nil
}
