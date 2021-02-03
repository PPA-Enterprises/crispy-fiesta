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
)

const (
	AlreadyExists = http.StatusConflict
	NotFound = http.StatusNotFound
)

type User struct{}

func (u User) Create(db *mongo.DBConnection, ctx context.Context, user *PPA.User) (*PPA.User, error) {
	coll := db.Use("users")

	if(emailExists(db, ctx, user.Email)) {
		return nil, PPA.NewAppError(AlreadyExists, "Email Taken")
	}

	if _, err := coll.InsertOne(ctx, user); err != nil {
		return nil, PPA.InternalError
	}

	return user, nil
}

func(u User) ViewById(db *mongo.DBConnection, ctx context.Context, oid primitive.ObjectID) (*PPA.User, error) {
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

func(u User) ViewByEmail(db *mongo.DBConnection, ctx context.Context, email string) (*PPA.User, error) {
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

func(u User) List(db *mongo.DBConnection, ctx context.Context) (*[]PPA.User, error) {
	coll := db.Use("users")

	//check error?
	cursor, err := coll.Find(ctx, bson.D{{}})
	defer cursor.Close(ctx)

	var users []PPA.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, PPA.InternalError
	}
	return &users, nil
}

func(u User) Delete(db *mongo.DBConnection, ctx context.Context, oid primitive.ObjectID) error {
	coll := db.Use("deleted_users")

	fetched, err := u.ViewById(db, ctx, oid); if err != nil {
		return PPA.NewAppError(NotFound, "User Not Found")
	}

	if _, insertErr := coll.InsertOne(ctx, fetched); insertErr != nil {
		return PPA.InternalError //insert err, db err
	}

	coll = db.Use("users")
	if _, delErr := coll.DeleteOne(ctx, bson.D{{"_id", oid}}); delErr != nil {
		return PPA.InternalError //db error
	}

	return nil
}

func (u User) Update(db *mongo.DBConnection, ctx context.Context, oid primitive.ObjectID, update *PPA.User) error {
	coll := db.Use("users")

	filter := bson.D{{"_id", oid}}
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

func (u User) LogEvent(db *mongo.DBConnection, ctx context.Context, update *PPA.User) {
	if err := u.Update(db, ctx, update.ID, update); err != nil {
		fmt.Println(err)
	}
}

func fetchByEmail(db *mongo.DBConnection, ctx context.Context, email string) (PPA.User, error) {
	coll := db.Use("users")

	var user PPA.User
	err := coll.FindOne(ctx, bson.D{{"email", email}}).Decode(&user)
	return user, err
}

func emailExists(db *mongo.DBConnection, ctx context.Context, email string) bool {
	if _, err := fetchByEmail(db, ctx, email); err != nil {
		return false
	}
	return true
}
