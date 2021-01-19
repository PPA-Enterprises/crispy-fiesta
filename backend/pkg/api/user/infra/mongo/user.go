package mongo
import (
	"PPA"
	"errors"
	"context"
	"pkg/common/mongo"
	"go.go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct{}

func (u User) Create(db *mongo.DBConnection, ctx context.Context, user *PPA.User) (PPA.User, error) {
	coll := db.Connection().Use(db.DefaultDatabase, "users")

	if(emailExists(db, ctx, user.Email)) {
		return nil, errors.New("")
	}

	if _, err := coll.InsertOne(ctx, user); err != nil {
		return nil, errors.New("")
	}

	return user, nil

}
func(u User) View()
func(u User) List()
func(u User) Delete()
func fetchByEmail(db *mongo.DBConnection, ctx context.Context, email string) (PPA.User, error) {
	coll := db.Connection().Use(db.DefaultDatabase, "users")

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
