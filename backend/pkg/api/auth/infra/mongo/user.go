package mongo

import(
	"github.com/gin-gonic/gin"
)

type User struct{}


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
//FindByToken?
