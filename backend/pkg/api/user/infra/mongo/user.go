package mongo
import (
	"PPA"
	"errors"
	"context"
	"pkg/common/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct{}

func (u User) Create(db *mongo.DBConnection, ctx context.Context, user *PPA.User) (*PPA.User, error) {
	coll := db.Use("users")

	if(emailExists(db, ctx, user.Email)) {
		return nil, errors.New("")
	}

	if _, err := coll.InsertOne(ctx, user); err != nil {
		return nil, errors.New("")
	}

	return user, nil

}

func(u User) ViewById(db *mongo.DBConnection, ctx context.Context, oid primitive.ObjectID) (*PPA.User, error) {
	coll := db.Use("users")

	var user PPA.User
	err := coll.FindOne(ctx, bson.D{{"_id", oid}}).Decode(&user)
	return &user, err
}

func(u User) ViewByEmail(db *mongo.DBConnection, ctx context.Context, email string) (*PPA.User, error) {
	coll := db.Use("users")

	var user PPA.User
	err := coll.FindOne(ctx, bson.D{{"email", email}}).Decode(&user)
	return &user, err
}

func(u User) List(db *mongo.DBConnection, ctx context.Context) (*[]PPA.User, error) {
	coll := db.Use("users")

	//check error?
	cursor, err := coll.Find(ctx, bson.D{{}})
	defer cursor.Close(ctx)

	var users []PPA.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, errors.New("")
	}
	return &users, nil
}

func(u User) Delete(db *mongo.DBConnection, ctx context.Context, oid primitive.ObjectID) error {
	coll := db.Use("deleted_users")

	fetched, err := u.ViewById(db, ctx, oid); if err != nil {
		return errors.New("")//user doesnt exist
	}

	if _, insertErr := coll.InsertOne(ctx, fetched); insertErr != nil {
		return errors.New("") //insert err, db err
	}

	coll = db.Use("users")
	if _, delErr := coll.DeleteOne(ctx, bson.D{{"_id", oid}}); delErr != nil {
		return errors.New("") //db error
	}

	return nil
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
