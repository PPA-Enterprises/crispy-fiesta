package common

import (
	"context"
	"time"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBConnection struct {
	Client *mongo.Client
}

var dbConnect *DBConnection

func Init(host string) {
	fmt.Println("hellow init")
	dbConnect = NewConnection(host)
}

func NewConnection(host string) (conn *DBConnection) {

	//TODO: Auth
	//TODO: Client does panic when there is no databse running. We dont want the
	//app to be running unless it has access to databse
	client, err := mongo.NewClient(options.Client().ApplyURI(host))
	if err != nil { panic(err) }
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)

	if err != nil {
		panic(err)
	}

	conn = &DBConnection{client}
	return conn

}

func (conn *DBConnection) Use(dbName, tableName string) *mongo.Collection {
	return conn.Client.Database(dbName).Collection(tableName)
}

func (conn *DBConnection) Session(opts ...*options.SessionOptions) (*mongo.Session, error) {
	var sessionOpts *options.SessionOptions
	if len(opts) <= 0 {
		sessionOpts = options.Session()
	} else {
		sessionOpts = options.MergeSessionOptions(opts...)
	}
	session, err := conn.Client.StartSession(sessionOpts)
	return &session, err
}
