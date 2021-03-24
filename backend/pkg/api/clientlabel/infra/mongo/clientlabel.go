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
	Collection = "clientlabels"
)

type ClientLabel struct{}

func (cl ClientLabel) Create(db *mongo.DBConnection, ctx context.Context, label *PPA.ClientLabel) (*PPA.ClientLabel, error) {
	coll := db.Use(Collection)
	label.IsDeleted = false

	if _, err := coll.InsertOne(ctx, label); err != nil {
		return nil, PPA.InternalError
	}

	return label, nil
}

func(cl ClientLabel) ViewById(db *mongo.DBConnection, ctx context.Context, oid primitive.ObjectID) (*PPA.ClientLabel, error) {
	coll := db.Use(Collection)

	var label PPA.ClientLabel
	if err := coll.FindOne(ctx, bson.D{{"_id", oid}, {"is_deleted", false}}).Decode(&label); err != nil {
		if err == mongodb.ErrNoDocuments {
			return nil, PPA.NewAppError(NotFound, "Label Not Found")
		}
		if err != nil {
			return nil, PPA.InternalError
		}
	}

	return &label, nil
}

func(cl ClientLabel) ViewByLabelName(db *mongo.DBConnection, ctx context.Context, labelName string) (*PPA.ClientLabel, error) {
	coll := db.Use(Collection)

	var label PPA.ClientLabel
	if err := coll.FindOne(ctx, bson.D{{"label_name", labelName}, {"is_deleted", false}}).Decode(&label); err != nil {
		if err == mongodb.ErrNoDocuments {
			return nil, PPA.NewAppError(NotFound, "Label Not Found")
		}
		if err != nil {
			return nil, PPA.InternalError
		}
	}

	return &label, nil
}

func(cl ClientLabel) List(db *mongo.DBConnection, ctx context.Context) (*[]PPA.ClientLabel, error) {
	coll := db.Use(Collection)

	//check error?
	cursor, err := coll.Find(ctx, bson.D{{"is_deleted", false}})
	defer cursor.Close(ctx)

	var labels []PPA.ClientLabel
	if err = cursor.All(ctx, &labels); err != nil {
		return nil, PPA.InternalError
	}
	return &labels, nil
}


func (cl ClientLabel) Update(db *mongo.DBConnection, ctx context.Context, oid primitive.ObjectID, update *PPA.ClientLabel) error {
	coll := db.Use(Collection)

	filter := bson.D{{"_id", oid}}
	updateDoc := bson.D{{"$set", update}}

	var oldDoc PPA.ClientLabel
	err := coll.FindOneAndUpdate(ctx, filter, updateDoc).Decode(&oldDoc)

	if err == mongodb.ErrNoDocuments {
		return PPA.NewAppError(NotFound, "Label Not Found")
	}

	if err != nil {
		return PPA.InternalError
	}
	return nil
}

func (cl ClientLabel) PutLabels(db *mongo.DBConnection, ctx context.Context, oid primitive.ObjectID, clientOIDs []primitive.ObjectID) error {
	coll := db.Use(Collection)

	filter := bson.D{{"_id", oid}}
	updateDoc := bson.D{{"$set", bson.D{{"clients", clientOIDs}} }}

	var oldDoc PPA.ClientLabel
	err := coll.FindOneAndUpdate(ctx, filter, updateDoc).Decode(&oldDoc)

	if err == mongodb.ErrNoDocuments {
		return PPA.NewAppError(NotFound, "Label Not Found")
	}

	if err != nil {
		return PPA.InternalError
	}
	return nil
}

func (cl ClientLabel) LogEvent(db *mongo.DBConnection, ctx context.Context, update *PPA.ClientLabel) {
	if err := cl.Update(db, ctx, update.ID, update); err != nil {
		fmt.Println(err)
	}
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
