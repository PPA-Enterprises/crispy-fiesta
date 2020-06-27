package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBConnection struct {
	client *mongo.Client
}

type DBCollection struct {
	Collection *mongo.Collection
}

func NewConnection(host string) (conn *DBConnection) {

	//TODO: Auth
	client, err := mongo.NewClient(options.Client().ApplyURI(host))
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)

	if err != nil {
		panic(err)
	}

	conn = &DBConnection{client}
	return conn

}

func (conn *DBConnection) Use(dbName, tableName string) *DBCollection {

	connect := conn.client.Database(dbName).Collection(tableName)

	return &DBCollection{Collection: connect}

}

func (col *DBCollection) Insert(data bson.D, options *options.InsertOneOptions) (*mongo.InsertOneResult, error) {

	if options == nil {
		result, err := col.Collection.InsertOne(context.Background(), data)
		return result, err

	}

	result, err := col.Collection.InsertOne(context.Background(), data, options)
	return result, err
}
