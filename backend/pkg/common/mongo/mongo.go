package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBConnection struct {
	client *mongo.Client
	databaseName string
}

func Init(host, db string) *DBConnection {
	dbConnect := NewConnection(host, db)
	return dbConnect
}

func NewConnection(host, db string) (conn *DBConnection) {
	//TODO: Auth
	//client, err := mongo.NewClient(options.Client().SetReplicaSet(repl).ApplyURI(host))
	client, err := mongo.NewClient(options.Client().ApplyURI(host))
	if err != nil { panic(err) }

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil { panic(err)}

	//ensure connection was successful
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = client.Ping(pingCtx, nil)
	if err != nil { panic(err)}

	conn = &DBConnection{client: client, databaseName: db}
	return conn
}

func (conn *DBConnection) Use(tableName string) *mongo.Collection {
	return conn.client.Database(conn.databaseName).Collection(tableName)
}

func (conn *DBConnection) Disconnect() {
	err := conn.client.Disconnect(context.Background())
	if err != nil { panic(err) }
}

func (conn *DBConnection) Session(opts ...*options.SessionOptions) (mongo.Session, error) {
	var sessionOpts *options.SessionOptions
	if len(opts) <= 0 {
		sessionOpts = options.Session()
	} else {
		sessionOpts = options.MergeSessionOptions(opts...)
	}
	session, err := conn.client.StartSession(sessionOpts)
	return session, err
}

func Populate(ctx context.Context, coll *mongo.Collection, ids []primitive.ObjectID) (*mongo.Cursor, error) {
	return coll.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
}
