package mongo

import(
	"PPA"
	"net/http"
	"context"
	"pkg/common/mongo"
	"go.mongodb.org/mongo-driver/bson"
	mongodb "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct{}

const (
	NotFound = http.StatusNotFound
)

func(u User) FindById(db *mongo.DBConnection, ctx context.Context, oid primitive.ObjectID) (*PPA.User, error) {
	coll := db.Use("users")

	var user PPA.User
	if err := coll.FindOne(ctx, bson.D{{"_id", oid}}).Decode(&user); err != nil {
		if err == mongodb.ErrNoDocuments {
			return nil, PPA.NewAppError(NotFound, "User Not Found")
		}
		if err != nil {
			return nil, PPA.InternalError
		}
	}

	return &user, nil
}

func(u User) FindByEmail(db *mongo.DBConnection, ctx context.Context, email string) (*PPA.User, error) {
	coll := db.Use("users")

	var user PPA.User
	err := coll.FindOne(ctx, bson.D{{"email", email}}).Decode(&user)
		if err == mongodb.ErrNoDocuments {
			return nil, PPA.NewAppError(NotFound, "User Not Found")
		}
		if err != nil {
			return nil, PPA.InternalError
		}
	return &user, nil
}

func (u User) Update(db *mongo.DBConnection, ctx context.Context, update *PPA.User) error {
	coll := db.Use("users")

	filter := bson.D{{"_id", update.ID}}
	updateDoc := bson.D{{"$set", update}}

	var oldDoc PPA.User
	err := coll.FindOneAndUpdate(ctx, filter, updateDoc).Decode(&oldDoc)

	if err == mongodb.ErrNoDocuments {
		return PPA.NewAppError(NotFound, "User Not Found")
	}

	if err != nil {
		return PPA.InternalError
	}
	return nil
}

func (u User) FindByToken(db *mongo.DBConnection, ctx context.Context, token string) (*PPA.User, error) {
	coll := db.Use("users")

	var user PPA.User
	err := coll.FindOne(ctx, bson.D{{"token", token}}).Decode(&user)
		if err == mongodb.ErrNoDocuments {
			return nil, PPA.NewAppError(NotFound, "User Not Found")
		}
		if err != nil {
			return nil, PPA.InternalError
		}
	return &user, nil
}
